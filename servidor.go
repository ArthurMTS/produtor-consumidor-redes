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

func handleProducer(conn net.Conn, data *Data) {
	defer conn.Close()

	for c := 0; c < len(data.Topics); c++ {
		addTopic(data.Topics[c])
	}

	//fmt.Println(topicNameList)

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

		//fmt.Println(topicList[i])
		i = (i + 1) % len(data.Topics)
	}
}

func handleConsumer(conn net.Conn, data *Data) {
	defer conn.Close()

	var i = 0
	buffer := new(bytes.Buffer)

	for {
		index := searchTopic(data.Topics[i])
		if index == -1 {
			continue
		}

		for c := 0; c < 50; c++ {
			if topicList[index][c] != "" {
				mensagem := Message{
					Topic: data.Topics[i],
					Message: topicList[index][c],
				}

				topicList[index][c] = ""
			
				gobobj := gob.NewEncoder(buffer)
				gobobj.Encode(mensagem)
			
				fmt.Println(mensagem.Message, mensagem.Topic)

				binary.Write(buffer, binary.BigEndian, mensagem)

				conn.Write(buffer.Bytes())
				time.Sleep(2 * time.Second)
			}
		}
		i = (i + 1) % len(data.Topics)
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