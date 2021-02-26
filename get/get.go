package get

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

//键：储存目标位置 值：网络地址
type Links map[string]string

var WG sync.WaitGroup

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
	WG.Add(1)
	defer WG.Done()
	var resp *http.Response
	resp,err=http.Get(url)
	if err!=nil {
		fmt.Println("GetSave:",url,"err =",err)
		return
	}

	defer resp.Body.Close()
	//准备存文件
	dire,_:=Spliter(path)
	//fmt.Println("dire:",dire)
	err=os.MkdirAll(dire,0777)//准备文件夹
	var file *os.File
	file,err=os.Create(path)//创建文件
	if err!=nil {
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
		result[name]=value[1]//文件名:连接
	}
	return
}

