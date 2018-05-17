package main

import (
	"io"
	"log"
	"net"
	"fmt"
	"os"
)

func connectTwoConn(conn1 net.Conn,conn2 net.Conn){
	if conn1==nil || conn2==nil{return }
	go func() {
		io.Copy(conn1,conn2)
	}()
	go func() {
		io.Copy(conn2,conn1)
	}()
}

func server(portStr string) net.Listener{
	listener, err := net.Listen("tcp",fmt.Sprintf("0.0.0.0:%s",portStr))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
	return listener

}

func listen(listener net.Listener,connChan chan net.Conn){
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection %v\n", conn)
		connChan <- conn
	}
}


func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage %s server1:port server2:port\n", os.Args[0]);
		return
	}
	var s1ConnChan=make(chan net.Conn)
	var s2ConnChan=make(chan net.Conn)
	var s1Conn net.Conn
	var s2Conn net.Conn
	var listener1= server(os.Args[1])
	var listener2= server(os.Args[2])
	go listen(listener1,s1ConnChan)
	go listen(listener2,s2ConnChan)
	for {
		select {
		case conn := <- s1ConnChan:
			if s1Conn!=nil{
				s1Conn.Close()
			}
			s1Conn=conn
			connectTwoConn(s1Conn,s2Conn)
		case conn :=<-s2ConnChan:
			if s2Conn!=nil{
				s2Conn.Close()
			}
			s2Conn=conn
			connectTwoConn(s2Conn,s1Conn)
		}
	}
}

