package log

import (
	"os"
	"strconv"
	"time"
)

func Log() (err error) {
	time:=time.Now()
	//创建目录
	dire:="log/"+strconv.FormatInt(int64(time.Year()),10)+"-"+strconv.FormatInt(int64(time.Month()),10)+"-"+strconv.FormatInt(int64(time.Day()),10)
	//fmt.Println(dire)
	os.MkdirAll(dire,0777)
	//创建文件
	fileName :=strconv.FormatInt(int64(time.Hour()),10)+"-"+strconv.FormatInt(int64(time.Minute()),10)+"-"+strconv.FormatInt(int64(time.Second()),10)+".log"
	//fmt.Println(fileName)
	var file *os.File
	file,err=os.Create(dire+"/"+fileName)
	if err!=nil {
		return
	}
	//重定向
	os.Stdout=file
	os.Stderr=file
	return
}
