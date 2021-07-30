package main
/*
NAIKOM arena timing - free to use by general public.
based on Alien F800 RFID 912 MHZ.

Work in progress.
*/

import (
	"gorm.io/gorm" // Database ORM package
	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres" // Gorm Postgres driver package
	"github.com/joho/godotenv" // Enviroment read package
	"./Models" // Our package with database models
	"log"
	"os"
	"fmt"
	"strconv"
)

func main() {

	// Load enviroment
	fmt.Println("Load enviroment")
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}



	// Check enviroment
	APP_ANTENNA_LISTENER_IP, exists := os.LookupEnv("APP_ANTENNA_LISTENER_IP")
	if !exists {
		log.Fatal(".env file is incorrect")
	}
	API_SERVER_LISTENER_IP, _ := os.LookupEnv("API_SERVER_LISTENER_IP")
	TIME_ZONE, _ :=  os.LookupEnv("TIME_ZONE")
	rTS, _ := os.LookupEnv("RACE_TIMEOUT_SEC")
	rTS8, _ := strconv.Atoi(rTS)


	// DB connection preferences
	DB_HOST, _ := os.LookupEnv("DB_HOST")
	DB_USER, _ := os.LookupEnv("DB_USER")
	DB_PASSWORD, _ := os.LookupEnv("DB_PASSWORD")
	DB_NAME, _ := os.LookupEnv("DB_NAME")
	DB_PORT, _ := os.LookupEnv("DB_PORT")

	// Try connect to DB
	fmt.Println("Connect to DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}); err != nil {
		log.Fatal("DB not opened. Check your .env settings. Status: ", err)
	} else {
		Models.DB = db
	}


	// Database Migrations
	fmt.Println("Apply migrations")
	Models.DB.AutoMigrate(&Models.Lap{}, &Models.User{}, &Models.Race{}, &Models.Checkin{}, &Models.Admin{})

	// Create new system administator if them not exists
	ADMIN_LOGIN, _ := os.LookupEnv("ADMIN_LOGIN")
	ADMIN_PASSWORD, _ := os.LookupEnv("ADMIN_PASSWORD")
	Models.CreateDefaultAdmin(ADMIN_LOGIN, ADMIN_PASSWORD)

	// Start RFID listener
	fmt.Println("Started RFID data listener")
	fmt.Println("Started tcp copy stream to 192.168.96.4:4000 motosponder")
	RFID_LISTEN_TIMEOUT, _ := os.LookupEnv("RFID_LISTEN_TIMEOUT")
	LAPS_SAVE_INTERVAL, _ := os.LookupEnv("LAPS_SAVE_INTERVAL")
	go Models.StartAntennaListener(APP_ANTENNA_LISTENER_IP, RFID_LISTEN_TIMEOUT, LAPS_SAVE_INTERVAL, TIME_ZONE, int64(rTS8))



	// Routing
	r := Models.SetupRouter()
	fmt.Println("Start API server")

	// Start API server
	r.Run(API_SERVER_LISTENER_IP)
}
