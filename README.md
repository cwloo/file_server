### golang 文件服务器 断点续传上传文件

* go mod download github.com/cwloo/uploader@latest

* 1.支持多用户上传
* 2.支持多文件批量上传
* 3.支持大文件上传
* 4.支持文件上传去重
* 5.支持断点续传，客户程序下一次启动会从上次上传位置继续上传
* 6.支持并发模型

* file_server         测试服务端
* file_client         测试客户端(子进程，必须 ./loader 父进程启动)
* file_client\loader 测试客户端(父进程)

* $ SET GOOS=linux
* $ SET GOARCH=amd64
* $ GOOS=linux GOARCH=amd64 go build


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_client.png)


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_server.png)
