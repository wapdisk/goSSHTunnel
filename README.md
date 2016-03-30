# goSSHTunnel
golang sshtunnel 实现免密码登录远程服务器，并映射远程端口到本地

使用方法`go get github.com/bibinbin/goSSHTunnel`

有机器A[10.16.31.56]，B[10.16.93.204]。现想A通过ssh免密码登录到B

1.在A机下生成公钥/私钥对  `ssh-keygen -t rsa -P ''`
2.把A机下的id_rsa.pub复制到B机下`scp ~/.ssh/id_rsa.pub root@10.16.93.204:/root/`
3.B机把从A机复制的id_rsa.pub添加到.ssh/authorzied_keys文件里`cat id_rsa.pub >> .ssh/authorized_keys`
4.把B机的authorized_keys文件夹加上600权限`chmod 600 .ssh/authorized_keys`
5.实现免密码登录验证`ssh root@10.16.93.204`

端口映射 
`ssh -C -f -N -g -L 本地端口:目标IP:目标端口 用户名@目标IP`
`ssh -C -f -N -g -L 5678:10.16.93.204:9200 root@10.16.93.204`

