package main

import (
	"DnslogMsg4Wx/config"
	"DnslogMsg4Wx/plugin"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//发送微信消息
func Pushplus(msg string,token string) bool{
	url := "http://www.pushplus.plus/send"
	// 超时时间：5秒
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url,"application/x-www-form-urlencoded",strings.NewReader("token="+token+"&content="+msg))
	if err != nil {
		panic(err)
		return false
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	//fmt.Println(string(content))
	return true
}

func main(){
	//大致逻辑 循环遍历ceye接口  每10分钟推送一次消息 //todo 后面可以改成白天五分钟，晚上一小时
	fmt.Println("begin")
	//time.AfterFunc(5*time.Second,test)//直接来个死循环就好了
	config.Init()
	for true {
		//test()
		//fmt.Println(config.Pushplus_token)
		msg ,err:= plugin.Ceye(config.Ceye_token,"dns")
		if err == "HaveMsg"{
			Pushplus(msg,config.Pushplus_token)
		}
		msg ,err= plugin.Ceye(config.Ceye_token,"http")
		if err == "" && msg!=""{
			Pushplus(msg,config.Pushplus_token)
		}
		time.Sleep(10*time.Minute)
	}
}
