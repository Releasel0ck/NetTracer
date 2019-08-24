package util

import (
	"io/ioutil"
	"os"
	"strings"
)

//检查错误
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//删除切片重复项
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

//检查文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//创建文件
func CreateFile(filename string) bool {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return false
	} else {
		return true
	}

}

//获取所有数据库文件名
func GetArbitraryDB() string {
	var fs []string
	files, _ := ioutil.ReadDir("./db")
	for _, f := range files {
		if f.Name()[0] == '.' {
			continue
		}
		fs = append(fs, strings.Replace(f.Name(), ".db", "", 1))
	}
	return strings.Join(fs, "$")
}

//读取端口配置
func ReadPortConfig() []string {
	var r []string
	name := "port.config"
	if contents, err := ioutil.ReadFile(name); err == nil {
		result := strings.Replace(string(contents), "\n", "", 1)
		r = strings.Split(result, ",")
		return r
	} else {
		return r
	}
}
