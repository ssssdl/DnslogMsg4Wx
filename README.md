# DnslogMsg4Wx
一个小小有用的脚本，dnslog消息推送程序，利用pushplus接口，将dnslog数据推送至微信

## 使用方法
- 关注pushplus 推送加公众号，获取到token添加到config/config.go 中
- 将ceye的token，添加到config/config.go 中
- 运行main即可，每十分钟自动推送一次


> 会在/tmp下创建log文件，用于记录ceye最新id，防止重复推送
