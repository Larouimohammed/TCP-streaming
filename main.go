package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Fileserver struct{}

func (fs *Fileserver) start() {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		fmt.Print("start faild")
		log.Fatal(err)
	}
	for {

		conn, err := ln.Accept()
		
		if err != nil {
			fmt.Println("conn non accepter")
			log.Fatal(err)

		}
		go fs.readloop(conn)
	}

}

func (fs *Fileserver) readloop(conn net.Conn) {
	//buf := make([]byte, 2048)
    buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn,binary.LittleEndian , &size)
		n,err :=io.CopyN(buf,conn,1000)
		//n, err := conn.Read(buf)
		if err != nil {
			fmt.Print("can't read from buf")
			log.Fatal(err)
		}

		//file := buf[:n]
		fmt.Println(buf.Bytes())
		fmt.Printf("%d byte received", n)
	}

}

func Sendfile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ":3030")
	if err != nil {
		return nil
	}
	
	binary.Write(conn,binary.LittleEndian ,int64(size))
	n,err :=io.CopyN(conn,bytes.NewReader(file),int64(size))
	
	//n, err := conn.Write(file)
	if err != nil {
		return err 
	}
	fmt.Printf("%d sended over the network", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		Sendfile(4000000)
	}()
	server := &Fileserver{}
	server.start()
}
