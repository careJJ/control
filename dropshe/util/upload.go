package util

import (
	"github.com/astaxie/beego"
	"path"
	"github.com/weilaihui/fdfs_client"
	"strings"
)
func RespFunc(this *beego.Controller, resp map[string]interface{}) {
	//1.把容器传递给前段
	this.Data["json"] = resp
	//2.指定传递方式 也就是值前端的ajax
	this.ServeJSON()
}

//fastdfs上传
func FdsUploadImage(this *beego.Controller,file string)string{
	//上传图片
	//在 beego 中你可以很容易的处理文件上传，就是别忘记在你的 form
	// 表单中增加这个属性 enctype="multipart/form-data"，否则你的浏览器不会传输你的上传文件。
	image,head,err:=this.GetFile(file)
	//获取图片
	//返回值 文件二进制流  文件头    错误信息
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = "add.html"

	}
	defer image.Close()
	//校验文件大小
	if head.Size >5000000{
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片数据过大"
		this.TplName = "add.html"

	}

	//校验格式 获取文件后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "上传文件格式错误"
		this.TplName = "add.html"

	}
	//防止重名
	//fileName := time.Now().Format("200601021504052222")
	//把上传的文件存储到项目文件夹
	//this.SaveToFile("images","./static/img/"+fileName+ext)
	//将图片存入fastdfs
	//先获取一个[]byte
	fileBuffer := make([]byte,head.Size)
	//把文件数据读入到fileBuffer中
	image.Read(fileBuffer)
	//获取client对象

	client,err := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	if err!=nil{
		beego.Error(err)
	}

	//上传
	fdfsresponse,_:=client.UploadByBuffer(fileBuffer,ext[1:])
	 return fdfsresponse.RemoteFileId
}

//将字符串拆分成切片的方法
func String2slice(str string)[]string{
	a:=strings.Split(str,",")
	return a
}


//切片换字符串
func Slice2string(str []string)string{
	a:=strings.Join(str,",")
	return a
}