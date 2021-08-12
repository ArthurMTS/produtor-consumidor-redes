package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"encoding/gob"
	"encoding/binary"
	"strings"
	"time"
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

func sendMessage(conn net.Conn, topic string, message string) {
	buffer := new(bytes.Buffer)

	mensagem := Message{
		Topic: topic,
		Message: message,
	}

	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(mensagem)

	binary.Write(buffer, binary.BigEndian, mensagem)

	conn.Write(buffer.Bytes())
}
 
func handleMessages(conn net.Conn, topics []string) {
	defer conn.Close()
	var i = 0
	var cod = 0
	var mensagem = ""

	for {
		mensagem = fmt.Sprintf("Mensagem %d", cod)

		fmt.Printf("Sending: %s [%s]\n", mensagem, topics[i])

		sendMessage(conn, topics[i], mensagem)

		i = (i + 1) % len(topics)
		cod = cod + 1

		time.Sleep(2 * time.Second)
	}
}

func handleConnection(conn *net.TCPConn, topics []string) {

	data := Data{
		Client: 1,
		Topics: topics,
	}

	buffer := new(bytes.Buffer)

	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(data)

	binary.Write(buffer, binary.BigEndian, data)

	conn.Write(buffer.Bytes())

	handleMessages(conn, topics)
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
