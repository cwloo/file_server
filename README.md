### golang 断点续传上传文件

* go mod download github.com/cwloo/uploader@latest

* 1.支持多用户上传
* 2.支持多文件批量上传
* 3.支持大文件上传
* 4.支持并发模型

* file_server         测试服务端
* file_client         测试客户端(子进程)
* file_client\loader 测试客户端(父进程)