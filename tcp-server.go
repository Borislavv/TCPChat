package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// ServerProto	- transport protocol
// ServerPort	- default listening port
var ServerProto string = "tcp"
var ServerPort string  = ":8081"

// Server-side of app. Starting TCP-Websocket-Server,
// accepting the connection, getting the message and answer on it with the same text
func main() {
	startServer(ServerProto, ServerPort)
}

func startServer(serverProto string, serverPort string) {
	connectionsNumber := 0
	connections := make(map[int]net.Conn, 1024)

	fmt.Println("Starting the TCP server")

	// Start listening the port
	netListener, netListenerErr := net.Listen(serverProto, serverPort)
	if netListenerErr != nil {
		log.Println("Error occurred while listening the port: ", netListenerErr.Error())
	}

	for {
		// Accept new connection
		connection, acceptErr := netListener.Accept()
		if acceptErr != nil {
			log.Println("Error occurred while opening the port: ", acceptErr.Error())
			return
		} else {
			log.Println("Client " + connection.RemoteAddr().String() + " connected")

			// Adding each connection into pool
			//connections = append(connections, connection)
			connections[connectionsNumber] = connection
			connectionsNumber++

			go handleConnection(connection, connections)
		}
	}
}

// handleConnection - is processing the received connection
func handleConnection(connection net.Conn, connections map[int]net.Conn) {
	for {
		// Listening all messages which end on the "\n"
		message, wsReaderErr := bufio.NewReader(connection).ReadString('\n')
		if wsReaderErr != nil {
			if wsReaderErr.Error() == "EOF" {
				log.Println("Client " + connection.RemoteAddr().String() + " disconnected")
			} else {
				log.Println("Error occurred while reading the client message: ", wsReaderErr.Error())
			}

			return
		}

		// Printing the message into terminal
		fmt.Println("Client " + connection.RemoteAddr().String() + " say: ", message)

		for _, eachConnection := range connections {
			if eachConnection != connection {
				// Writing the received message back into websocket
				_, writerErr := eachConnection.Write([]byte(strings.ToUpper(message) + "\n"))
				if writerErr != nil {
					log.Println("Error occurred while writing message into websocket: ", writerErr.Error())
					return
				}
			}
		}
	}
}

