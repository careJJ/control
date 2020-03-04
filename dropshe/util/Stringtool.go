package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"image/color"
	"image/jpeg"
	"time"
	"crypto/md5"
)
//有关数据、图片等的处理


//将字符串加密成 md5
func String2md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)

	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}
/*
//RandomString 在数字、大写字母、小写字母范围内生成num位的随机字符串
func RandomString(length int) string {
	// 48 ~ 57 数字
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
	// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
	// 小于36在大写范围内随机，其他在小写范围随机
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	for i := 0; i < length; i++ {
		t := rand.Intn(62)
		if t < 10 {
			result = append(result, strconv.Itoa(rand.Intn(10)))
		} else if t < 36 {
			result = append(result, string(rand.Intn(26)+65))
		} else {
			result = append(result, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(result, "")
}
*/
const Salt  ="dropshe888"

//加盐后再次MD5  转换16进制是否有问题
func AddSalt2string(str string)string{
	saltstr:=str+Salt
	data := []byte(saltstr)
	newmd5pwd:=md5.Sum(data)
	return fmt.Sprintf("%x", newmd5pwd)
}


//生成图片验证码
func CheckImages(this *beego.Controller) {
	cap := captcha.New()
	//通过句柄调用 字体文件
	if err := cap.SetFont("comic.ttf"); err != nil {
		beego.Info("没有字体文件")
		panic(err.Error())
	}
	//设置图片的大小
	cap.SetSize(91, 41)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	//SetFrontColor(colors ...color.Color)  这两个颜色设置的函数属于不定参函数
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	//生成图片 返回图片和 字符串(图片内容的文本形式)
	img, str := cap.Create(4, captcha.NUM)

	beego.Info(str)

	a := *img      //解引用

	email := this.GetSession("email")

	//连接redis
	bm, err := RedisOpen(Server_name, Redis_addr,
		Redis_port, Redis_dbnum)
	if err != nil {
		beego.Error("连接数据库失败")
	}
	//存入redis，有效期时间5min
	bm.Put(email.(string), str, time.Second*300)
	//将图片转成base64
	imgbuff:=bytes.NewBuffer(nil)
	//将图片写入到buff
	jpeg.Encode(imgbuff,a,nil)
	//开辟存储空间
	dist:=make([]byte,50000)
	//将buff转成base64
	base64.StdEncoding.Encode(dist,imgbuff.Bytes())
}
