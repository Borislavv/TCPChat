package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// ClientProto	- transport protocol
// ClientPort	- default listening port
var ClientProto string = "tcp"
var ClientPort string  = ":8081"

// Client-side of app. Connecting to the server,
// sending messages and receive the answers from the server.
func main() {
	// Connecting to the server
	connection, connectionErr := net.Dial(ClientProto, ClientPort)
	if connectionErr != nil {
		log.Println("Error occurred while connecting to the server: ", connectionErr.Error())
		return
	}

	// Printing the message, that the client is connected to the server
	fmt.Println("Connection established with the server [Proto:" + ClientProto + ", Port:" + ClientPort + "]")

	for {
		// Awaiting while user type the new input
		inputReader := bufio.NewReader(os.Stdin)

		// Ask user a new message
		fmt.Println("Your message:")

		// Reading the user input
		inputText, inputErr := inputReader.ReadString('\n')
		if inputErr != nil {
			log.Println("Error occurred while reading the user input: ", inputErr.Error())
			return
		}

		// Writing the user input to the websocket
		_, inputWriterErr := connection.Write([]byte(inputText + "\n")); if inputWriterErr != nil {
			log.Println("Error occurred while writing user input to the websocket: ", inputWriterErr.Error())
			return
		}

		// Listening the answer from the server
		message, readingErr := bufio.NewReader(connection).ReadString('\n')
		if readingErr != nil {
			log.Println("Error occurred while reading the answer from the server: ", readingErr.Error())
			return
		}

		// Printing the message to terminal
		fmt.Println("Message from the server: ", message)
	}
}