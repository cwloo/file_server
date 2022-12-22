##### golang 分布式文件上传服务，用于图片，语音，视频等文件上传阿里云等 高效稳定

###### go mod download github.com/cwloo/uploader@latest

###### 1.monitor 监控父进程，监控子进程状态并拉起，服务保活，防宕机

###### 2.http_gate http文件服务网关(子进程，多进程模型)

###### 3.file_server 文件上传节点(子进程，多进程模型)

* $ SET GOOS=linux
* $ SET GOARCH=amd64
* $ GOOS=linux GOARCH=amd64 go build

##### file_server 启动

* $ cd file_server/loader
* $ ./loader --config=uploader/deploy/config/conf.ini

###### c 清屏指令

###### l 查看子服务

###### 55383 [file:0 192.168.0.113:8086 rpc:127.0.0.1:5236 uploader/src/file_server/ ./file_server --config=uploader/src/config/conf.ini --log_dir=]
###### 55384 [file:1 192.168.0.113:8087 rpc:127.0.0.1:5237 uploader/src/file_server/ ./file_server --config=uploader/src/config/conf.ini --log_dir=]
###### 55385 [file:2 192.168.0.113:8089 rpc:127.0.0.1:5238 uploader/src/file_server/ ./file_server --config=uploader/src/config/conf.ini --log_dir=]
###### 55386 [file:3 192.168.0.113:8090 rpc:127.0.0.1:5239 uploader/src/file_server/ ./file_server --config=uploader/src/config/conf.ini --log_dir=]
###### 55381 [gate.http:0 192.168.0.113:7787 rpc:127.0.0.1:5233 uploader/src/http_gate/ ./http_gate --config=uploader/src/config/conf.ini --log_dir=]
###### 55382 [gate.http:1 192.168.0.113:7788 rpc:127.0.0.1:5235 uploader/src/http_gate/ ./http_gate --config=uploader/src/config/conf.ini --log_dir=]


###### k pid kill子服务，会自动拉起

###### q  killAll子服务，并退出监控

##### launch.json debug 调试

###### {
######     "version": "0.2.0",
######     "configurations": [
######         {
######             "name": "loader",
######             "type": "go",
######             "request": "launch",
######             "mode": "debug",
######             "program": "${workspaceRoot}/src/loader/loader.go",
######             "args": [
######                 "-c=${workspaceRoot}/deploy/config/conf.ini",
######                 "-l=${workspaceRoot}/deploy/log",
######                 "-server=192.168.0.103:8000",
######                 "-rpc=192.168.0.103:5228",
######             ],
######         },
######         {
######             "name": "gate",
######             "type": "go",
######             "request": "launch",
######             "mode": "debug",
######             "program": "${workspaceRoot}/src/gate/gate.go",
######             "args": [
######                 "-c=${workspaceRoot}/deploy/config/conf.ini",
######                 "-l=${workspaceRoot}/deploy/log",
######                 "-server=192.168.0.103:7786",
######                 "-rpc=192.168.0.103:5232",
######             ],
######         },
######         {
######             "name": "http_gate",
######             "type": "go",
######             "request": "launch",
######             "mode": "debug",
######             "program": "${workspaceRoot}/src/http_gate/http_gate.go",
######             "args": [
######                 "-c=${workspaceRoot}/deploy/config/conf.ini",
######                 "-l=${workspaceRoot}/deploy/log",
######                 "-server=192.168.0.103:7788",
######                 "-rpc=192.168.0.103:5235",
######             ],
######         },
######         {
######             "name": "file_server",
######             "type": "go",
######             "request": "launch",
######             "mode": "debug",
######             "program": "${workspaceRoot}/src/file_server/file_server.go",
######             "args": [
######                 "-c=${workspaceRoot}/deploy/config/conf.ini",
######                 "-l=${workspaceRoot}/deploy/log",
######                 "-server=192.168.0.103:8086",
######                 "-rpc=192.168.0.103:5236",
######             ],
######         }
######     ]
###### }



##### file_client 启动

* $ cd file_client/loader
* $ ./loader


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_client.png)


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_server.png)
