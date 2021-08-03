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

var topicNameList [50]string
var topicList [50][50]string
var i = 0

func addTopic(topicName string) {
	for c := 0; c <= i; c++ {
		if topicNameList[c] == topicName {
			return
		}
	}

	topicNameList[i] = topicName
	i++
}

// func searchTopic(topicName string) int {
// 	for c := 0; c <= i; c++ {
// 		if topicNameList[c] == topicName {
// 			return c
// 		}
// 	}
// }

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
	fmt.Println(len(data.Topic))

	if data.Client == 1 {
		for c := 0; c < len(data.Topic); c++ {
			addTopic(data.Topic[c])
		}

		// receber mensagens do produtor

	} else if data.Client == 2 {
		// lidar com o consumidor
	}

	fmt.Println(topicNameList)
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