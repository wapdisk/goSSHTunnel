# goSSHTunnel
golang sshtunnel 实现免密码登录远程服务器，并映射远程端口到本地

使用方法`go get github.com/bibinbin/goSSHTunnel`

#linux和linux直接使用ssh就可以实现ssh隧道
有机器A[10.16.31.56]，B[10.16.93.204]。现想A通过ssh免密码登录到B

1.在A机下生成公钥/私钥对  `ssh-keygen -t rsa -P ''`

2.把A机下的id_rsa.pub复制到B机下`scp ~/.ssh/id_rsa.pub root@10.16.93.204:/root/`

3.B机把从A机复制的id_rsa.pub添加到.ssh/authorzied_keys文件里`cat id_rsa.pub >> .ssh/authorized_keys`

4.把B机的authorized_keys文件夹加上600权限`chmod 600 .ssh/authorized_keys`

5.实现免密码登录验证`ssh root@10.16.93.204`

端口映射 
`ssh -C -f -N -g -L 本地端口:目标IP:目标端口 用户名@目标IP`
`ssh -C -f -N -g -L 5678:10.16.93.204:9200 root@10.16.93.204`
#goSSHTunnel 使用
我们的目的是实现在windows上的 所以才有了gosshtunnel

将上面生成的id_rsa路径填写到SSHTunnel.conf中 并在里面配置你要映射端口的服务器用户名、ssh端口、要监听的端口，和你要映射到本地那个端口
```
{
    "Username": "root", 
    "PublicKeyPath": "/Users/bibinbin/.ssh/id_rsa", 
    "ServerAddrString": "10.16.93.204:22", 
    "LocalAddrString": "127.0.0.1:3690", 
    "RemoteAddrString": "127.0.0.1:9200"
}
```
使用方法`go get github.com/bibinbin/goSSHTunnel`

根据你的需要编译成不同平台下的程序
```
GO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -o gosshtunnel.exe github.com/goSSHTunnel
GO_ENABLED=0 GOOS=windows GOARCH=386  go build -o gosshtunnel.exe github.com/goSSHTunnel
GO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o gosshtunnel.exe github.com/goSSHTunnel
GO_ENABLED=0 GOOS=linux GOARCH=386  go build -o gosshtunnel.exe github.com/goSSHTunnel
```
