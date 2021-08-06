package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"encoding/gob"
	"encoding/binary"
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

func reveiveMessage(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)

	for {
		n, err := conn.Read(buffer[0:])
		if err != nil {
			return
		}

		mensagem := string(buffer[0:n])

		fmt.Println(mensagem)
	}
}

func handleConnection(conn *net.TCPConn) {

	data := Data{
		Client: 2,
		Topic: []string{"Tópico A"},
	}

	buffer := new(bytes.Buffer)

	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(data)

	binary.Write(buffer, binary.BigEndian, data)

	conn.Write(buffer.Bytes())

	reveiveMessage(conn)
}

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)
	
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	handleConnection(conn)

	conn.Close()

	os.Exit(0)
}
