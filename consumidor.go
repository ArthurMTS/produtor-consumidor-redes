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
// função para verificação de erro
func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}
// função para lidar com o recebimento de mensagens do servidor
func receiveMessages(conn net.Conn, topics []string) {
	defer conn.Close()
	// cria o buffer e a mensagem
	buffer := make([]byte, 512)
	mensagem := new(Message)
	// loop para ficer recebendo mensagens do servidor
	for {
		_, err := conn.Read(buffer[0:])
		if err != nil {
			continue
		}
		// decodifica a mensagem e coloca em 'mensagem'
		tmpbuff := bytes.NewBuffer(buffer)

		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(mensagem)
		// printa a mensagem recebida e o tópico dela
		fmt.Println(mensagem.Topic, mensagem.Message)
	}
}
// função para lidar com a primeira conexão com o servidor
func handleConnection(conn net.Conn, topics []string) {
	// cria os dados que serão enviados
	data := Data{
		Client: 2,
		Topics: topics,
	}

	buffer := new(bytes.Buffer)
	// codifica os dados para enviar para o servidor
	gobobj := gob.NewEncoder(buffer)
	gobobj.Encode(data)

	binary.Write(buffer, binary.BigEndian, data)
	// envia os dados para o server
	conn.Write(buffer.Bytes())
	// chama a função para lidar com o recebimento de mensagens
	receiveMessages(conn, topics)
}

func main() {
	// espera uma string com o nome dos tópicos, separando os tópicos por vírgula
	if len(os.Args) != 2 {
		fmt.Printf("Erro! too few arguments. You should try: %s \"topic A, topic B\" \n", os.Args[0])
		os.Exit(1)
	}
	// transforma a string em um vetor
	topics := strings.Split(os.Args[1], ",")
	// cria o addr para a conexão
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)
	// realiza a conexão com o servidor
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	// vai para a função que lida com a primeira conexão com o server
	handleConnection(conn, topics)
	// fecha a conexão
	conn.Close()
	// sai do código
	os.Exit(0)
}
