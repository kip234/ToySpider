package get

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//键：储存目标位置 值：网络地址
type Links map[string]string

const(
	//<a  href="/zh-cn/docs/quickstart/" class="....">快速入门</a>
	RxHtml=`<a.*?href="(.*?)"`

	RxHttp=`(https?://.*?\.com)`


	//<link rel="apple-touch-icon" href="/favicons/apple-touch-icon-180x180.png" sizes="180x180">
	//RxImg=`<img[\s\S]+?src="(http[\s\s]+?)"`
	//RxImg=`<img[\s\S]+?src="(http[\s\s]+?.[a-z]{3})"`
	RxImg=`<link.*?href="(.*?\.png)"`

	//<link rel="preload" href="/scss/main.min.5b9bca55b4e29e6e9c7553311b48ebb1a83b26687587393ed52e2ee07084507a.css"
	RxCss=`<link.*?href="(.*?\.css)"`

	//<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"
	RxJs=`<script src="(.*?\.js)"`
	RxPhone=`(1[3456789]\d)(\d{4})(\d{4})`
	RxEmail=`\w+@\w+\.com`
)

//从URL的连接里匹配信息
func GetHttpMsg(url string, standard *regexp.Regexp,num int) (msg [][]string,err error) {
	var contain *http.Response//网页内容
	contain,err=http.Get(url)
	temp1:=make([]byte,4*1024)//读取缓存
	var temp2 string
	var n int
	n,err=contain.Body.Read(temp1)//读取出来
	for n!=0{
		temp2+=string(temp1[:n])
		n,err=contain.Body.Read(temp1)
	}
	msg=standard.FindAllStringSubmatch(temp2,num)
	return
}

//将网址拆分为目录+文件名
func Spliter(data string) (dire,name string) {
	re:=strings.Split(data,"/")
	num:=len(re)
	for i:=0;i<num-1;i++ {
		dire+=re[i]+"/"
	}
	name = re[num-1]
	return
}

//下载链接的内容至path
func GetSave(url,path string,over chan int) (err error) {
	//获取http信息
	//fmt.Println("URL:",url)

	var resp *http.Response
	resp,err=http.Get(url)
	if err!=nil {
		fmt.Println("GetSave:",url,"err =",err)
		over<-0
		return
	}

	defer resp.Body.Close()
	//准备存文件
	dire,_:=Spliter(path)
	//fmt.Println("dire:",dire)
	err=os.MkdirAll(dire,0777)
	//if err!=nil {
	//	fmt.Println("os.Mkdir err =",err)
	//	return
	//}
	var file *os.File
	file,err=os.Create(path)
	if err!=nil {
		//fmt.Println("GetHttpSave",id,"err =",err)
		//over<-id
		return
	}
	defer file.Close()

	buf:=make([]byte,1024)
	var n int

	//转存
	n,err=resp.Body.Read(buf)
	for n!=0 {
		file.Write(buf[:n])
		n,err=resp.Body.Read(buf)
	}
	over<-1
	return
}

//获取网页文本信息
func GetUrlText(url string) (result string,err error) {
	var resp *http.Response
	resp,err=http.Get(url)
	re,err1:=ioutil.ReadAll(resp.Body)
	err=err1
	result=string(re)
	return
}

//获取字符串中的链接
func GetLink(data string,standard *regexp.Regexp) (result Links) {
	temp:=standard.FindAllStringSubmatch(data,-1)
	result=make(map[string]string,len(temp))
	for _,value:=range temp{
		_,name:=Spliter(value[1])
		result[value[1]]=name//连接：文件名
	}
	return
}

//下载JS,IMG,CSS到path下面的子目录
func OnePage(url,path,name,LinkHead string,over1 chan int)  {
	text,_:=GetUrlText(url)//获取文本

	//获取CSS
	//fmt.Println("CSS")
	cssstd:=regexp.MustCompile(RxCss)
	css:=GetLink(text,cssstd)

	//获取img
	//fmt.Println("IMG")
	imgstd:=regexp.MustCompile(RxImg)
	img:=GetLink(text,imgstd)

	//获取js
	//fmt.Println("JS")
	jsstd:=regexp.MustCompile(RxJs)
	js:=GetLink(text,jsstd)

	//修改HTML里面的链接为相对路径
	for link,name1:=range css{
		text=strings.Replace(text,link,"css/"+name1,-1)
	}
	for link,name1:=range img{
		text=strings.Replace(text,link,"img/"+name1,-1)
	}
	for link,name1:=range js{
		text=strings.Replace(text,link,"js/"+name1,-1)
	}

	//保存HTML
	os.MkdirAll(path,0777)
	temp,_:=os.Create(path+"/"+name)
	temp.Write([]byte(text))
	defer temp.Close()

	over:=make(chan int)

	//开始批量下载
	for link,name:=range css{
		go GetSave(LinkHead+"/"+link,path+"/css/"+name,over)
	}
	for link,name:=range img{
		go GetSave(LinkHead+"/"+link,path+"/img/"+name,over)
	}

	for link,name:=range js{
		go GetSave(link,path+"/js/"+name,over)
	}

	//等待
	for i:=0;i<len(css)+len(img)+len(js);i++ {
		<-over
	}
	over1<-1
}