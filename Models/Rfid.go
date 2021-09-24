package Models

/**
* This package module have some methods for storage RFID labels
* and store them into database
*/

import (
	"../Packages/Config"
	"../Packages/Proxy"
	"encoding/xml"
	"encoding/csv"
	"strings"
	"strconv"
	"bytes"
	"sync"
	"time"
	"log"
	"net"
	"fmt"
	"io"
)

// Buffer for new RFID requests
var laps []Lap

// Laps locker
var lapsLocker sync.Mutex

// Laps save into DB interval
var lapsSaveInterval int

// Check RFID mute timeout map
var rfidTimeoutMap map[string]time.Time

// Mute timeout duration (stored in .env)
var rfidLapMinimalTime int

// Check RFID mute timeout locker
var rfidTimeoutLocker sync.Mutex

// Start antenna listener
func StartAntennaListener() {

	if Config.PROXY_ACTIVE=="true" {
		fmt.Println("Started tcp proxy restream to", Config.PROXY_HOST,"and port:",Config.PROXY_PORT )
	}

	// Start buffer synchro with database
	go startSaveLapsBufferToDatabase()

	// Create RFID mute timeout
	rfidTimeoutMap = make(map[string]time.Time)

	// Prepare rfidLapMinimalTime
	rfidTimeout, rfidTimeoutErr := strconv.Atoi(Config.MINIMAL_LAP_TIME)
	if rfidTimeoutErr != nil {
		log.Panicln("Incorrect MINIMAL_LAP_TIME parameter in .env file", rfidTimeoutErr)
	}
	rfidLapMinimalTime = int(rfidTimeout)

	// Prepare lapsSaveInterval
	lapsInterval, lapsIntervalErr := strconv.Atoi(Config.LAPS_SAVE_INTERVAL)
	if lapsIntervalErr != nil {
		log.Panicln("Incorrect LAPS_SAVE_INTERVAL parameter in .env file", lapsIntervalErr)
	}
	lapsSaveInterval = int(lapsInterval)

	// Start listener
	l, err := net.Listen("tcp", Config.APP_ANTENNA_LISTENER_IP)
	if err != nil {
		log.Panicln("Can't start the antenna listener", err)
	}
	defer l.Close()

	// Listen new connections
	for {
		conn, err := l.Accept()
		if err != nil {
			if err != io.EOF {
				log.Panicln("tcp connection error:", err)
			}
		}

		go newAntennaConnection(conn)
	}
}

// Save laps buffer to database
func startSaveLapsBufferToDatabase() {
	for range time.Tick(time.Duration(lapsSaveInterval) * time.Second) {
		lapsLocker.Lock()
		var lapStruct Lap
		var currentlapRaceID uint 
		var currentlapLapNumber int
		lastRaceID, lastLapTime := GetLastRaceIDandTime(&lapStruct)
		if lastRaceID == 0 {
			currentlapRaceID = 1
		} else {
			raceTimeOut, _ := strconv.Atoi(Config.RACE_TIMEOUT_SEC)
			if (time.Now().UnixNano()/int64(time.Millisecond)-(int64(raceTimeOut)*1000) > lastLapTime.UnixNano()/int64(time.Millisecond)) {
				//last lap data was created more than RACE_TIMEOUT_SEC seconds ago
				//RaceID++ (create new race)
				currentlapRaceID = lastRaceID + 1

			} else {
				//last lap data was created less than RACE_TIMEOUT_SEC seconds ago
				currentlapRaceID = lastRaceID
			}
		}


		// Save laps to database
		for _,lap := range laps {
			previousLapNumber, previousLapTime, previousDiscoveryUnixTime, previousRaceTotalTime := GetPreviousLapDataFromRaceByTagID(lap.TagID, currentlapRaceID)
			if previousLapNumber != -1 {
				//set lap.LapIsCurrent = 0 for previous lap
				//set previos lap "non current"
				ExpireMyPreviousLap(lap.TagID, currentlapRaceID)
			}
			if previousLapNumber == -1 {
				currentlapLapNumber = 0
			} else {
				currentlapLapNumber = previousLapNumber + 1
			}
			//set this lap actual (current)
			lap.LapIsCurrent = 1
			lap.LapNumber = currentlapLapNumber
			lap.RaceID = currentlapRaceID
			lap.DiscoveryUnixTime = lap.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond);
			if previousLapNumber == -1 {
				//if this is first lap results:
				//#7 issue - first lap time
				leaderFirstLapDiscoveryUnixTime := GetLeaderFirstLapDiscoveryUnixTime(currentlapRaceID)
				if (leaderFirstLapDiscoveryUnixTime == 0) {
					//you are the leader set LapTime=0;
					lap.LapTime = 0
					lap.LapPosition = 1
					lap.CurrentRacePosition = 1
				} else {
					//you are not the leader of the first lap
					//calculate against the leader
					lap.LapTime = lap.DiscoveryUnixTime - leaderFirstLapDiscoveryUnixTime
					lap.LapPosition = GetLapPosition(currentlapRaceID, currentlapLapNumber, lap.TagID)
				}
			} else {
				lap.LapTime = lap.DiscoveryUnixTime - previousDiscoveryUnixTime
				lap.LapPosition = GetLapPosition(currentlapRaceID, currentlapLapNumber, lap.TagID)
			}
			lap.RaceTotalTime = previousRaceTotalTime + lap.LapTime
			if previousDiscoveryUnixTime == 0 {
				lap.BetterOrWorseLapTime = 0
			} else {
				lap.BetterOrWorseLapTime = previousLapTime-lap.LapTime
			}
			if err := AddNewLap(&lap); err != nil {
				fmt.Println("Error. Lap not added to database", err)
			} else {
			        currentRacePosition := GetCurrentRacePosition(currentlapRaceID, lap.TagID)
				DB.Model(&lap).Update("CurrentRacePosition", currentRacePosition)
			}
			fmt.Printf("Saved! tag: %s, position: %d, lap: %d, lap time: %d, total time: %d \n", lap.TagID, lap.CurrentRacePosition, lap.LapNumber, lap.LapTime, lap.RaceTotalTime)
		}


		// Clear lap buffer
		var cL []Lap
		laps = cL
		lapsLocker.Unlock()

	}
}

