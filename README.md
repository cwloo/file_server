### golang 分布式文件上传服务，用于图片，语音，视频等文件上传阿里云等 高效稳定

* go mod download github.com/cwloo/uploader@latest

* 1.monitor 监控父进程，监控子进程状态并拉起，服务保活，防宕机

* 2.http_gate http文件服务网关(子进程，多进程模型)

* 3.file_server 文件上传节点(子进程，多进程模型)

* file_server 启动

* $ cd file_server/loader
* $ ./loader --config=/mnt/hgfs/uploader/deploy/config/conf.ini

* c 清屏指令

* l 查看子服务

21663 [1 ./file_server /mnt/hgfs/uploader/src/file_server/ file_server /mnt/hgfs/uploader/src/config/conf.ini ]
21664 [2 ./file_server /mnt/hgfs/uploader/src/file_server/ file_server /mnt/hgfs/uploader/src/config/conf.ini ]
21831 [3 ./file_server /mnt/hgfs/uploader/src/file_server/ file_server /mnt/hgfs/uploader/src/config/conf.ini ]
21660 [0 ./http_gate /mnt/hgfs/uploader/src/http_gate/ http_gate /mnt/hgfs/uploader/src/config/conf.ini ]
21661 [1 ./http_gate /mnt/hgfs/uploader/src/http_gate/ http_gate /mnt/hgfs/uploader/src/config/conf.ini ]
21887 [0 ./file_server /mnt/hgfs/uploader/src/file_server/ file_server /mnt/hgfs/uploader/src/config/conf.ini ]


* k 21663 kill子服务，会自动拉起

* q  killAll子服务，并退出监控

* file_client 启动

* $ cd file_client/loader
* $ ./loader

* $ SET GOOS=linux
* $ SET GOARCH=amd64
* $ GOOS=linux GOARCH=amd64 go build



![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_client.png)


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_server.png)
