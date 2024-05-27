package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
)

func main() {
	// 绑定路由
	http.HandleFunc("/", sign)
	// 启动监听
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println("服务器启动失败！")
	}
}
func sign(response http.ResponseWriter, request *http.Request) {
	//解析参数
	err := request.ParseForm()
	if err != nil {
		return
	}
	var token string = "jxtd_server_token"
	// 获取参数
	signature := request.FormValue("signature")
	timestamp := request.FormValue("timestamp")
	nonce := request.FormValue("nonce")
	echostr := request.FormValue("echostr")
	//将token、timestamp、nonce三个参数进行字典序排序
	var tempArray = []string{token, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	h := sha1.New()
	h.Write([]byte(sha1String))
	sha1String = hex.EncodeToString(h.Sum([]byte("")))
	//获得加密后的字符串可与signature对比
	if sha1String == signature {
		_, err := response.Write([]byte(echostr))
		if err != nil {
			fmt.Println("response failure")
		} else {
			fmt.Println("verify success")
		}
	} else {
		fmt.Println("verify failure")
	}
}
