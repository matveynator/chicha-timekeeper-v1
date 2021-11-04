package Config

import (
	"github.com/joho/godotenv" // Enviroment read package
	"log"
	"os"
	"flag"
	"strconv"
)

var APP_ANTENNA_LISTENER_IP,API_SERVER_LISTENER_IP,TIME_ZONE,DB_TYPE,DB_HOST,DB_USER,DB_PASSWORD,DB_NAME,DB_PORT,ADMIN_LOGIN,ADMIN_PASSWORD,PROXY_HOST,PROXY_PORT,VERSION string
var AVERAGE_RESULTS,PROXY_ACTIVE bool
var RACE_TIMEOUT_SEC,MINIMAL_LAP_TIME_SEC,RESULTS_PRECISION_SEC,LAPS_SAVE_INTERVAL_SEC int64

func readIntFromConfig(name string) (value int64) {
	val, ok := os.LookupEnv(name)
	if ok {
		value, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Printf("ERROR, %s value invalid: %s\n", name, err)
			os.Exit(1)
		} else {
			return value
		}
	} else {
		log.Printf("ERROR, configuration undefined: %s\n", name)
		os.Exit(1)
	}
	return
}

func readBoolFromConfig(name string) bool {
  var boolValue bool
	value, ok := os.LookupEnv(name)
	if ok {
		if value == "true" {
			boolValue = true
		} else if value == "false" {
			boolValue = false
		} else {
			log.Printf("ERROR, configuration not boolean: %s\n", name)
			os.Exit(1)
		}
	} else {
		log.Printf("ERROR, configuration undefined: %s\n", name)
		os.Exit(1)
	}
	return boolValue
}


func init()  {

	flagVersion := flag.Bool("version", false, "Output version information")
	flag.Parse()
	if *flagVersion  {
		if VERSION != "" {
			log.Println("Version:", VERSION)
		}
		os.Exit(0)
	}
	// Init ConfigMap here
	// Load enviroment
	// PROXY settings
	log.Println("Welcome to CHICHA, the competition timekeeper (chronograph).")
	log.Println("https://github.com/matveynator/chicha")
	if VERSION != "" {
		log.Println("Version:", VERSION)
	}
	if err := godotenv.Load("chicha.conf"); err != nil {
		log.Fatal("Configuration file chicha.conf not found", err)
	}

	// PROXY settings
	PROXY_ACTIVE = readBoolFromConfig("PROXY_ACTIVE")
	PROXY_HOST, _ = os.LookupEnv("PROXY_HOST")
	PROXY_PORT, _ = os.LookupEnv("PROXY_PORT")
	// Check enviroment
	APP_ANTENNA_LISTENER_IP, _ = os.LookupEnv("APP_ANTENNA_LISTENER_IP")
	API_SERVER_LISTENER_IP, _ = os.LookupEnv("API_SERVER_LISTENER_IP")
	TIME_ZONE, _ =  os.LookupEnv("TIME_ZONE")

	// DB connection preferences
	DB_HOST, _ = os.LookupEnv("DB_HOST")
	DB_USER, _ = os.LookupEnv("DB_USER")
	DB_PASSWORD, _ = os.LookupEnv("DB_PASSWORD")
	DB_NAME, _ = os.LookupEnv("DB_NAME")
	DB_PORT, _ = os.LookupEnv("DB_PORT")
	DB_TYPE, _ = os.LookupEnv("DB_TYPE")

	ADMIN_LOGIN, _ = os.LookupEnv("ADMIN_LOGIN")
	ADMIN_PASSWORD, _ = os.LookupEnv("ADMIN_PASSWORD")


	//INT configuration variables:
	MINIMAL_LAP_TIME_SEC = readIntFromConfig("MINIMAL_LAP_TIME_SEC")
	RACE_TIMEOUT_SEC = readIntFromConfig("RACE_TIMEOUT_SEC")
	AVERAGE_RESULTS = readBoolFromConfig("AVERAGE_RESULTS")
	RESULTS_PRECISION_SEC = readIntFromConfig("RESULTS_PRECISION_SEC")
	LAPS_SAVE_INTERVAL_SEC = readIntFromConfig("LAPS_SAVE_INTERVAL")

	log.Println("Loaded configuration from chicha.conf file.")
	return
}
