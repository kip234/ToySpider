package main

import (
	"./get"
	"fmt"
	"regexp"
)

//<a class="nav-link active" href="/zh-cn/docs/"><span class="active">文档</span></a>
//链接的正则
var Reurl =`<a.*?href="(/zh-cn/docs.*?)".*?>(.*?)<`
var LinkHead=`https://gin-gonic.com/`

func main() {
	over:=make(chan int)
	text,_:=get.GetUrlText("https://gin-gonic.com/zh-cn/docs/")
	linkstd:=regexp.MustCompile(Reurl)
	links:=linkstd.FindAllStringSubmatch(text,-1)

	for _,data:=range links{
		go get.OnePage("https://gin-gonic.com"+data[1],data[1],data[2]+".html",LinkHead,over)
		//fmt.Println(data[1],data[2])
	}

	num:=len(links)
	for i:=1;i<=num;i++{
		<-over
		fmt.Println("已完成",i,":",num,"个页面")
		//OnePage("https://gin-gonic.com"+data[1],data[1],data[2]+".html",over)
		//fmt.Println(data[1],data[2])
	}

	//OnePage("https://gin-gonic.com/zh-cn/docs/","text","text.html")
}