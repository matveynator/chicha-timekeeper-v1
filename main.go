package main
/*
NAIKOM arena timing - free to use by general public.
based on Alien F800 RFID 912 MHZ.

Work in progress. 
*/

import (
	"flag"
	"log"
	"net"
	"strconv"
	"encoding/xml"
	"strings"
	"fmt"
	"time"
)

type AlienRFIDTag struct {
	TagID   string     `xml:"TagID"`
	DiscoveryTime string `xml:"DiscoveryTime"`
	LastSeenTime string `xml:"LastSeenTime"`
	Antenna int     `xml:"Antenna"`
	ReadCount int   `xml:"ReadCount"`
	Protocol int `xml:"Protocol"`
}


func main() {

	port := flag.Int("port", 4000, "Port to accept connections on.")
	host := flag.String("host", "0.0.0.0", "Host or IP to bind to")
	flag.Parse()

	l, err := net.Listen("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Listening to connections at '"+*host+"' on port", strconv.Itoa(*port))
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("Accepted new connection.")
	defer conn.Close()
	defer fmt.Println("Closed connection.")

	for {
		buf := make([]byte, 8192)
		size, err := conn.Read(buf)
		if err == nil {
			data := buf[:size]

			//echo raw data for debug:
			//fmt.Printf("%s\n", data);

			var rfid AlienRFIDTag
			err := xml.Unmarshal(data, &rfid)
			if err != nil {
				//received data of type TEXT (parse TEXT).
				fmt.Printf("%s\n", data)

			} else {
				//received data of type XML (parse XML)
				//2021/05/17 16:33:18.960
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.Parse(xmlTimeFormat, rfid.DiscoveryTime)
				unixMillyTime:=discoveryTime.UnixNano()/int64(time.Millisecond)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Printf("%s, %d, %d\n", strings.ReplaceAll(rfid.TagID, " ", ""), unixMillyTime, rfid.Antenna)
			}
		}

	}
}

