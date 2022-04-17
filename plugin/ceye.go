package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Msg struct {
	Meta Meta `json:"meta"`
	Data []Data `json:"data"`
}
type Meta struct {
	Code int `json:"code"`
	Message string `json:"message"`
}
type Data struct {
	ID string `json:"id"`
	Name string `json:"name"`
	RemoteAddr string `json:"remote_addr"`
	CreatedAt string `json:"created_at"`
}

func Ceye(token string,types string) (string,string){
	url := "http://api.ceye.io/v1/records?token="
	// 超时时间：5秒
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url+token+"&type="+types)
	if err != nil {
		panic(err)
		return "",err.Error()
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
			return "",err.Error()
		}
	}
	res := &Msg{}
	json.Unmarshal(result.Bytes(),res)
	if res.Meta.Code == 200 &&len(res.Data)>0 {
		//在tmp下创建一个文件 记录最新的id
		lastID := 0
		filePath := "/tmp/Dnslog.log"
		lastIDstr, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("read fail", err)
		}
		lastID, _ =strconv.Atoi(string(lastIDstr))
		msgFlag := false
		msg :="<html>\n<body>\n\n<table border=\"1\">\n  <tr>\n    <th>host</th>\n    <th>RemoteAddr</th>\n  </tr> "
		for i:=0;i < len(res.Data); i++ {
			id, _ := strconv.Atoi(res.Data[i].ID)
			if lastID < id {
				//需要先解决一个重复的问题，每次只取出最新的消息
				msg = msg+" <tr><td>"+res.Data[i].Name+"</td><td>"+res.Data[i].RemoteAddr+"</td></tr>"
				msgFlag = true
			}
		}
		msg = msg+"</table>\n\n</body>\n</html>"
		//覆盖写入最新的id
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("file create failed. err: " + err.Error())
		} else {
			// offset
			//os.Truncate(filename, 0) //clear
			n, _ := f.Seek(0, os.SEEK_END)
			_, err = f.WriteAt([]byte(res.Data[0].ID), n)
			defer f.Close()
		}
		if msgFlag{
			return msg,"HaveMsg"
		}else {
			return "","NoMsg"
		}
	}else{
		return "","NoMsg"
	}
}
