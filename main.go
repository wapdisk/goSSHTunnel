package main

import (
	"360.cn/SSHTunnel/Tunnel"
	"360.cn/armory/glog"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

var (
	serverConf = ServerConf{}
)

const (
	SELF_SSHTUNNEL_CONFIGE = `.` + string(os.PathSeparator) + `SSHTunnel.conf`
)

type ServerConf struct {
	Username         string
	PublicKeyPath    string
	ServerAddrString string
	LocalAddrString  string
	RemoteAddrString string
}

// read serverconf
func ReadServerConf() {
	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := os.Open(filepath.Join(selfDir, SELF_SSHTUNNEL_CONFIGE))
	if err != nil {
		glog.Fatalln(SELF_SSHTUNNEL_CONFIGE, ` OPEN ERR:`, err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&serverConf)
	if err != nil {
		glog.Fatalln(SELF_SSHTUNNEL_CONFIGE, ` PARSE ERR:`, err)
	}
}

//配置文件输出
func TableRender() {
	data := [][]string{
		[]string{"A", "Username", serverConf.Username},
		[]string{"B", "PublicKeyPath", serverConf.PublicKeyPath},
		[]string{"C", "ServerAddrString", serverConf.ServerAddrString},
		[]string{"D", "LocalAddrString", serverConf.LocalAddrString},
		[]string{"E", "RemoteAddrString", serverConf.RemoteAddrString},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Num", "Name", "Value"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
func main() {
	// 打印出当前的版本信息:
	fmt.Println(`SSHTunnel v1.0.0`)
	// 允许go使用所有的CPU:
	runtime.GOMAXPROCS(runtime.NumCPU())
	ReadServerConf() //读取配置文件
	TableRender()    //打印配置文件
	for true {
		if serverConf.PublicKeyPath == `` { //检查是否提供了服务器密码
			fmt.Println(`Please provide the publickey for the connection:`)
			fmt.Scanln(serverConf.PublicKeyPath)
		} else {
			break
		}
	}
PRIVATEKEYS:
	privateKey, err := ioutil.ReadFile(serverConf.PublicKeyPath)
	if err != nil {
		glog.Error("id_rsa file not found; code:", err)
		timer := time.NewTimer(time.Second * 5)
		<-timer.C
		goto PRIVATEKEYS
	}
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		glog.Error("The privateKey format is not correct code:", err)
		timer := time.NewTimer(time.Second * 5)
		<-timer.C
		goto PRIVATEKEYS
	}
	config := &ssh.ClientConfig{
		User: serverConf.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	localListener := Tunnel.CreateLocalEndPoint(serverConf.LocalAddrString)                               // 创建一个本地连接:
	Tunnel.AcceptClients(localListener, config, serverConf.ServerAddrString, serverConf.RemoteAddrString) // 去偷偷用偷梁换柱大法吧:
	chExit := make(chan os.Signal, 1)                                                                     // 捕获ctrl-c,平滑退出
	signal.Notify(chExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-chExit:
	}
}
