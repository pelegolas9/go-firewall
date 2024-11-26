// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// func main() {
// 	if len(os.Args) != 4 {
// 		fmt.Println("Usage: script <quantity> <protocol> <interval>")
// 		return
// 	}

// 	// parsing arguments
// 	quantity, err := strconv.Atoi(os.Args[1])
// 	if err != nil || quantity <= 0 {
// 		fmt.Println("Error: quantity must be greater than zero. ")
// 		return
// 	}

// 	protocol := strings.ToLower(os.Args[2])
// 	if protocol != "tcp" && protocol != "udp" {
// 		fmt.Println("Error: protocol must be 'tcp' or 'udp'. ")
// 		return
// 	}

// 	interval, err := strconv.Atoi(os.Args[3])
// 	if err != nil || interval <= 0 {
// 		fmt.Println("Error: interval must be greater than zero. ")
// 		return
// 	}

// 	// target server details
// 	targetHost := "localhost"
// 	targetPort := "3333"

// 	fmt.Printf("Sending %d %s packets to %s:%s every %d second(s)...\n", quantity, protocol, targetHost, targetPort, interval)

// 	for i := 0; i < quantity; i++ {
// 		switch protocol {
// 		case "tcp":
// 			sendTCPPacket(targetHost, targetPort)
// 		case "udp":
// 			sendUDPPacket(targetHost, targetPort)
// 		}

// 		time.Sleep(time.Duration(interval) * time.Second)
// 	}
// }

// func sendTCPPacket(host, port string) {
// 	conn, err := net.Dial("tcp", host+":"+port)
// 	if err != nil {
// 		fmt.Println("Error connecting to TCP server:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Example packet payload
// 	payload := fmt.Sprintf("TCP Packet at %s\n", time.Now().Format(time.RFC3339))
// 	_, err = conn.Write([]byte(payload))
// 	if err != nil {
// 		fmt.Println("Error sending TCP packet:", err)
// 	} else {
// 		fmt.Println("TCP packet sent:", payload)
// 	}
// }

// func sendUDPPacket(host, port string) {
// 	conn, err := net.Dial("udp", host+":"+port)
// 	if err != nil {
// 		fmt.Println("Error connecting to UDP server:", err)
// 		return
// 	}
// 	defer conn.Close()

//		// Example packet payload
//		payload := fmt.Sprintf("UDP Packet at %s\n", time.Now().Format(time.RFC3339))
//		_, err = conn.Write([]byte(payload))
//		if err != nil {
//			fmt.Println("Error sending UDP packet:", err)
//		} else {
//			fmt.Println("UDP packet sent:", payload)
//		}
//	}
package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: script <quantity> <protocol> <interval> [<srcPort>]")
		return
	}

	// Parsing arguments
	quantity, err := strconv.Atoi(os.Args[1])
	if err != nil || quantity <= 0 {
		fmt.Println("Error: quantity must be greater than zero.")
		return
	}

	protocol := strings.ToLower(os.Args[2])
	if protocol != "tcp" && protocol != "udp" {
		fmt.Println("Error: protocol must be 'tcp' or 'udp'.")
		return
	}

	interval, err := strconv.Atoi(os.Args[3])
	if err != nil || interval <= 0 {
		fmt.Println("Error: interval must be greater than zero.")
		return
	}

	var srcPort string
	if len(os.Args) == 5 {
		srcPort = os.Args[4]
		_, err := strconv.Atoi(srcPort)
		if err != nil {
			fmt.Println("Error: srcPort must be a valid port number.")
			return
		}
	}

	// Target server details
	targetHost := "localhost"
	targetPort := "3333"

	fmt.Printf("Sending %d %s packets to %s:%s every %d second(s)", quantity, protocol, targetHost, targetPort, interval)
	if srcPort != "" {
		fmt.Printf(" from source port %s", srcPort)
	}
	fmt.Println("...")

	for i := 0; i < quantity; i++ {
		switch protocol {
		case "tcp":
			sendTCPPacket(targetHost, targetPort, srcPort)
		case "udp":
			sendUDPPacket(targetHost, targetPort, srcPort)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func sendTCPPacket(host, port, srcPort string) {
	var conn net.Conn
	var err error

	if srcPort != "" {
		// Bind to a specific source port
		localAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+srcPort)
		if err != nil {
			fmt.Println("Error resolving local TCP address:", err)
			return
		}
		dialer := &net.Dialer{LocalAddr: localAddr}
		conn, err = dialer.Dial("tcp", host+":"+port)
	} else {
		// Use any available source port
		conn, err = net.Dial("tcp", host+":"+port)
	}

	if err != nil {
		fmt.Println("Error connecting to TCP server:", err)
		return
	}
	defer conn.Close()

	// Example packet payload
	payload := fmt.Sprintf("TCP Packet at %s\n", time.Now().Format(time.RFC3339))
	_, err = conn.Write([]byte(payload))
	if err != nil {
		fmt.Println("Error sending TCP packet:", err)
	} else {
		fmt.Println("TCP packet sent:", payload)
	}
}

func sendUDPPacket(host, port, srcPort string) {
	var conn net.PacketConn
	var err error

	if srcPort != "" {
		// Bind to a specific source port
		localAddr, err := net.ResolveUDPAddr("udp", "localhost:"+srcPort)
		if err != nil {
			fmt.Println("Error resolving local UDP address:", err)
			return
		}
		conn, err = net.ListenPacket("udp", localAddr.String())
		if err != nil {
			fmt.Println("Error binding to UDP source port:", err)
			return
		}
	} else {
		// Use any available source port
		conn, err = net.ListenPacket("udp", ":0")
	}

	defer conn.Close()

	remoteAddr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		fmt.Println("Error resolving remote UDP address:", err)
		return
	}

	// Example packet payload
	payload := fmt.Sprintf("UDP Packet at %s\n", time.Now().Format(time.RFC3339))
	_, err = conn.WriteTo([]byte(payload), remoteAddr)
	if err != nil {
		fmt.Println("Error sending UDP packet:", err)
	} else {
		fmt.Println("UDP packet sent:", payload)
	}
}
