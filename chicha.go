package main

/*
Work in progress.
*/

import (
	"chicha/Packages/Config"
	"chicha/Packages/View"
	"chicha/Packages/Models" // Our package with database models
	"embed"
	//"fmt"
	"log"

	//"gorm.io/driver/postgres" // Gorm Postgres driver package
	//"gorm.io/driver/sqlite"   //gorm sqlite driver
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"            // Database ORM package
	"gorm.io/gorm/logger"

	//"github.com/sethvargo/go-password/password" //password generator
	//profiling CPU:
	//"github.com/pkg/profile"

)

//go:embed static
var static embed.FS

func main() {
	//profiling CPU: https://hackernoon.com/go-the-complete-guide-to-profiling-your-code-h51r3waz
	//defer profile.Start(profile.ProfilePath(".")).Stop()

		//DEFAULT: if Config.DB_TYPE == "sqlite"
		dsn := "chicha.sqlite"
		if db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			SkipDefaultTransaction: true,
		}); err != nil {
			log.Println("ERROR: Connect to local SQLite database failed at", dsn, err)
			panic(err)

		} else {
			Models.DB = db
			log.Println("Connected to local SQLite database at", dsn)
		}

	// Database Migrations
	log.Println("Creating or changing database structures (applying migrations)...")
	Models.DB.AutoMigrate(&Models.Lap{}, &Models.User{}, &Models.Race{}, &Models.Checkin{}, &Models.Admin{})

	// Create new system administator if them not exists

	//adminPass, err := password.Generate(8, 1, 3, true, true)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("Creating system administrator account if not exists with name =", Config.ADMIN_LOGIN, "and password =", adminPass)
	//	Models.CreateDefaultAdmin(Config.ADMIN_LOGIN, adminPass)
	//}

	updCh := make(chan struct{}, 1)
	Models.UpdateChan = updCh


	// Start RFID listener
	go Models.StartAntennaListener()
	log.Printf("Data collector IP = %s, db save interval = %d sec, minimal lap time = %d sec.\n", Config.APP_ANTENNA_LISTENER_IP, Config.LAPS_SAVE_INTERVAL_SEC, Config.MINIMAL_LAP_TIME_SEC)

	// Routing
	r := Models.SetupRouter()
	// view
	view.New(r, static, updCh)

	// Start API server
	log.Printf("WEB API server IP = %s\n", Config.API_SERVER_LISTENER_IP)
	errAPI := r.Run(Config.API_SERVER_LISTENER_IP)
	if errAPI != nil {
		log.Println("ERROR: API server start failed:", errAPI)
	}
}
