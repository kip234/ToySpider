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
	var file *os.File
	file,err=os.Create(path)
	if err!=nil {
		over<-0
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
	if err!=nil {
		return
	}
	re,err1:=ioutil.ReadAll(resp.Body)
	err=err1
	if err!=nil {
		return
	}
	result=string(re)
	return
}

//获取字符串中的链接
func GetLink(data string,standard *regexp.Regexp) (result Links) {
	temp:=standard.FindAllStringSubmatch(data,-1)
	result=make(Links,len(temp))
	for _,value:=range temp{
		_,name:=Spliter(value[1])
		result[value[1]]=name//连接：文件名
	}
	return
}

