package main

import (
	"./config"
	"./get"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var conf config.Config

//下载文本，链接更换
func SaveHTML(url,name string,html,css,js,img get.Links,over chan int) (err error) {
	text,_:=get.GetUrlText(url)//获取文本

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
	for link,name1:=range html{
		text=strings.Replace(text,"\""+conf.LinkHead+link+"\"","\""+name1+"\"",-1)
		text=strings.Replace(text,"\""+link+"\"","\""+name1+"\"",-1)
	}

	//保存
	os.MkdirAll(conf.Directory,0777)
	file,err1:=os.Create(conf.Directory+name)
	if err1!=nil {
		err=err1
		over<-0
		return
	}
	defer file.Close()
	file.Write([]byte(text))
	over<-1
	return
}

//查找本页面需要的CSS、JS、IMG资源
func OnePageSources(url string,cssLinks,imgLinks,jsLinks get.Links)  {
	text,_:=get.GetUrlText(url)//获取文本

	//获取CSS
	//fmt.Println("CSS")
	cssstd:=regexp.MustCompile(conf.RxCss)
	css:=get.GetLink(text,cssstd)
	for link,data:=range css{
		cssLinks[link]=data
	}
	//获取img
	//fmt.Println("IMG")
	imgstd:=regexp.MustCompile(conf.RxImg)
	img:=get.GetLink(text,imgstd)
	for link,data:=range img{
		imgLinks[link]=data
	}

	//获取js
	//fmt.Println("JS")
	jsstd:=regexp.MustCompile(conf.RxJs)
	js:=get.GetLink(text,jsstd)
	for link,data:=range js{
		jsLinks[link]=data
	}
}

func main() {

	//m:="Multipart/Urlencoded 表单"
	//for _,illegal:=range "\\/:*?\"<>|"{//替换作为文件名时的非法字符
	//	m=strings.Replace(m,string(illegal),"+",-1)
	//	fmt.Println(string(illegal))
	//}
	//m=m+".html"
	//fmt.Println(m)


	conf.Init()

	cssLink:=make(get.Links)
	imgLink:=make(get.Links)
	jsLink:=make(get.Links)
	htmlLink:=make(get.Links)

	over:=make(chan int)
	text,_:=get.GetUrlText(conf.Goal)
	htmlstd:=regexp.MustCompile(conf.RxHtml)
	links:=htmlstd.FindAllStringSubmatch(text,-1)

	for _,data:=range links{//收集资源
		//htmlLink[data[1]]=data[2]+".html"
		for _,illegal:=range "\\/:*?\"<>|"{//替换作为文件名时的非法字符
			data[2]=strings.Replace(data[2],string(illegal),"+",-1)
		}
		htmlLink[data[1]]=data[2]+".html"
		OnePageSources(conf.LinkHead+data[1],cssLink,imgLink,jsLink)
	}

	for link,value:=range htmlLink{
		fmt.Println(link,"=",value)
	}

	num:=len(cssLink)+len(jsLink)+len(imgLink)+len(htmlLink)
	fmt.Println("共",num,"个链接")

	//保存HTML文本
	for a1,a2:=range htmlLink{
		//a:=strings.Index(a1,"http")//有没有网址
		//if a<0{
			go SaveHTML(conf.LinkHead+a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		//}else {
		//	go SaveHTML(a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		//}
	}

	//下载其他资源
	for a1,a2:=range cssLink{
		a:=strings.Index(a1,"http")//有没有网址
		if a<0{
			go get.GetSave(conf.LinkHead+a1,conf.Directory+"css/"+a2,over)
		}else {
			go get.GetSave(a1,conf.Directory+"css/"+a2,over)
		}
	}
	for a1,a2:=range jsLink{
		a:=strings.Index(a1,"http")
		if a<0{
			go get.GetSave(conf.LinkHead+a1,conf.Directory+"js/"+a2,over)
		}else {
			go get.GetSave(a1,conf.Directory+"js/"+a2,over)
		}
	}
	for a1,a2:=range imgLink{
		a:=strings.Index(a1,"http")
		if a<0{
			go get.GetSave(conf.LinkHead+a1,conf.Directory+"img/"+a2,over)
		}else{
			go get.GetSave(a1,conf.Directory+"img/"+a2,over)
		}
	}


	for i:=1;i<=num;i++{
		<-over
		fmt.Println("已完成",i,":",num,"")
	}
}