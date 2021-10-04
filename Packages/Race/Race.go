package ChichaRace

import (
    "../../Packages/Config"
    "database/sql"
    _ "github.com/lib/pq" 
    "fmt"
    "log"
    )


func ConnectDB() {


connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Config.DB_HOST, Config.DB_PORT, Config.DB_USER, Config.DB_PASSWORD, Config.DB_NAME)
		    fmt.Println(connectionString)
		    db, err := sql.Open(Config.DB_TYPE, connectionString)
		    if err != nil {
		      log.Fatal(err)
		    }
		  fmt.Println("CONNECTED OK")
		    defer db.Close()
}

var (
    race_id int
    lap_number int
    discovery_time string
    tag_id string
    )

func FetchData() {

connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Config.DB_HOST, Config.DB_PORT, Config.DB_USER, Config.DB_PASSWORD, Config.DB_NAME)


		  fmt.Println(connectionString)
		    db, err := sql.Open(Config.DB_TYPE, connectionString)
		    if err != nil {
		      log.Fatal(err)
		    }

err = db.Ping()
if err != nil {
fmt.Println("// do something here")
}


		  rows, errR := db.Query("select race_id,lap_number,discovery_time,tag_id from laps where race_id=27 order by tag_id")
		    if errR != nil {
		      log.Fatal(errR)
		    }
		  defer rows.Close()
		    for rows.Next() {
err := rows.Scan(&race_id, &lap_number, &discovery_time, &tag_id)
       if err != nil {
	 log.Fatal(err)
       }
     fmt.Println(race_id,lap_number,discovery_time,tag_id)
		    }
		  err = rows.Err()
		    if err != nil {
		      log.Fatal(err)
		    }
}
/*
// Get laps by race ID
func GetAllLapsByRaceId(u *[]LapSmall, raceid_string string) (err error) {
raceid_int, _ := strconv.Atoi (raceid_string)
result := DB.Select("race_id", "lap_number", "discovery_time", "tag_id").Where("race_id = ?" , raceid_int).Order("discovery_time asc").Find(u)
//result := DB.Where("race_id = ?" , raceid_int).Order("discovery_time asc").Find(u)
return result.Error
}
 */
