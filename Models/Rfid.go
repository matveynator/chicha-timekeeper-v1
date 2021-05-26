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
func StartAntennaListener(appAntennaListenerIp, rfidListenTimeoutString, lapsSaveIntervalString string) {

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

		go newAntennaConnection(conn)
	}
}

// Save laps buffer to database
func startSaveLapsBufferToDatabase() {
    for range time.Tick(time.Duration(lapsSaveInterval) * time.Second) {
        lapsLocker.Lock()

        // Save laps to database
        for _,lap := range laps {
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
func newAntennaConnection(conn net.Conn) {

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

                fmt.Println("XML error", err)

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
                    fmt.Println("Recive incorrect Antenna position value", antennaErr)
                    continue
                }

                // Prepare date
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.Parse(xmlTimeFormat, CSV[1])

                if err != nil {
                    fmt.Println("Recive incorrect time in antenna data", err)
                    continue
                }

                lap.DiscoveryTimePrepared = discoveryTime
                lap.TagID = CSV[0]
                lap.Antenna = uint8(antennaPosition)

            // XML data processing
			} else {

                // Prepare date
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.Parse(xmlTimeFormat, lap.DiscoveryTime)

                // If time is incorrect than skip them
                if err != nil {
                    continue
                }

                lap.DiscoveryTimePrepared = discoveryTime
			}

            // Additional preparing for TagID
            lap.TagID = strings.ReplaceAll(lap.TagID, " ", "")

            // Add current Lap to Laps buffer
            go addNewLapToLapsBuffer(lap)
		}
	}
}
