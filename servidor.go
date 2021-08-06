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

type Message struct {
  Topic string
	Message string
}

var topicNameList [50]string
var topicList [50][50]string
var i = 0

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}

func addTopic(topicName string) {
	if searchTopic(topicName) != -1 {
		return
	}

	topicNameList[i] = topicName
	i++
}

func addMessage(topico string, message string) {
	index := searchTopic(topico)

	for c := 0; c < 50; c++ {
		if topicList[index][c] == "" {
			topicList[index][c] = message
			break
		}
	}
}

func searchTopic(topicName string) int {
	for c := 0; c <= i; c++ {
		if topicNameList[c] == topicName {
			return c
		}
	}

	return -1
}

func handleProducer(conn net.Conn, data *Data) {
	defer conn.Close()

	for c := 0; c < len(data.Topic); c++ {
		addTopic(data.Topic[c])
	}

	buffer := make([]byte, 512)

	_, err := conn.Read(buffer[0:])
	if err != nil {
		return
	}

	mensagem := new(Message)

	tmpbuff := bytes.NewBuffer(buffer)

	gobobj := gob.NewDecoder(tmpbuff)
	gobobj.Decode(mensagem)

	addMessage(mensagem.Topic, mensagem.Message)
}

func handleConsumer(conn net.Conn, data *Data) {
	defer conn.Close()

	index := searchTopic(data.Topic[0])

	for c := 0; c < 50; c++ {
		if (topicList[index][c]) != "" {
			conn.Write([]byte(topicList[index][c]))
			topicList[index][c] = ""
		}
	}
}

func handleClient(conn net.Conn)  {
	buffer := make([]byte, 512)

	_, err := conn.Read(buffer[0:])
	if err != nil {
		return
	}

	data := new(Data)

	tmpbuff := bytes.NewBuffer(buffer)

	gobobj := gob.NewDecoder(tmpbuff)
	gobobj.Decode(data)

	if data.Client == 1 {
		handleProducer(conn, data)
	} else if data.Client == 2 {
		handleConsumer(conn, data)
	} else {
		conn.Close()
	}
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