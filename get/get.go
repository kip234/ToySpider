package get

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//键：链接 值：名字
type Links map[string]string

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
func GetSave(url,path string) (err error) {
	//获取http信息
	fmt.Println("URL:",url)

	var resp *http.Response
	resp,err=http.Get(url)

	defer resp.Body.Close()
	//准备存文件
	dire,_:=Spliter(path)
	fmt.Println("dire:",dire)
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
	//over<-id
	return
}