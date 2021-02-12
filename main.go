package main

import (
	"./get"
	"fmt"
	"os"
	"regexp"
	"strings"
)

//<a class="nav-link active" href="/zh-cn/docs/"><span class="active">文档</span></a>
//链接的正则
var Reurl =`<a.*?href="(/zh-cn/docs/.*?)".*?>(.*?)<`
var LinkHead=`https://gin-gonic.com/`
const directory=`Gin/`

var Home=get.HomePage{//首页
"/zh-cn/docs/",
"文档.html",
}

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
		text=strings.Replace(text,"https://gin-gonic.com"+link,name1,-1)
		text=strings.Replace(text,link,name1,-1)
	}
	text=strings.Replace(text,"https://gin-gonic.com"+Home.Url,Home.Title,-1)
	text=strings.Replace(text,Home.Url,Home.Title,-1)



	//保存
	os.MkdirAll(directory,0777)
	file,err1:=os.Create(directory+name)
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
	cssstd:=regexp.MustCompile(get.RxCss)
	css:=get.GetLink(text,cssstd)
	for link,data:=range css{
		cssLinks[link]=data
	}
	//获取img
	//fmt.Println("IMG")
	imgstd:=regexp.MustCompile(get.RxImg)
	img:=get.GetLink(text,imgstd)
	for link,data:=range img{
		imgLinks[link]=data
	}

	//获取js
	//fmt.Println("JS")
	jsstd:=regexp.MustCompile(get.RxJs)
	js:=get.GetLink(text,jsstd)
	for link,data:=range js{
		jsLinks[link]=data
	}
	////修改HTML里面的链接为相对路径
	//for link,name1:=range css{
	//	text=strings.Replace(text,link,"css/"+name1,-1)
	//}
	//for link,name1:=range img{
	//	text=strings.Replace(text,link,"img/"+name1,-1)
	//}
	//for link,name1:=range js{
	//	text=strings.Replace(text,link,"js/"+name1,-1)
	//}
}

func main() {

	cssLink:=make(get.Links)
	imgLink:=make(get.Links)
	jsLink:=make(get.Links)
	htmlLink:=make(get.Links)

	over:=make(chan int)
	text,_:=get.GetUrlText("https://gin-gonic.com/zh-cn/docs/")
	htmlstd:=regexp.MustCompile(Reurl)
	links:=htmlstd.FindAllStringSubmatch(text,-1)

	for _,data:=range links{//收集资源
		htmlLink[data[1]]=data[2]+".html"
		OnePageSources("https://gin-gonic.com"+data[1],cssLink,imgLink,jsLink)
		//go get.OnePage("https://gin-gonic.com"+data[1],data[1],data[2]+".html",LinkHead,over)
		//fmt.Println(data[1],data[2])
	}
	delete(htmlLink,Home.Url)//把首页去掉

	//反馈
	fmt.Println("HTML")
	for a1,a2:=range htmlLink{
		fmt.Println(a1,":",a2)
	}
	fmt.Println("CSS")
	for a1,a2:=range cssLink{
		fmt.Println(a1,":",a2)
	}
	fmt.Println("JS")
	for a1,a2:=range jsLink{
		fmt.Println(a1,":",a2)
	}
	fmt.Println("IMG")
	for a1,a2:=range imgLink{
		fmt.Println(a1,":",a2)
	}

	num:=len(cssLink)+len(jsLink)+len(imgLink)+len(htmlLink)+1
	fmt.Println("共",num,"个链接")

	//保存HTML文本
	for a1,a2:=range htmlLink{
		go SaveHTML("https://gin-gonic.com"+a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
	}
	go SaveHTML("https://gin-gonic.com"+Home.Url,Home.Title,htmlLink,cssLink,jsLink,imgLink,over)

	//下载其他资源
	for a1,a2:=range cssLink{
		a:=strings.Index(a1,"http")//有没有网址
		if a<0{
			go get.GetSave("https://gin-gonic.com"+a1,directory+"css/"+a2,over)
		}else {
			go get.GetSave(a1,directory+"css/"+a2,over)
		}
	}
	for a1,a2:=range jsLink{
		a:=strings.Index(a1,"http")
		if a<0{
			go get.GetSave("https://gin-gonic.com"+a1,directory+"js/"+a2,over)
		}else {
			go get.GetSave(a1,directory+"js/"+a2,over)
		}
	}
	for a1,a2:=range imgLink{
		a:=strings.Index(a1,"http")
		if a<0{
			go get.GetSave("https://gin-gonic.com"+a1,directory+"img/"+a2,over)
		}else{
			go get.GetSave(a1,directory+"img/"+a2,over)
		}
	}


	for i:=1;i<=num;i++{
		<-over
		fmt.Println("已完成",i,":",num,"")
		//OnePage("https://gin-gonic.com"+data[1],data[1],data[2]+".html",over)
		//fmt.Println(data[1],data[2])
	}

	//OnePage("https://gin-gonic.com/zh-cn/docs/","text","text.html")
}