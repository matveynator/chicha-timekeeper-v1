package main
/*
NAIKOM arena timing - free to use by general public.
based on Alien F800 RFID 912 MHZ.

Work in progress.
*/

import (
	"./Packages/Config"
	"gorm.io/gorm" // Database ORM package
	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres" // Gorm Postgres driver package
	"github.com/sethvargo/go-password/password" //password generator
	"./Models" // Our package with database models
	"fmt"
)


func main() {

	//Database section
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", Config.DB_HOST, Config.DB_USER, Config.DB_PASSWORD, Config.DB_NAME, Config.DB_PORT)
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		fmt.Println("ERROR: Connect to database failed at", Config.DB_HOST, Config.DB_PORT, "with database name =", Config.DB_NAME, "and user =", Config.DB_USER, err)
	} else {
		Models.DB = db
		fmt.Println("Connected to database at", Config.DB_HOST, Config.DB_PORT, "with database name =", Config.DB_NAME, "and user =", Config.DB_USER)
	}


	// Database Migrations
	fmt.Println("Creating or changing database structures (applying migrations)...")
	Models.DB.AutoMigrate(&Models.Lap{}, &Models.User{}, &Models.Race{}, &Models.Checkin{}, &Models.Admin{})

	// Create new system administator if them not exists

	adminPass, err := password.Generate(8, 1, 3, true, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Creating system administrator account if not exists with name =", Config.ADMIN_LOGIN, "and password =", adminPass)
		Models.CreateDefaultAdmin(Config.ADMIN_LOGIN, adminPass)
	}

	// Start RFID listener
	go Models.StartAntennaListener(Config.APP_ANTENNA_LISTENER_IP, Config.MINIMAL_LAP_TIME, Config.LAPS_SAVE_INTERVAL, Config.TIME_ZONE, int64(Config.RTS8))
	fmt.Println("Started RFID data listener at", Config.APP_ANTENNA_LISTENER_IP, "with laps save interval =", Config.LAPS_SAVE_INTERVAL, "and lap minimal duration =", Config.MINIMAL_LAP_TIME, "seconds")

	// Routing
	r := Models.SetupRouter()

	// Start API server
	err = r.Run(Config.API_SERVER_LISTENER_IP)

	if err != nil {
		fmt.Println("ERROR: API server start failed:", err)
	} else {
		fmt.Println("OK: Started API server at", Config.API_SERVER_LISTENER_IP)
	}
}
