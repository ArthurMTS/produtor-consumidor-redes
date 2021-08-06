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

type Message struct {
  Topic string
	Message string
}

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}
 
func handleMessages(conn net.Conn) {
	defer conn.Close()

	buffer := new(bytes.Buffer)

	mensagem := Message{
		Topic: "Tópico A",
		Message: "Eu sou o dougras!",
	}

	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(mensagem)

	binary.Write(buffer, binary.BigEndian, mensagem)

	conn.Write(buffer.Bytes())
}

func handleConnection(conn *net.TCPConn) {

	data := Data{
		Client: 1,
		Topic: []string{"Tópico A"},
	}

	buffer := new(bytes.Buffer)

	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(data)

	binary.Write(buffer, binary.BigEndian, data)

	conn.Write(buffer.Bytes())

	handleMessages(conn)
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
