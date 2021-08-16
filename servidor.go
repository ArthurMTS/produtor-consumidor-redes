package main

import (
	"fmt"
	"net"
	"os"
	"bytes"
	"encoding/gob"
	"encoding/binary"
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

var topicNameList [50]string
var topicList [50][50]string

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}

func searchTopic(topic string) int {
	for c := 0; c < len(topicNameList); c++ {
		if topic == topicNameList[c] {
			return c
		}
	}

	return -1
}

func addTopic(topic string) {
	if searchTopic(topic) != -1 {
		return
	}

	for c := 0; c < len(topicNameList); c++ {
		if topicNameList[c] == "" {
			topicNameList[c] = topic
			break
		}
	}
}

func addMessage(topic string, message string) {
	index := searchTopic(topic)

	for c := 0; c < len(topicList[index]); c++ {
		if topicList[index][c] == "" {
			topicList[index][c] = message
			break
		}
	}
}

func handleProducer(conn net.Conn, topics []string) {
	defer conn.Close()

	for c := 0; c < len(topics); c++ {
		addTopic(topics[c])
	}

	fmt.Println(topicNameList)

	buffer := make([]byte, 512)
	mensagem := new(Message)

	var i = 0

	for {
		_, err := conn.Read(buffer[0:])
		if err != nil {
			continue
		}

		tmpbuff := bytes.NewBuffer(buffer)

		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(mensagem)

		addMessage(mensagem.Topic, mensagem.Message)

		fmt.Println(topicList[i])
		i = (i + 1) % len(topics)
	}
}

func handleConsumer(conn net.Conn, topics []string) {
	defer conn.Close()

	var i = 0

	for {
		index := searchTopic(topics[0])
		if index == -1 {
			continue
		}

		for c := 0; c < 50; c++ {
			if topicList[index][c] != "" {
				time.Sleep(2 * time.Second)

				buffer := new(bytes.Buffer)

				mensagem := Message{
					Topic: topics[0],
					Message: topicList[index][c],
				}

				gobobj := gob.NewEncoder(buffer)
				gobobj.Encode(mensagem)

				binary.Write(buffer, binary.BigEndian, mensagem)

				fmt.Println(mensagem.Message, mensagem.Topic)

				conn.Write(buffer.Bytes())

				topicList[index][c] = ""
			}
		}
		i = (i + 1) % len(topics)
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
		handleProducer(conn, data.Topics)
	} else if data.Client == 2 {
		handleConsumer(conn, data.Topics)
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