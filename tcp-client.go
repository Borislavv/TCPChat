package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// ClientProto	- transport protocol
// ClientPort	- default listening port
var ClientProto string = "tcp"
var ClientPort string  = ":8081"

// Client-side of app. Connecting to the server,
// sending messages and receive the answers from the server.
func main() {
	// Init user input channel
	inputCh := make(chan string)
	defer close(inputCh)

	// Connecting to the server
	connection, connectionErr := net.Dial(ClientProto, ClientPort)
	if connectionErr != nil {
		log.Println("Error occurred while connecting to the server: ", connectionErr.Error())
		return
	}

	// Printing the message, that the client is connected to the server
	fmt.Println("Connection established with the server [Proto:" + ClientProto + ", Port:" + ClientPort + "]")

	// Reading user input
	go readInput(inputCh)
	// Reading websocket message
	go readSocket(connection)

	for {
		message, ok := <-inputCh
		if !ok {
			time.Sleep(time.Second)
			continue
		}

		// Ask user a new message
		fmt.Println("You: ")

		// Writing the user input to the websocket
		_, inputWriterErr := connection.Write([]byte(message + "\n")); if inputWriterErr != nil {
			log.Println("Error occurred while writing user input to the websocket: ", inputWriterErr.Error())
			return
		}
	}
}

// readInput - reading the user input
func readInput(inputCh chan string) {
	for {
		inputText, inputErr := bufio.NewReader(os.Stdin).ReadString('\n')
		if inputErr != nil {
			log.Println("Error occurred while reading user input: ", inputErr.Error())
			continue
		}

		inputCh <- inputText
	}
}

// readSocket - listening the answer from the server
func readSocket(connection net.Conn) {
	for {
		message, readingErr := bufio.NewReader(connection).ReadString('\n')
		if readingErr != nil {
			log.Println("Error occurred while reading the answer from the server: ", readingErr.Error())
			return
		}

		// Printing the message to terminal
		fmt.Println("Companion: ", message)
	}
}