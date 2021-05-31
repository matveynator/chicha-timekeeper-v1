package Models

/**
* This package module have some methods for storage RFID labels
* and store them into database
*/

import (
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
)

// Buffer for new RFID requests
var laps []Lap

// Laps locker
var lapsLocker sync.Mutex

// Laps save into DB interval
var lapsSaveInterval int64

// Check RFID mute timeout map
var rfidTimeoutMap map[string]time.Time

// Mute timeout duration (stored in .env)
var rfidListenTimeout int64

// Check RFID mute timeout locker
var rfidTimeoutLocker sync.Mutex



// Start antenna listener
func StartAntennaListener(appAntennaListenerIp, rfidListenTimeoutString, lapsSaveIntervalString string, TIME_ZONE string) {

	// Start buffer synchro with database
	go startSaveLapsBufferToDatabase()

	// Create RFID mute timeout
	rfidTimeoutMap = make(map[string]time.Time)

	// Prepare rfidListenTimeout
	rfidTimeout, rfidTimeoutErr := strconv.Atoi(rfidListenTimeoutString)
	if rfidTimeoutErr != nil {
		log.Panicln("Incorrect RFID_LISTEN_TIMEOUT parameter in .env file")
	}
	rfidListenTimeout = int64(rfidTimeout)

	// Prepare lapsSaveInterval
	lapsInterval, lapsIntervalErr := strconv.Atoi(lapsSaveIntervalString)
	if lapsIntervalErr != nil {
		log.Panicln("Incorrect LAPS_SAVE_INTERVAL parameter in .env file")
	}
	lapsSaveInterval = int64(lapsInterval)

	// Start listener
	l, err := net.Listen("tcp", appAntennaListenerIp)
	if err != nil {
		log.Panicln("Can't start the antenna listener", err)
	}
	defer l.Close()

	// Listen new connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go newAntennaConnection(conn, TIME_ZONE)
	}
}

// Save laps buffer to database
func startSaveLapsBufferToDatabase() {
	for range time.Tick(time.Duration(lapsSaveInterval) * time.Second) {
		lapsLocker.Lock()

		var lapStruct Lap
		var currentRaceID, currentLapNumber uint
		lastRaceID, lastLapTime := GetLastRaceIDandTime(&lapStruct)
		if lastRaceID == 0 {
			currentRaceID = 1
		} else {
			currentRaceID = lastRaceID
			if (time.Now().UnixNano()/int64(time.Millisecond)-300000 > lastLapTime.UnixNano()/int64(time.Millisecond)) {
				//last lap data was created more than 300 seconds ago
				//RaceID++ (create new race)
				currentRaceID = (lastRaceID+1)

			} else {
				//last lap data was created less than 300 seconds ago
				currentRaceID = lastRaceID
			}
		}


		// Save laps to database
		for _,lap := range laps {

			lastLapNumber := GetLastLapNumberFromRaceByTagID(&lapStruct, lap.TagID, currentRaceID)
			if lastLapNumber == 0 {
				currentLapNumber = 1
			} else {
				currentLapNumber = lastLapNumber+1
			}
			lap.LapNumber=currentLapNumber
			lap.RaceID=currentRaceID

			if err := AddNewLap(&lap); err != nil {
				fmt.Println("Error. Lap not added to database")
			}
		}


		// Clear lap buffer
		var cL []Lap
		laps = cL
		lapsLocker.Unlock()

	}
}

// Add new lap to laps buffer (private func)
func addNewLapToLapsBuffer(lap Lap) {

	// Check RFID timeout (we save only first rfid data. timeout value stored in .env file as RFID_LISTEN_TIMEOUT parameter)

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

	newExpiredTime := time.Now().Add(time.Duration(rfidListenTimeout) * time.Second)
	rfidTimeoutLocker.Lock()
	rfidTimeoutMap[tagID] = newExpiredTime
	rfidTimeoutLocker.Unlock()

}

// New antenna connection (private func)
func newAntennaConnection(conn net.Conn, TIME_ZONE string) {

	defer conn.Close()

	// Read connection in lap
	for {
		buf := make([]byte, 8192)
		size, err := conn.Read(buf)
		if err == nil {
			data := buf[:size]
			var lap Lap
			err := xml.Unmarshal(data, &lap)

			// CSV data processing
			if err != nil {

				fmt.Println("Received data is not XML, trying CSV text...", err)

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
				antennaPosition, antennaErr := strconv.Atoi(CSV[2])
				if antennaErr != nil {
					fmt.Println("Recived incorrect Antenna position value:", antennaErr)
					continue
				}

				// Prepare date
				loc, _ := time.LoadLocation(TIME_ZONE)
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.ParseInLocation(xmlTimeFormat, CSV[1], loc)

				if err != nil {
					fmt.Println("Recived incorrect time from RFID reader:", err)
					continue
				}

				lap.DiscoveryTimePrepared = discoveryTime
				lap.TagID = CSV[0]
				lap.Antenna = uint8(antennaPosition)

				// XML data processing
			} else {

				// Prepare date
				loc, _ := time.LoadLocation(TIME_ZONE)
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.ParseInLocation(xmlTimeFormat, lap.DiscoveryTime, loc)



				//unixMillyTime:=discoveryTime.UnixNano()/int64(time.Millisecond)
				// If time is incorrect than skip them
				if err != nil {
					continue
				}

				lap.DiscoveryTimePrepared = discoveryTime
			}

			// Additional preparing for TagID
			lap.TagID = strings.ReplaceAll(lap.TagID, " ", "")

			//Debug all received data from RFID reader
			fmt.Printf("%s, %d, %d\n", lap.TagID, lap.DiscoveryTimePrepared, lap.Antenna)

			// Add current Lap to Laps buffer
			go addNewLapToLapsBuffer(lap)
		}
	}
}
