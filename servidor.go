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

// array para relacionar o nome de um tópico a um indice
var topicNameList [50]string
// matriz para armazenar as diferentes mensagens de cada tópico
var topicList [50][50]string

// função para verificação de erro
func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s\n", err.Error())
		os.Exit(1)
	}
}

// função responsável por verificar se um tópico existe
// se sim ela retorna seu indice, se não retorna -1
func searchTopic(topic string) int {
	for c := 0; c < len(topicNameList); c++ {
		if topic == topicNameList[c] {
			return c
		}
	}

	return -1
}

// função para adicionar um tópico ao array de nomes de tópicos
func addTopic(topic string) {
	// verificação para não colocar o mesmo tópico duas vezes
	if searchTopic(topic) != -1 {
		return
	}

	// caso o tópico não exista, percorre o array e coloca o tópico na primeira posição disponível
	for c := 0; c < len(topicNameList); c++ {
		if topicNameList[c] == "" {
			topicNameList[c] = topic
			break
		}
	}
}

// função para adicioanr uma mensagem a matriz em seu tópico correto
func addMessage(topic string, message string) {
	// pega o indice do tópico ao qual a mensagem deve ser adicionada
	index := searchTopic(topic)

	// coloca o mensagem no indice do tópico em uma posição disponível
	for c := 0; c < len(topicList[index]); c++ {
		if topicList[index][c] == "" {
			topicList[index][c] = message
			break
		}
	}
}

// função para lidar com as funcionalidade do cliente produtor
func handleProducer(conn net.Conn, topics []string) {
	defer conn.Close()

	// pegando todos os tópicos que o produtor mandou e adicionando-os
	// obs: se o tópico já existir ele não adiciona
	for c := 0; c < len(topics); c++ {
		addTopic(topics[c])
	}

	// printando os tópicos existentes atualmente
	fmt.Println(topicNameList)

	// criando um buffer e uma nova mensagem
	buffer := make([]byte, 512)
	mensagem := new(Message)

	var i = 0

	// loop para ficar recebendo mensagens do produtor
	for {
		_, err := conn.Read(buffer[0:])
		if err != nil {
			continue
		}

		// tudo isso é para decodificar a mensagem
		tmpbuff := bytes.NewBuffer(buffer)

		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(mensagem)

		// com a mensagem decodificada ela é adicionada ao seu respectivo tópico
		addMessage(mensagem.Topic, mensagem.Message)

		// printando as mensagens dos tópicos
		fmt.Println(topicList[i])
		i = (i + 1) % len(topics)
	}
}

// função para lidar com as funcionalidade do cliente consumidor
func handleConsumer(conn net.Conn, topics []string) {
	defer conn.Close()

	var i = 0

	// loop para ficar enviando as mensagens para os consumidores
	for {
		// pegando o indice da lista do tópico que ele vai enviar
		index := searchTopic(topics[i])
		if index == -1 {
			continue
		}

		// percorre a lista para enviar todas as mensagens nela para o consumidor
		for c := 0; c < 50; c++ {
			// se tiver uma mensagem
			if topicList[index][c] != "" {
				// espera 2 segundos
				time.Sleep(2 * time.Second)

				// cria um buffer para enviar a mensagem
				buffer := new(bytes.Buffer)
				// cria a mensagem com o tópico dela e o conteúdo
				mensagem := Message{
					Topic: topics[i],
					Message: topicList[index][c],
				}
				// codifica a mensagem para enviar
				gobobj := gob.NewEncoder(buffer)
				gobobj.Encode(mensagem)

				binary.Write(buffer, binary.BigEndian, mensagem)
				// printa a mensagem que irá enviar
				fmt.Println(mensagem.Message, mensagem.Topic)
				// envia a mensagem
				conn.Write(buffer.Bytes())
				// tira a mensagem da lista
				topicList[index][c] = ""
			}
		}
		// isso é para trocar o tópico
		i = (i + 1) % len(topics)
	}
}

// função para receber a primeira conexão do cliente
func handleClient(conn net.Conn)  {
	buffer := make([]byte, 512)

	_, err := conn.Read(buffer[0:])
	if err != nil {
		return
	}
	// cria a mensagem e decodifica ela
	data := new(Data)

	tmpbuff := bytes.NewBuffer(buffer)

	gobobj := gob.NewDecoder(tmpbuff)
	gobobj.Decode(data)
	// se a mensagem tiver o código 1 é produtor
	// se tiver o código 2 é consumidor e se por acaso não tiver código ele fecha a conexão
	if data.Client == 1 {
		// chama a função para lidar com o produtor
		handleProducer(conn, data.Topics)
	} else if data.Client == 2 {
		// chama a função para lidar com o consumidor
		handleConsumer(conn, data.Topics)
	} else {
		conn.Close()
	}
}

func main() {
	// criando o addr
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)
	// criando o ouvinte
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	// loop para ficar recebendo conexões dos clientes
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		
		go handleClient(conn)
	}
}