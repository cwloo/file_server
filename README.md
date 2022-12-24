##### golang 分布式文件上传服务，用于图片，语音，视频等文件上传阿里云等 高效稳定

###### go mod download github.com/cwloo/uploader@latest

###### 1.loader 父进程，监控子进程状态并拉起，服务保活，防宕机

###### 2.http_gate 子进程，文件网关节点(多进程)

###### 3.file_server 子进程，文件上传节点(多进程)

##### file_server 启动

* $ cd loader
* $ ./loader --config=/mnt/hgfs/uploader/deploy/config/conf.ini
* `or`
* $ ./loader --dir-level=2 --conf-name=config/conf.ini

###### c 清屏指令

###### l 查看子服务

###### 55383 [file_server:1 192.168.0.113:8086 rpc:192.168.0.113:5236 file_server/ ./file_server --config=config/conf.ini --log_dir=]
###### 55384 [file_server:2 192.168.0.113:8087 rpc:192.168.0.113:5237 file_server/ ./file_server --config=config/conf.ini --log_dir=]
###### 55385 [file_server:3 192.168.0.113:8089 rpc:192.168.0.113:5238 file_server/ ./file_server --config=config/conf.ini --log_dir=]
###### 55386 [file_server:4 192.168.0.113:8090 rpc:192.168.0.113:5239 file_server/ ./file_server --config=config/conf.ini --log_dir=]
###### 55381 [http_gate:1 192.168.0.113:7787 rpc:192.168.0.113:5233 http_gate/ ./http_gate --config=config/conf.ini --log_dir=]
###### 55382 [http_gate:2 192.168.0.113:7788 rpc:192.168.0.113:5235 http_gate/ ./http_gate --config=config/conf.ini --log_dir=]

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
######             "program": "${workspaceRoot}/loader/loader.go",
######             "args": [
######                 "-dir-level=2",
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
######             "program": "${workspaceRoot}/gate/gate.go",
######             "args": [
######                 "-dir-level=2",
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
######             "program": "${workspaceRoot}/http_gate/http_gate.go",
######             "args": [
######                 "-dir-level=2",
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
######             "program": "${workspaceRoot}/file_server/file_server.go",
######             "args": [
######                 "-dir-level=2",
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

###### $ SET GOOS=linux
###### $ SET GOARCH=amd64
###### $ GOOS=linux GOARCH=amd64 go build

![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_client.png)


![image](https://github.com/cwloo/gonet/blob/master/tool/res/uploader_server.png)
