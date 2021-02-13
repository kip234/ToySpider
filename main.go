package main

import (
	"./config"
	"./get"
	"fmt"
	"os"
	"regexp"
	"strings"
	"./log"
)

var conf config.Config

//下载文本，链接更换
func SaveHTML(url,name string,html,css,js,img get.Links,over chan int)  {
	text,err1:=get.GetUrlText(url)//获取文本
	if err1!=nil {
		fmt.Println("!Error",url,":",err1)
	}

	//修改HTML里面的链接为相对路径
	for name1,link:=range css{
		text=strings.Replace(text,link,"css/"+name1,-1)
	}
	for name1,link:=range img{
		text=strings.Replace(text,link,"img/"+name1,-1)
	}
	for name1,link:=range js{
		text=strings.Replace(text,link,"js/"+name1,-1)
	}
	for name1,link:=range html{//即使同一个页面也要给我整2个花样...
		text=strings.Replace(text,"\"https://"+conf.Link+link+"\"","\""+name1+"\"",-1)
		text=strings.Replace(text,"\"http://"+conf.Link+link+"\"","\""+name1+"\"",-1)
		text=strings.Replace(text,"\""+link+"\"","\""+name1+"\"",-1)
	}

	//保存
	os.MkdirAll(conf.Directory,0777)
	var file *os.File
	file,err1=os.Create(conf.Directory+name)
	if err1!=nil {
		fmt.Println("!Error",url,":",err1)
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
	text,err:=get.GetUrlText(url)//获取文本

	if err!=nil {
		fmt.Println("!Error ",url,":",err)
	}

	//获取CSS
	//fmt.Println("CSS")
	cssstd:=regexp.MustCompile(conf.RxCss)
	css:=get.GetLink(text,cssstd)
	for data,link:=range css{
		data=strings.Trim(data," \t\n\r\f")//除去空白字符
		cssLinks[data]=link
	}
	//获取img
	//fmt.Println("IMG")
	imgstd:=regexp.MustCompile(conf.RxImg)
	img:=get.GetLink(text,imgstd)
	for data,link:=range img{
		data=strings.Trim(data," \t\n\r\f")//除去空白字符
		imgLinks[data]=link
	}

	//获取js
	//fmt.Println("JS")
	jsstd:=regexp.MustCompile(conf.RxJs)
	js:=get.GetLink(text,jsstd)
	for data,link:=range js{
		data=strings.Trim(data," \t\n\r\f")//除去空白字符
		jsLinks[data]=link
	}
}

func main() {

	err:=log.Log()
	if err!=nil {
		fmt.Println(err)
		return
	}

	conf.Init()//初始化配置

	cssLink:=make(get.Links)//记录所有页面的CSS资源链接
	imgLink:=make(get.Links)//记录所有页面的img资源链接
	jsLink:=make(get.Links)//记录所有页面的js资源链接
	htmlLink:=make(get.Links)//记录所有页面的链接

	over:=make(chan int)//并发信道
	text,_:=get.GetUrlText(conf.Goal)//获取文本
	htmlstd:=regexp.MustCompile(conf.RxHtml)//匹配HTML的正则对象
	links:=htmlstd.FindAllStringSubmatch(text,-1)//匹配HTML链接

	for _,data:=range links{//收集资源
		for _,illegal:=range "\\/:*?\"<>|"{//替换作为文件名时的非法字符
			data[2]=strings.Replace(data[2],string(illegal),"+",-1)
		}
		//除去空白字符
		data[2]=strings.Trim(data[2]," \t\n\r\f")
		htmlLink[data[2]+".html"]=data[1]
		OnePageSources(conf.LinkHead+data[1],cssLink,imgLink,jsLink)
	}

	//反馈
	for value,link:=range htmlLink{
		fmt.Println(link,"=>",value)
	}
	for value,link:=range cssLink{
		fmt.Println(link,"=>",value)
	}
	for value,link:=range jsLink{
		fmt.Println(link,"=>",value)
	}
	for value,link:=range imgLink{
		fmt.Println(link,"=>",value)
	}

	num:=len(cssLink)+len(jsLink)+len(imgLink)+len(htmlLink)//统计链接总数
	fmt.Println("共",num,"个链接")

	//保存HTML文本
	for a2,a1:=range htmlLink{
		a:=strings.Index(a1,"http://")//有没有协议
		a+=strings.Index(a1,"https://")//找不到返回-1
		if a>= -1{//有其中一个-啥也不缺
			go SaveHTML(a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		}else if a=strings.Index(a1,".com");a>=0{//只缺协议
			a1=strings.Trim(a1,"/")//除去多余的 //
			go SaveHTML("https://"+a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		}else {//都缺
			go SaveHTML(conf.LinkHead+a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		}
	}

	//下载其他资源
	for a2,a1:=range cssLink{
		a:=strings.Index(a1,"http://")//有没有
		a+=strings.Index(a1,"https://")//有没有
		if a>= -1{
			go get.GetSave(a1,conf.Directory+"css/"+a2,over)
		}else if a=strings.Index(a1,".com");a>=0{
			a1=strings.Trim(a1,"/")//除去多余的 //
			go get.GetSave("https://"+a1,conf.Directory+"css/"+a2,over)
		}else{
			go get.GetSave(conf.LinkHead+a1,conf.Directory+"css/"+a2,over)
		}
	}
	for a2,a1:=range jsLink{
		a:=strings.Index(a1,"http://")
		a+=strings.Index(a1,"https://")
		if a>= -1{
			go get.GetSave(a1,conf.Directory+"js/"+a2,over)
		}else if a=strings.Index(a1,".com");a>=0{
			a1=strings.Trim(a1,"/")//除去多余的 //
			go get.GetSave("https://"+a1,conf.Directory+"js/"+a2,over)
		}else{
			go get.GetSave(conf.LinkHead+a1,conf.Directory+"js/"+a2,over)
		}
	}
	for a2,a1:=range imgLink{
		a:=strings.Index(a1,"http://")//判断有没有网址
		a+=strings.Index(a1,"https://")//找不到返回-1
		if a>= -1{
			go get.GetSave(a1,conf.Directory+"img/"+a2,over)
		}else if a=strings.Index(a1,".com");a>=0{
			a1=strings.Trim(a1,"/")//除去多余的 //
			go get.GetSave("https://"+a1,conf.Directory+"img/"+a2,over)
		}else{
			go get.GetSave(conf.LinkHead+a1,conf.Directory+"img/"+a2,over)
		}
	}


	for i:=1;i<=num;i++{//等待协程退出
		<-over
		fmt.Println("已完成",i,":",num,"")
	}
}