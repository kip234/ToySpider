# ToySpider

----

## 目录

[目录](#目录)

[写在前面](#写在前面)

[如何使用](#如何使用)

[业务逻辑](#业务逻辑)

[吐槽](#吐槽)

[写在最后](#写在最后)

## 写在前面

[目录](#目录)

鄙人是菜鸟，网警同志应该不会来抓我吧。

>顾名思义，玩具而已。
>
>用于爬取网站上的教程
>
>力求完美复制：HTML，CSS，JavaScript，图片。
>
>应该不会有人伤心病狂到对用爱发电的网站下黑手吧。
>
>下面的内容你可能会看到“graph”开头的代码，那是因为GitHub好像不支持mermaid，下载到本地后应该能正常显示。

最终效果

>爬取成套的教程，且各项指标正常，页面可正常转跳

## 如何使用

[目录](#目录)

​	[Lv0](#Lv0)

​	[Lv1](#Lv1)

​	[Lv2](#Lv2)

### Lv0

看见那个.exe结尾的东西了吗？点它！然后等，正常情况下是没有任何输出的，在log文件夹里面有输出的日志。

默认是Gin的中文教程，json文件夹里面有其他教程对应的配置。

然后就复制粘贴

### Lv1

通过更改配置文件进行修改

你需要

> JSON会看就行
>
> 一点点前端基础
>
> “深厚”的**正则表达式**功底

```json
{
  "Directory": "Gin/",
  "LinkHead": "https://gin-gonic.com",
  "RxImg": "<link.*?href=\"(.*?\\.png)\"",
  "RxHtml": "<a.*?href=\"(/zh-cn/docs/.*?)\".*?>(.*?)<",
  "RxCss": "<link.*?href=\"(.*?\\.css)\"",
  "RxJs": "<script src=\"(.*?\\.js)\"",
  "Goal": "https://gin-gonic.com/zh-cn/docs/",
  "Link":"gin-gonic.com",
   "Debug":false,
   "RoutingGroup":""
}
```

| “键”         | “值”                                                         |
| ------------ | ------------------------------------------------------------ |
| Directory    | 保存文件的目录，一定要以“/”结尾，如果以“/”开头则会在当前盘符的根目录下面创建文件夹 |
| LinkHead     | 看提示、顾名思义，填协议+域名                                |
| Goal         | 目标网页地址，将会在这个网页上匹配其他网页的链接             |
| Link         | 顾名思义还是算了吧，填域名                                   |
| Rx…          | 填对应的正则表达式、如果处理不好对方服务器崩了来找你可别提起我 |
| Debug        | 调试模式下不会下载                                           |
| RoutingGroup | 链接路由组-页面嵌入的链接没有则填上                          |

==**正则表达式一定要描述准确，尽量保证只匹配你真正需要的**==

==**一定要注意字符的转义，正则和字符串两个层面都要考虑！**==

### Lv2

直接改代码

你需要


* GOLang
    > 并发编程
    > channel通信
    >
    > 文本处理
    >
    > 文件管理
    >
    > 对正则的支持
    >
    > JSON处理
    >
    > HTTP请求
    
* 心态
	>我的代码可能反人类
	
	代码结构：
	
	```mermaid
	graph TD
	/==>/log==>log.go==>fun1(func log)
	/==>/get==>go1(get.go)==>type1(type Links)
	go1==>fun2(func Spliter)
	go1==>fun3(func GetSave)
	go1==>fun4(func GetUrlText)
	go1==>fun5(func GetLink)
	/==>/config==>go2(base.go)==>type2(type Config struct)
	go2==>fun6(func Init)
	/==>/json==>.json
	/==>go3(main.go)==>fun7(func main)
	go3==>fun8(func SaveHTML)
	go3==>fun9(func OnePageSources)
	```
	
	

## 业务逻辑

[目录](#目录)

业务逻辑

```mermaid
graph TD
url(URL)==>1(获取HTML文本)
1==>2(统计HTML链接)==>6(提取标题)
1==>3(统计CSS链接)
1==>4(统计IMG链接)
1==>5(统计JS链接)
3==>7(提取文件名)
4==>7
5==>7
6==>8(将HTML链接改为对应的标题,CSS等资源链接改为文件夹+文件名)
7==>8
8==>9(以标题名作为文件名保存HTML文本)
8==>10(下载资源到对应目录)
```

目录结构

```mermaid
graph TD
...==>1(指定文件夹)
1==>.html
1==>/css==>.css
1==>/js==>.js
1==>/img==>.jpg/.png
```

工作流程

```mermaid
graph LR
下载HTML文本==>提取HTML链接==>统计所有页面的资源==>下载
```



## 吐槽

[目录](#目录)

### 图片

关于图片的引用，在我印象里我见到了三种不同的方式

1.缺域名和协议的

```html
 <img width="128" height="128" src="/wp-content/themes/runoob/assets/images/qrcode.png" />
```

> 这种浏览器会默认在本站内查找

2.缺协议的

```html
<img src="//static.runoob.com/images/dashang/close.jpg" alt="取消" />
```

> 访问时浏览器自动添加

3.啥也不缺的

```html
<img src="//www.runoob.com/wp-content/uploads/2019/03/01986C87-7E19-4497-878E-AE996AFC088E.jpg">
```

> 如果全长这样该多好

由于鄙人只会用<kbd>http.Get()</kbd>,然而这东西只对上面第三种有效，于是就造就了下面的判断结构

```go
//保存HTML文本
	for a2,a1:=range htmlLink{
		a:=strings.Index(a1,"http://")//有没有协议
		a+=strings.Index(a1,"https://")//找不到返回-1
		if a>= -1{//有其中一个-啥也不缺
			go SaveHTML(a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		}else if a=strings.Index(a1,".com");a>=0{//只缺协议
			a1=strings.Trim(a1,"/")//除去多余的 /
			go SaveHTML("https://"+a1,a2,htmlLink,cssLink,jsLink,imgLink,over)//http和https好像都行
		}else {//缺域名和协议的-LinkHead上面有提及
			go SaveHTML(conf.LinkHead+a1,a2,htmlLink,cssLink,jsLink,imgLink,over)
		}
	}
```

> 这大过年的，整的我简直不要太那啥。

### 链接

#### 之不同的页面

```html
<a  href="/zh-cn/docs/" class="......">文档</a>
<a  href="/zh-cn/docs/introduction/" class="......">介绍</a>
```

> 不知道发现没有，这俩货的链接是目录，你中有我
>
> 由于鄙人为避免链接重复使用map储存，而map在迭代时是无序的

于是在改链接的时候就出现过下面的情况

```html
<a  href="文档.html" class="......">文档</a>
<a  href="文档.htmlintroduction/" class="......">介绍</a>
```

> 第2行的链接浏览器要是能找到我跟它姓

最后的解决办法：把双引号一起匹配

即：

```go
text=strings.Replace(text,"\""+link+"\"","\""+name1+"\"",-1)
```

#### 之同一个页面

```html
<a href="https://gin-gonic.com/zh-cn/docs/">文档</a>
<a  href="/zh-cn/docs/" class="....">文档</a>
```

> 这…看来我只能来硬的

```go
for name1,link:=range html{
	text=strings.Replace(text,"\"https://"+conf.Link+link+"\"","\""+name1+"\"",-1)
	text=strings.Replace(text,"\"http://"+conf.Link+link+"\"","\""+name1+"\"",-1)
	text=strings.Replace(text,"\""+link+"\"","\""+name1+"\"",-1)
}
```

#### 路由组

> 在https://gorm.io/zh_CN/docs/indexes.html中的链接

```html
<a href="index.html" class="sidebar-link">概述</a>
```

> 点进去后地址栏是:https://gorm.io/zh_CN/docs/index.html,它没有回到根目录，任然处在当前路由下

### 页面标题

```html
 <a class="..." id="..." href="/zh-cn/docs/examples/rendering/">XML/JSON/YAML/ProtoBuf 渲染</a>
```

> 这…自带文件夹？

```go
for _,illegal:=range "\\/:*?\"<>|"{//替换作为文件名时的非法字符
	data[2]=strings.Replace(data[2],string(illegal),"+",-1)
}
```

## 写在最后

[目录](#目录)

编写于鄙人寒假，对我是学生，所以更新什么的就不可能了，我作业还没写来着，应该也没有人会看这个吧，毕竟GitHub上星星排前面的简直太强了，而且几乎全是中国人。不由得想起前些年频频出现的“现象”。