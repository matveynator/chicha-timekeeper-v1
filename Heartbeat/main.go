package main

import (
	"encoding/xml"
	"net"
	"fmt"
	"log"
)

//Heartbeat xml struct
type Heartbeat struct {
	IPAddress   string `xml:"DiscoveryTime"`
	IPv6Address  string `xml:"IPv6Address"`
	MACAddress string `xml:"MACAddress"`
	CommandPort string `xml:"CommandPort"`
}

// Start heartbeat listener
func StartHeartbeatListener() {


	// Start listener
	l, err := net.ListenPacket("udp", ":3988")
	if err != nil {
		log.Panicln("Can't start the Heartbeat listener", err)
	}
	defer l.Close()

	// Listen new connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go newHeartbeatConnection(conn)
	}
}


// New heartbeat connection (private func)
func newHeartbeatConnection(conn net.Conn) {

	defer conn.Close()

	// Read connection in lap
	for {
		buf := make([]byte, 8192)
		size, err := conn.Read(buf)
		if err == nil {
			data := buf[:size]
			var heartbeat Heartbeat
			err := xml.Unmarshal(data, &heartbeat)

			// CSV data processing
			if err != nil {

				fmt.Println("Received data is not XML:", err)

			} 


			//Debug all received heartbeat data from RFID reader
			fmt.Printf("%s, %s, %s, %s\n", heartbeat.IPAddress, heartbeat.IPv6Address, heartbeat.MACAddress, heartbeat.CommandPort)

		}
	}
}


func main() {
	// Start heartbeat listener
	fmt.Println("Start heartbeat data listener")
	StartHeartbeatListener()
}
