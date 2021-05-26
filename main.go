package main
/*
NAIKOM arena timing - free to use by general public.
based on Alien F800 RFID 912 MHZ.

Work in progress.
*/

import (
	"gorm.io/gorm" // Database ORM package
	"gorm.io/driver/postgres" // Gorm Postgres driver package
	"github.com/joho/godotenv" // Enviroment read package
	"./Models" // Our package with database models
	"flag"
	"log"
	"os"
	"net"
	"strconv"
	"encoding/xml"
	"encoding/csv"
	"strings"
	"bytes"
	"fmt"
	"time"
)

type tagXML struct {
	TagID   string     `xml:"TagID"`
	DiscoveryTime string `xml:"DiscoveryTime"`
	LastSeenTime string `xml:"LastSeenTime"`
	Antenna int     `xml:"Antenna"`
	ReadCount int   `xml:"ReadCount"`
	Protocol int `xml:"Protocol"`
}

type tagCSV struct {
	TagID string
	UnixTime string
	Antenna int
}

func main() {

	// Load enviroment
	fmt.Println("Load enviroment")
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}

	// Check enviroment
	DB_HOST, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Fatal(".env file is incorrect")
	}

	// DB connection preferences
	DB_USER, _ := os.LookupEnv("DB_USER")
	DB_PASSWORD, _ := os.LookupEnv("DB_PASSWORD")
	DB_NAME, _ := os.LookupEnv("DB_NAME")
	DB_PORT, _ := os.LookupEnv("DB_PORT")

	// Try connect to DB
	fmt.Println("Connect to DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatal("DB not opened. Check your .env settings. Status: ", err)
	} else {
		Models.DB = db
	}

	// Database Migrations
	fmt.Println("Apply migrations")
	Models.DB.AutoMigrate(&Models.Lap{})


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

			var rfid tagXML
			err := xml.Unmarshal(data, &rfid)
			if err != nil {
				//received data of type TEXT (parse TEXT).
				r := csv.NewReader(bytes.NewReader(data))
				r.Comma = ','
				r.FieldsPerRecord = 3
				CSV, err := r.Read()
				if err == nil {
					fmt.Printf("%s,%s,%s\n", CSV[0], CSV[1], CSV[2])
				}

			} else {
				//received data of type XML (parse XML)
				//2021/05/17 16:33:18.960
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.Parse(xmlTimeFormat, rfid.DiscoveryTime)
				unixMillyTime:=discoveryTime.UnixNano()/int64(time.Millisecond)
				if err == nil {
					fmt.Printf("%s, %d, %d\n", strings.ReplaceAll(rfid.TagID, " ", ""), unixMillyTime, rfid.Antenna)
				}

			}
		}

	}
}
