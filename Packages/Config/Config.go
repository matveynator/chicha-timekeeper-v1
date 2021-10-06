package Config

import (
	"fmt"
	"github.com/joho/godotenv" // Enviroment read package
	"log"
	"os"
)

var APP_ANTENNA_LISTENER_IP,API_SERVER_LISTENER_IP,TIME_ZONE,RACE_TIMEOUT_SEC,DB_TYPE,DB_HOST,DB_USER,DB_PASSWORD,DB_NAME,DB_PORT,ADMIN_LOGIN,ADMIN_PASSWORD,MINIMAL_LAP_TIME,LAPS_SAVE_INTERVAL,PROXY_ACTIVE,PROXY_HOST,PROXY_PORT,VERSION string

func init()  {
	// Init ConfigMap here
	// Load enviroment
	// PROXY settings
	fmt.Println("Welcome to CHICHA, the competition timekeeper (chronograph)!")
	fmt.Println("https://github.com/matveynator/chicha")
	fmt.Println("Version:", VERSION)
	fmt.Println("Loading chicha.conf configuration...")
	if err := godotenv.Load("chicha.conf"); err != nil {
		log.Fatal("Configuration file chicha.conf not found", err)
	}

	// PROXY settings
	PROXY_ACTIVE, _ = os.LookupEnv("PROXY_ACTIVE")
	PROXY_HOST, _ = os.LookupEnv("PROXY_HOST")
	PROXY_PORT, _ = os.LookupEnv("PROXY_PORT")
	// Check enviroment
	APP_ANTENNA_LISTENER_IP, _ = os.LookupEnv("APP_ANTENNA_LISTENER_IP")
	API_SERVER_LISTENER_IP, _ = os.LookupEnv("API_SERVER_LISTENER_IP")
	TIME_ZONE, _ =  os.LookupEnv("TIME_ZONE")
	RACE_TIMEOUT_SEC, _ = os.LookupEnv("RACE_TIMEOUT_SEC")

	// DB connection preferences
	DB_HOST, _ = os.LookupEnv("DB_HOST")
	DB_USER, _ = os.LookupEnv("DB_USER")
	DB_PASSWORD, _ = os.LookupEnv("DB_PASSWORD")
	DB_NAME, _ = os.LookupEnv("DB_NAME")
	DB_PORT, _ = os.LookupEnv("DB_PORT")
	DB_TYPE, _ = os.LookupEnv("DB_TYPE")

	ADMIN_LOGIN, _ = os.LookupEnv("ADMIN_LOGIN")
	ADMIN_PASSWORD, _ = os.LookupEnv("ADMIN_PASSWORD")

	MINIMAL_LAP_TIME, _ = os.LookupEnv("MINIMAL_LAP_TIME")
	LAPS_SAVE_INTERVAL, _ = os.LookupEnv("LAPS_SAVE_INTERVAL")

	return
}