// Add new lap to laps buffer (private func)
func addNewLapToLapsBuffer(lap Lap) {

	// Check minimal lap time (we save only laps grater than MINIMAL_LAP_TIME from .env file)

	if expiredTime, ok := rfidTimeoutMap[lap.TagID]; !ok {

		// First time for this TagID, save lap to buffer
		lapsLocker.Lock()
		laps = append(laps, lap)
		lapsLocker.Unlock()

		// Add new value to timeouts checker map
		setNewExpriredDataForRfidTag(lap.TagID)


	} else {

		// Check previous time
		tN := time.Now()
		if tN.After(expiredTime)  {

			// Time is over, save lap to buffer
			lapsLocker.Lock()
			laps = append(laps, lap)
			lapsLocker.Unlock()


			// Generate new expired time
			setNewExpriredDataForRfidTag(lap.TagID)


		} 
	}
}

// Set new expired date for rfid Tag
func setNewExpriredDataForRfidTag(tagID string) {

	newExpiredTime := time.Now().Add(time.Duration(rfidLapMinimalTime) * time.Second)
	rfidTimeoutLocker.Lock()
	rfidTimeoutMap[tagID] = newExpiredTime
	rfidTimeoutLocker.Unlock()

}

//string to time.Unix milli
func timeFromUnixMillis(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
} 

func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// New antenna connection (private func)
func newAntennaConnection(conn net.Conn) {

	defer conn.Close()

	// Read connection in lap
	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("conn.Read(buf) error:", err)
				continue
			}
			if err == io.EOF {
				//fmt.Println("Message EOF detected - closing LAN connection.")
				break
			}
		} else {
			data := buf[:size]
			var lap Lap
			// CSV data processing
			if !IsValidXML(data) {
				//fmt.Println("Received data is not XML, trying CSV text...")
				//received data of type TEXT (parse TEXT).
				r := csv.NewReader(bytes.NewReader(data))
				r.Comma = ','
				r.FieldsPerRecord = 3
				CSV, err := r.Read()
				if err != nil {
					fmt.Println("Recived incorrect CSV data", err)
					continue
				}

				// Prepare antenna position
				antennaPosition, err := strconv.Atoi(strings.TrimSpace(CSV[2]))
				if err != nil {
					fmt.Println("Recived incorrect Antenna position CSV value:", err)
					continue
				}
				_, err =  strconv.Atoi(strings.TrimSpace(CSV[1]))
				if err != nil {
					fmt.Println("Recived incorrect discovery unix time CSV value:", err)
					continue
				} else {
					lap.DiscoveryTimePrepared, _ = timeFromUnixMillis(strings.TrimSpace(CSV[1]))
				}
				lap.TagID = strings.TrimSpace(CSV[0])
				lap.Antenna = uint8(antennaPosition)

				// XML data processing
			} else {
				// XML data processing
				// Prepare date
				//fmt.Println("Received data is valid XML")
				err := xml.Unmarshal(data, &lap)
				if err != nil {
					fmt.Println("xml.Unmarshal ERROR:", err)
					continue
				}
				//fmt.Println("TIME_ZONE=", Config.TIME_ZONE)
				loc, err := time.LoadLocation(Config.TIME_ZONE)
				if err != nil {
					fmt.Println(err)
					continue
				}
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.ParseInLocation(xmlTimeFormat, lap.DiscoveryTime, loc)
				if err != nil {
					fmt.Println("time.ParseInLocation ERROR:", err)
					continue
				}
				lap.DiscoveryTimePrepared = discoveryTime
				// Additional preparing for TagID
				lap.TagID = strings.ReplaceAll(lap.TagID, " ", "")
			}

			//Debug all received data from RFID reader
			fmt.Printf("%s, %d, %d\n", lap.TagID, lap.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond), lap.Antenna)


			if Config.PROXY_ACTIVE=="true" {
				go Proxy.ProxyDataToMotosponder(lap.TagID, lap.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond), lap.Antenna )
			}
			// Add current Lap to Laps buffer
			go addNewLapToLapsBuffer(lap)
		}
	}
}
