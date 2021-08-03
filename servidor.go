package main

import (
	"fmt"
	"net"
	"os"
	"bytes"
	"encoding/gob"
)

type Data struct {
  Client int
	Topic []string
}

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn)  {
	defer conn.Close()

	buffer := make([]byte, 512)

	_, err := conn.Read(buffer[0:])
	if err != nil {
		return
	}

	data := new(Data)

	tmpbuff := bytes.NewBuffer(buffer)

	gobobj := gob.NewDecoder(tmpbuff)
	gobobj.Decode(data)

	fmt.Println("Cliente: ", data.Client)
	fmt.Println("Tópico: ", data.Topic[0])
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}