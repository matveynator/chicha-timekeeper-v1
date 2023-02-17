package Config

import (
	"flag"
	"log"
	"fmt"
	"os"
	"hash/fnv"
)

var APP_NAME, APP_ANTENNA_LISTENER_IP, API_SERVER_LISTENER_IP, TIME_ZONE, COLLECTOR_LISTENER_ADDRESS_HASH, DB_TYPE, DB_FILE_PATH, DB_FULL_FILE_PATH, PROXY_ADDRESS, VERSION string
var AVERAGE_RESULTS bool
var RACE_TIMEOUT_SEC, MINIMAL_LAP_TIME_SEC, RESULTS_PRECISION_SEC, LAPS_SAVE_INTERVAL_SEC int

func init() {
	ParseFlags()

	log.Println("Welcome to CHICHA, the competition timekeeper (chronograph).")
	log.Println("https://github.com/matveynator/chicha")
	if VERSION != "" {
		log.Println("Version:", VERSION)
	}

}


func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func hash(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return fmt.Sprint(hash.Sum32())
}

func ParseFlags()  {
	APP_NAME = "chicha"
	flagVersion := flag.Bool("version", false, "Output version information")
	flag.StringVar(&APP_ANTENNA_LISTENER_IP, "collector", "0.0.0.0:4000", "Provide IP address and port to collect and parse data from RFID and timing readers.")
	flag.StringVar(&API_SERVER_LISTENER_IP, "web", "0.0.0.0:80", "Provide IP address and port to listen for HTTP connections from clients.")
	flag.StringVar(&PROXY_ADDRESS, "proxy", "", "Proxy received data to another collector. For example: -proxy '10.9.8.7:4000'.")
	flag.StringVar(&TIME_ZONE, "timezone", "UTC", "Set race timezone. Example: Europe/Paris, Africa/Dakar, UTC, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")
	flag.IntVar(&RACE_TIMEOUT_SEC, "timeout", 120, "Set race timeout in seconds. After this time if nobody passes the finish line the race will be stopped. ")
	flag.IntVar(&MINIMAL_LAP_TIME_SEC, "laptime", 45, "Minimal lap time in seconds. Results smaller than this duration would be considered wrong." )
	flag.IntVar(&RESULTS_PRECISION_SEC, "average-duration", 2, "Duration in seconds to calculate average results. Results passed to reader during this duration will be calculated as average result." )
	flag.IntVar(&LAPS_SAVE_INTERVAL_SEC, "save-interval", 30, "Duration in seconds to save results to database.")
	flag.BoolVar(&AVERAGE_RESULTS, "average", true, "Calculate average results instead of only first results.")

	//db
	flag.StringVar(&DB_FILE_PATH, "dbpath", ".", "Provide path to writable directory to store database data.")
	flag.StringVar(&DB_TYPE, "dbtype", "sqlite", "For now it is sqlite only")

	//process all flags
	flag.Parse()

	//делаем хеш от порта коллектора чтобы использовать в уникальном названии файла бд
	COLLECTOR_LISTENER_ADDRESS_HASH = hash(APP_ANTENNA_LISTENER_IP)

	//путь к файлу бд
	DB_FULL_FILE_PATH = fmt.Sprintf(DB_FILE_PATH+"/"+APP_NAME+"."+COLLECTOR_LISTENER_ADDRESS_HASH+".db."+DB_TYPE)

	if *flagVersion  {
		if VERSION != "" {
			fmt.Println("Version:", VERSION)
		}
		os.Exit(0)
	}
	return
}

