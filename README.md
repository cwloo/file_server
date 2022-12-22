### golang 分布式文件上传服务，用于图片，语音，视频等文件上传阿里云等 高效稳定

* go mod download github.com/cwloo/uploader@latest

* 1.monitor 监控父进程，监控子进程状态并拉起，服务保活，防宕机

* 2.http_gate http文件服务网关(子进程，多进程模型)

* 3.file_server 文件上传节点(子进程，多进程模型)

* file_server         测试服务端(必须 ./loader 父进程启动)

* file_server 服务端启动

* $ cd file_server/loader
* $ ./loader --config=/mnt/hgfs/uploader/deploy/config/conf.ini

* file_client 启动

* $ cd file_client/loader
* $ ./loader

* $ SET GOOS=linux
* $ SET GOARCH=amd64
* $ GOOS=linux GOARCH=amd64 go build


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_client.png)


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_server.png)
