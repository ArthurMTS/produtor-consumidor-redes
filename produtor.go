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
// função para verificação de erro
func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}
// função para enviar uma mensagem do produtor ao servidor
func sendMessage(conn net.Conn, topic string, message string) {
	buffer := new(bytes.Buffer)
	// cria a mensagem
	mensagem := Message{
		Topic: topic,
		Message: message,
	}
	// codifica a mensagem
	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(mensagem)

	binary.Write(buffer, binary.BigEndian, mensagem)
	// envia a mensagem
	conn.Write(buffer.Bytes())
}
// função para lidar com o envio das mensagens
func handleMessages(conn net.Conn, topics []string) {
	defer conn.Close()
	// variáveis auxiliares
	var i = 0
	var cod = 0
	var mensagem = ""
	// loop para ficar enviando mensagens ao servidor
	for {
		// espera por 2 segundos
		time.Sleep(2 * time.Second)
		// faz uma mensagem qualquer
		mensagem = fmt.Sprintf("Mensagem %d", cod)
		// printa a mensagem que será enviada
		fmt.Printf("Sending: %s [%s]\n", mensagem, topics[i])
		// chama a função que lida com o envio
		sendMessage(conn, topics[i], mensagem)
		// para mandar mensagens de diferentes tópicos
		i = (i + 1) % len(topics)
		cod = cod + 1
	}
}
// função para realizar a primeira conexão com o servidor
func handleConnection(conn *net.TCPConn, topics []string) {
	// cria os dados que serão enviados
	data := Data{
		Client: 1,
		Topics: topics,
	}

	buffer := new(bytes.Buffer)
	// codifica para enviar
	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(data)

	binary.Write(buffer, binary.BigEndian, data)
	// envia
	conn.Write(buffer.Bytes())
	// parte para o envio de mensagens
	handleMessages(conn, topics)
}

func main() {
	// espera uma string com o nome dos tópicos, separando os tópicos por vírgula
	if len(os.Args) != 2 {
		fmt.Printf("Erro! too few arguments. You should try: %s \"topic A, topic B\" \n", os.Args[0])
		os.Exit(1)
	}
	// transforam a string em um vetor
	topics := strings.Split(os.Args[1], ",")
	// cria o addr para a conexão
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)
	// realiza a conexão
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	// vai para a função que lida com a primeira conexão
	handleConnection(conn, topics)
	// fecha a conexão
	conn.Close()
  // fecha o programa
	os.Exit(0)
}
