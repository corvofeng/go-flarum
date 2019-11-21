# goyoubbs

```
CGO_ENABLED=0 GOOS=linux go build
```

## 轻论坛功能

- 用户：用户名密码登录、微博、QQ 登录
- 用户上传文件存储：本地、七牛、又拍云
- 根据标题自动提取tag 或管理员手工设置tag
- 根据tag 获取相关文章
- 站内搜索标题、内容
- 内容里链接点击计数
- 自动安装HTTPS并自动更新

## 快速使用

即使你没有接触过golang， 按照下面步骤也能快速部署

以linux 64位系统为例，依次输入下面几行命令即可：

下载主程序包、静态文件包
（最新版下载[https://github.com/ego008/goyoubbs/releases](https://github.com/ego008/goyoubbs/releases) 选择适合你系统的包）
```
wget https://github.com/ego008/goyoubbs/releases/download/current/goyoubbs-linux-amd64.zip
wget https://github.com/ego008/goyoubbs/releases/download/current/site.zip
unzip goyoubbs-linux-amd64.zip
unzip site.zip
./goyoubbs
```

如果出现类似下面的提示，说明已正常启动：

```
2017/12/06 16:24:42 MainDomain: http://127.0.0.1:8082
2017/12/06 16:24:42 youdb Connect to mydata.db
2017/12/06 16:24:42 Web server Listen to http://127.0.0.1:8082
```
在浏览器打开上面提示里`Web server Listen to` 的网址 `http://127.0.0.1:8082` 就可以看到网站首页