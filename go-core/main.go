/*
how do we monitor the network?
we need to listen for incoming request on a host and port

when a request is sent, break it down to check:
- source ip
- destination ip
- port
- protocol
- payload?

then, we check the rules for it on rules.go file


we managed to accept connections, now we need to get the data sent and parse it
let's create a packet structure to hold the incoming packets

we did a simple CLI behavior, now we want the rules to be persistent.
let's save the rules that are created in a file and read from it
steps:

when a rules is added, it opens the file rules.txt
scans the file to check if this rule already exists
if exists, tell the user
if not, write the rule to the file and close it

writing is ok, now we need to check the added rules to the file
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// packet structure
type Packet struct {
	PROTOCOL    string
	SOURCE_IP   string
	SOURCE_PORT string
}

func main() {
	// ip and port variables
	var HOST = "localhost"
	var PORT = "3333"

	fmt.Println("Starting firewall...")
	fmt.Println("Type 'help' for a list of commands.")

	// goroutine to start monitoring the network
	go func() {
		err := monitorNetwork(HOST, PORT)
		if err != nil {
			fmt.Println("Error monitoring network: ", err)
		}
	}()

	// starter variables for input scanning
	input_scanner := bufio.NewScanner(os.Stdin)

	// for loop unitl user exits
	for {
		fmt.Print("firewall> ")
		if input_scanner.Scan() {
			input := input_scanner.Text()
			getUserCommand(input)
		}
	}
}

func getUserCommand(userInput string) {

	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		fmt.Println("No command entered. Type 'help' for available commands.")
		return
	}

	args := strings.Split(userInput, " ")
	command := args[0]

	switch command {
	case "addrule":
		if len(args) < 5 {
			fmt.Println("Usage: addrule <protocol> <source_ip> <source_port> <allow|block>")
			return
		}
		protocol := args[1]
		sourceIP := args[2]
		sourcePort := args[3]
		allow := strings.ToLower(args[4]) == "allow"

		rule := Rule{
			Protocol: protocol,
			SourceIP: sourceIP,
			Port:     sourcePort,
			Allow:    allow,
		}

		// check if the rule passed alreay exists
		_, err := CheckFileRules(rule)
		if err != nil {
			fmt.Println("Rule already exists! ")
			return
		}

		// add rule
		err = AddRule(rule)
		if err != nil {
			fmt.Println("Error adding rule: ", err)
		} else {
			fmt.Printf("Rule added: %+v\n", rule)
		}

	case "exit":
		fmt.Println("Exiting...")
		os.Exit(0)

	case "help":
		fmt.Println("Commands:")
		fmt.Println("  addrule <protocol> <src_ip> <port> <allow|block> - Add a firewall rule")
		fmt.Println("  exit                                 - Exit the firewall")
	default:
		fmt.Println("Unknow command. Type 'help' for a list of commands")
	}
}

func monitorNetwork(host, port string) error {
	// first, listen for request on host:port, ie address
	address := host + ":" + port

	// now we setup a listener with TCP protocol
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start listener on %s: %v", address, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s...\n", address)

	// for loop to monitor requests
	for {
		// we accept connections on the address
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection ", err)
			continue
		}

		// connection details
		srcIP, srcPort := splitAddress(conn.RemoteAddr().String())

		// separating the host and port to get only the IP
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error parsing source address: ", err)
			return nil
		}

		if host == "localhost" {
			host = "127.0.0.1"
		}

		packet := Packet{
			PROTOCOL:    "tcp",
			SOURCE_IP:   srcIP,
			SOURCE_PORT: srcPort,
		}

		// check packet against the rules
		action := CheckRules(packet)

		if action {
			fmt.Println("Packet allowed: ", packet)
		} else {
			fmt.Println("Packet blocked: ", packet)
		}

		conn.Close()

	}

}

// handleConnection only parses the packets for now
// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	fmt.Printf("Received connection from %s\n", conn.RemoteAddr().String())

// 	// getting connection data and details
// 	srcAddr := conn.RemoteAddr().String()
// 	dstAddr := conn.LocalAddr().String()

// 	// splitting the ip and port for source and destination addresses using the custom splitter
// 	srcIP, srcPort := splitAddress(srcAddr)
// 	dstIP, dstPort := splitAddress(dstAddr)

// 	// creating a buffer to read incoming data
// 	buffer := bufio.NewReader(conn)
// 	data, err := buffer.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Error reading data: ", err)
// 		return
// 	}

// 	packet := Packet{
// 		SrcIP:   srcIP,
// 		DstIP:   dstIP,
// 		SrcPort: srcPort,
// 		DstPort: dstPort,
// 		Data:    strings.TrimSpace(data),
// 	}

// 	// logging packet details
// 	fmt.Printf("Packet received: %+v\n", packet)

// }

func splitAddress(addr string) (string, string) {
	// split "ip:port"
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return addr, ""
	}

	return parts[0], parts[1]
}
