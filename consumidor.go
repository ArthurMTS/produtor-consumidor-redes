package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"encoding/gob"
	"encoding/binary"
	"strings"
)

type Data struct {
  Client int
	Topics []string
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

func receiveMessages(conn net.Conn, topics []string) {
	defer conn.Close()

	buffer := make([]byte, 512)
	mensagem := new(Message)

	for {
		_, err := conn.Read(buffer[0:])
		if err != nil {
			continue
		}

		tmpbuff := bytes.NewBuffer(buffer)

		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(mensagem)

		fmt.Println(mensagem.Topic, mensagem.Message)
	}
}

func handleConnection(conn net.Conn, topics []string) {

	data := Data{
		Client: 2,
		Topics: topics,
	}

	buffer := new(bytes.Buffer)

	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(data)

	binary.Write(buffer, binary.BigEndian, data)

	conn.Write(buffer.Bytes())

	receiveMessages(conn, topics)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Erro! too few arguments. You should try: %s \"topic A, topic B\" \n", os.Args[0])
		os.Exit(1)
	}

	topics := strings.Split(os.Args[1], ",")

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)
	
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	handleConnection(conn, topics)

	conn.Close()

	os.Exit(0)
}
