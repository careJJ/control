package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"strconv"
	"math"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowIndex() {
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles [] models.Article
	qs.All(&articles)
	typeName := this.GetString("select")
	var count int64
	//总文章数
	if typeName == "" {
		count, _ = qs.RelatedSel("ArticleType").Count()
	} else {
		count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}

	//每页文章数
	pageindex := 2
	//总页数
	pageCount := math.Ceil((float64(count)) / float64(pageindex))
	pageNum, err := this.GetInt("pageNum")
	if err != nil {
		pageNum = 1
	}
	beego.Info("数据总页数未:", pageNum)



	if typeName == "" {
		qs.Limit(pageindex, pageindex*(pageNum-1)).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pageindex, pageindex*(pageNum-1)).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}

	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes

	this.Data["articles"] = articles
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = pageNum

	this.Data["TypeName"]=typeName
	this.TplName = "index.html"

}

func (this *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes
	this.TplName = "add.html"

}

func (this *ArticleController) HandleAddArticle() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	if content == "" || articleName == "" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "add.html"
		return
	}
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	//条件判断
	if head.Size > 500000 {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "文件过大"
		this.TplName = "add.html"
		return
	}
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "文件个是错误"
		this.TplName = "add.html"
		return
	}
	Filename := time.Now().Format("2006-01-02-15-04-05-2222")
	//把上传的文件存储到项目文件夹
	this.SaveToFile("uploadname", "./static/img"+Filename+ext)

	o := orm.NewOrm()
	var article models.Article
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")
	article.ArticleType = &articleType
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img" + Filename + ext
	_, err = o.Insert(&article)
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "数据插入失败"
		this.TplName = "add.html"
		return
	}
	this.Redirect("/index", 302)
}

func (this *ArticleController) ShowContent() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取文章ID错误")
		this.Redirect("/index", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	article.ReadCount += 1
	o.Update(&article)
	this.Data["article"] = article
	this.TplName = "content.html"
}

func (this *ArticleController) ShowUpdate() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取ID失败")
		this.Redirect("/index", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	this.Data["article"] = article
	this.TplName = "update.html"

}

func UploadFile(this *ArticleController, filePath string, errhtml string) string {
	file, head, err := this.GetFile(filePath)
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = errhtml
		return ""
	}
	defer file.Close()
	//条件判断
	if head.Size > 5000000 {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "文件过大"
		this.TplName = errhtml
		return ""
	}
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "文件个是错误"
		this.TplName = errhtml
		return ""
	}
	Filename := time.Now().Format("2006-01-02-15-04-05-2222")
	//把上传的文件存储到项目文件夹
	this.SaveToFile(filePath, "./static/img/"+Filename+ext)
	return "/static/img/" + Filename + ext
}

func (this *ArticleController) HandleUpdate() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	savePath := UploadFile(this, "uploadname", "update.html")
	id, _ := this.GetInt("id")
	if articleName == "" || content == "" || savePath == "" {
		beego.Error("获取文件错误")

		this.Redirect("/index?id="+strconv.Itoa(id), 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	article.Title = articleName
	article.Content = content
	article.Img = savePath
	o.Update(&article)
	this.Redirect("/index", 302)
}

func (this *ArticleController) HandleDelete() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取ID失败")
		this.Redirect("/index", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article

	article.Id = id

	o.Delete(&article, "id")
	this.Redirect("/index", 302)

}

func (this *ArticleController) HandleDeleteType() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取ID失败")
		this.Redirect("/addType", 302)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType, "id")
	this.Redirect("/addType", 302)

}

func (this *ArticleController) ShowAddType() {

	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes

	this.TplName = "addType.html"
}

func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Error("类型名称传输失败")
		this.Redirect("/addType", 302)
		return
	}
	//插入操作
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)
	this.Redirect("/addType", 302)
}
