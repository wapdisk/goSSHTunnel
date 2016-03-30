package Tunnel

import (
	"360.cn/armory/glog"
	"golang.org/x/crypto/ssh"
	"io"
	//"log"
	"net"
	"time"
)

const (
	maxRetriesLocal  = 10000000 // 允许创建多少个本地终端？
	maxRetriesRemote = 10000000 // 允许创建多少个远程服务终端？
	maxRetriesServer = 10000000 // 允许创建多少个ssh服务？
)

var (
	currentRetriesLocal = 0 // 检查有多少个本地端口监听
)

func AcceptClients(connection net.Listener, config *ssh.ClientConfig, serverAddrString, remoteAddrString string) {
	// 死死死死死死循环，直到天荒地久
	for {
		// 接受的本地端口监听。。。。。听说过凌静门吗？？？:
		if localConn, err := connection.Accept(); err != nil {
			glog.Errorln("Accepting a client failed:", err.Error())
		} else {
			//看我的偷梁换柱。。。。。。。。
			go forward(localConn, config, serverAddrString, remoteAddrString)
		}
	}
}

func forward(localConn net.Conn, config *ssh.ClientConfig, serverAddrString, remoteAddrString string) {
	defer localConn.Close()
	currentRetriesServer := 0
	currentRetriesRemote := 0
	var sshClientConnection *ssh.Client = nil
	// 不断的循环重试
	for {
		// 尝试通过ssh connect 远程服务器
		if sshClientConn, err := ssh.Dial(`tcp`, serverAddrString, config); err != nil {
			currentRetriesServer++
			glog.Errorln("Was not able to connect with the SSH server ", serverAddrString, ":", err.Error())
			// 如果连接数小于最大连接数。。。等待1s
			if currentRetriesServer < maxRetriesServer {
				glog.Infoln(`Retry...`)
				time.Sleep(1 * time.Second)
			} else {
				//过多则直接返回
				glog.Infoln(`No more retries for connecting the SSH server.`)
				return
			}
		} else {
			// Success:
			glog.Infoln(`Connected to the SSH server ` + serverAddrString)
			sshClientConnection = sshClientConn
			defer sshClientConnection.Close()
			break
		}
	}
	for {
		if sshConn, err := sshClientConnection.Dial(`tcp`, remoteAddrString); err != nil {
			currentRetriesRemote++
			glog.Errorln("Was not able to connect with the SSH server ", serverAddrString, ":", err.Error())
			if currentRetriesRemote < maxRetriesRemote {
				glog.Infoln(`Retry...`)
				time.Sleep(1 * time.Second)
			} else {
				glog.Errorln(`No more retries for connecting the SSH server.`)
				return
			}
		} else {
			//端口已经连上
			glog.Infof("The remote end-point %s is connected.\n", remoteAddrString)
			defer sshConn.Close()
			quit := make(chan bool)
			// 建立双向转移（两个新的线程为此创建）:
			go transfer(localConn, sshConn, `Local => Remote`, quit)
			go transfer(sshConn, localConn, `Remote => Local`, quit)
			isRunning := true
			for isRunning {
				select {
				case <-quit:
					glog.Infoln(`At least one transfer was stopped.`)
					isRunning = false
					break
				}
			}
			glog.Infoln(`Close now all connections.`)
			return
		}
	}
}

// 偷梁换柱大法
func transfer(fromReader io.Reader, toWriter io.Writer, name string, quit chan bool) {
	glog.Infof("%s transfer started.", name)
	if _, err := io.Copy(toWriter, fromReader); err != nil {
		glog.Errorln(name, "transfer failed: \n", err.Error())
	} else {
		glog.Infof("%s transfer closed.\n", name)
	}
	quit <- true
}

//创建监听本地端口
func CreateLocalEndPoint(localAddrString string) (localListener net.Listener) {
	for {
		if localListenerObj, err := net.Listen(`tcp`, localAddrString); err != nil {
			currentRetriesLocal++
			glog.Errorf("Was not able to create the local end-point %s: %s\n", localAddrString, err.Error())
			if currentRetriesLocal < maxRetriesLocal {
				glog.Infoln(`Retry...`)
				time.Sleep(1 * time.Second)
			} else {
				glog.Errorln(`No more retries for the local end-point: ` + localAddrString)
			}
		} else {
			glog.Infoln(`Listen to local address ` + localAddrString)
			localListener = localListenerObj
			return
		}
	}
}
