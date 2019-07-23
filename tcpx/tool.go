package tcpx

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// port  string
func GetTCPListener(port string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+port)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp4", tcpAddr) //监听
	if err != nil {
		return nil, err
	}
	return tcpListener, nil
}

//监听某个端口，并且返回一个chan，从chan里面可以拿到过来的客户端连接
func ListenAndGetConnsChan(port string) (chan *net.TCPConn, error) {
	c := make(chan *net.TCPConn, 0)
	listener, err := GetTCPListener(port)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if pa := recover(); pa != nil {
				log.Println(pa, "panicx")
			}
		}()

		defer listener.Close()
		defer close(c)

		for {
			tcpConn, err := listener.AcceptTCP()
			//fmt.Printf("The client:%s has connected!\n",tcpConn.RemoteAddr().String())
			if err != nil {
				log.Println(err)
				continue
			}
			c <- tcpConn
		}

	}()

	return c, nil

}

//addrport := "127.0.0.1:8282"
func GetConnByAddrPort(addrport string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addrport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ResolveTCPAddr Fatal error: %s", err.Error())
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, " DialTCP Fatal error: %s", err.Error())
		return nil, err
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 2)
	return conn, nil

}
