package Models

/**
* This package module have some methods for storage RFID labels
* and store them into database
*/

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"sort"

	"chicha/Packages/Config"
	"chicha/Packages/Proxy"
)

// Buffer for new RFID requests
var laps []Lap

// Laps locker channgel
var lapsChannelLocker = make(chan int, 1)

// Check RFID mute timeout map
var rfidTimeoutMap map[string]time.Time

// Check RFID mute timeout locker
var rfidTimeoutLocker sync.Mutex

// Start antenna listener
func StartAntennaListener() {

	if Config.PROXY_ACTIVE == "true" {
		log.Println("Started tcp proxy restream to", Config.PROXY_HOST, "and port:", Config.PROXY_PORT)
	}

	// Start buffer synchro with database
	lapsChannelLocker <- 0 //Put the initial value into the channel to unlock operations
	//go saveLapsBufferToDB()

	// Create RFID mute timeout
	rfidTimeoutMap = make(map[string]time.Time)

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
func saveLapsBufferToDB() {
	for range time.Tick(time.Duration(Config.LAPS_SAVE_INTERVAL_SEC) * time.Second) {
		//<-lapsChannelLocker //grab the ticket via channel (lock others)
		var currentlapRaceID uint
		var currentlapLapNumber int

		// Save laps to database
		for _, lap := range laps {
			previousLapNumber, previousDiscoveryUnixTime, previousRaceTotalTime := GetPreviousLapDataFromRaceByTagID(lap.TagID, currentlapRaceID)
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
			lap.DiscoveryUnixTime = lap.DiscoveryTimePrepared.UnixNano() / int64(time.Millisecond)
			if previousLapNumber == -1 {
				//if this is first lap results:
				//#7 issue - first lap time
				leaderFirstLapDiscoveryUnixTime, err := GetLeaderFirstLapDiscoveryUnixTime(currentlapRaceID)
				if err == nil {
					//you are not the leader of the first lap
					//calculate against the leader
					lap.LapTime = lap.DiscoveryUnixTime - leaderFirstLapDiscoveryUnixTime
					lap.LapPosition = GetLapPosition(currentlapRaceID, currentlapLapNumber, lap.TagID)
				} else {
					//you are the leader set LapTime=0;
					lap.LapTime = 0
					lap.LapPosition = 1
					lap.CurrentRacePosition = 1
				}
			} else {
				lap.LapTime = lap.DiscoveryUnixTime - previousDiscoveryUnixTime
				lap.LapPosition = GetLapPosition(currentlapRaceID, currentlapLapNumber, lap.TagID)
			}

			//race total time
			lap.RaceTotalTime = previousRaceTotalTime + lap.LapTime
			//log.Println("race total time:", lap.RaceTotalTime, "lap time", lap.LapTime)

			leaderRaceTotalTime := GetLeaderRaceTotalTimeByRaceIdAndLapNumber(lap.RaceID, lap.LapNumber)
			if leaderRaceTotalTime == 0 {
				//first lap
				//log.Println("leaderRaceTotalTime = 0 - first lap detected, TimeBehindTheLeader = lap.LapTime:", lap.LapTime)
				if lap.LapPosition == 1 {
					lap.TimeBehindTheLeader = 0
				} else {
					lap.TimeBehindTheLeader = lap.LapTime
				}
			} else {
				lap.TimeBehindTheLeader = lap.RaceTotalTime - leaderRaceTotalTime
			}

			//START: лучшее время и возможные пропуски в учете на воротах RFID (lap.LapIsStrange):
			if lap.LapNumber == 0 {
				//едем нулевой круг
				lap.BestLapTime = lap.LapTime
				lap.BetterOrWorseLapTime = 0
				_, err := GetBestLapTimeFromRace(currentlapRaceID)
				if err == nil {
					//если кто то проехал уже 2 круга а мы едем только нулевой
					//не нормально - помечаем что круг странный (возможно не считалась метка)
					lap.LapIsStrange = 1
				} else {
					//нормально - еще нет проехавших второй круг
					lap.LapIsStrange = 0
				}
			} else if lap.LapNumber == 1 {
				//едем первый полный круг
				lap.BestLapTime = lap.LapTime
				lap.BetterOrWorseLapTime = 0
				//узнаем лучшее время круга у других участников:
				currentRaceBestLapTime, _ := GetBestLapTimeFromRace(currentlapRaceID)
				lapIsStrange := int(math.Round(float64(lap.LapTime) / float64(currentRaceBestLapTime)))
				if lapIsStrange >= 2 {
					//если наше время в 2 или более раз долльше лучего времени этого круга у других участников
					//отметим что круг странный (возможно не считалась метка)
					lap.LapIsStrange = 1
				} else {
					//нормально - наше время не очень долгое (вероятно правильно считалось)
					lap.LapIsStrange = 0
				}
			} else {
				//едем второй полный круг и все последующие
				//запросим свое предыдущее лучшее время круга:
				myPreviousBestLapTime, _ := GetBestLapTimeFromRaceByTagID(lap.TagID, currentlapRaceID)
				if lap.LapTime > myPreviousBestLapTime {
					lap.BestLapTime = myPreviousBestLapTime
				} else {
					lap.BestLapTime = lap.LapTime
				}
				//улучшил или ухудшил свое предыдущее лучшее время?
				lap.BetterOrWorseLapTime = lap.LapTime - myPreviousBestLapTime
				lapIsStrange := int(math.Round(float64(lap.LapTime) / float64(lap.BestLapTime)))
				if lapIsStrange >= 2 {
					//если наше время в 2 и более раз дольше чем наше лучшее время круга
					//отметим что круг странный (метка возможно просто не считалась)
					lap.LapIsStrange = 1
				} else {
					lap.LapIsStrange = 0
				}
			}
			//END: лучшее время и возможные пропуски в учете на воротах RFID (lap.LapIsStrange):


			errL := DB.Where("id = ?", lap.ID).First(&lap).Error
			if errL != nil {
				DB.Create(&lap)
			} 

			err := DB.Save(&lap).Error

			if err != nil {
				log.Println("Error. Lap not added to database", err)
			} else {
				log.Printf("Saved! tag: %s, lap: %d, lap time: %d, total time: %d \n", lap.TagID, lap.LapNumber, lap.LapTime, lap.RaceTotalTime)
				spErr := UpdateCurrentStartPositionsByRaceId(currentlapRaceID)
				if spErr != nil {
					log.Println("UpdateCurrentStartPositionsByRaceId(currentlapRaceID) Error", spErr)
				}
				upErr := UpdateCurrentResultsByRaceId(currentlapRaceID)
				if upErr != nil {
					log.Println("UpdateCurrentResultsByRaceId(currentlapRaceID) Error", upErr)
				}

				//refresh my results
				golERR := DB.Where("id = ?", lap.ID).First(&lap).Error
				if golERR == nil {
					if lap.CurrentRacePosition == 1 {
						//if I am the leader - update other riders results - set lap.StageFinished=0
						err := UpdateAllStageNotYetFinishedByRaceId(currentlapRaceID)
						if err != nil {
							log.Println("UpdateAllStageNotYetFinishedByRaceId(currentlapRaceID) ERROR:", err)
						}
					}

					//update that your lap is finished lap.StageFinished=1 in any case
					lap.StageFinished = 1

					//save final results
					sfErr := DB.Save(&lap).Error
					if sfErr != nil {
						log.Println("lap.StageFinished=1 Error. Lap not added to database", sfErr)
					} else {
						err := PrintCurrentResultsByRaceId(currentlapRaceID)
						if err != nil {
							log.Println("PrintCurrentResultsByRaceId(currentlapRaceID) ERROR:", err)
						}
					}
				} else {
					log.Println("GetOneLap(&lap) ERROR:", golERR)
				}
			}
		}

		// Clear lap buffer
		//var cL []Lap
		//laps = cL
		//lapsChannelLocker <- 1 //give ticket back via channel (unlock operations)
	}
}


func getMyLastLapFromBuffer(newLap Lap) (myLastLap Lap) {
	//block 1: get my previous results from this race - start block.
	var myLastLaps []Lap
	//gather all my laps from previous results:
	for _, savedLap := range laps {
		if savedLap.TagID == newLap.TagID {
			myLastLaps = append(myLastLaps, savedLap)
		}
	}

	if len(myLastLaps) == 1 {
		//allready have one lap
		//get my last result:
		myLastLap = myLastLaps[0]

	} else if len(myLastLaps) > 1 {
		//allready have more than one lap
		sort.Slice(myLastLaps, func(i, j int) bool {
			//sort ascending by DisoveryUnixTime
			return myLastLaps[i].DiscoveryUnixTime < myLastLaps[j].DiscoveryUnixTime
		})
		//get my last result (newest DisoveryUnixTime result)
		myLastLap = myLastLaps[len(myLastLaps)-1]
	}
	return
	//block 1: get my previous results from this race - finish block.
}

func getLastLapFromBuffer() (lastLap Lap) {
	//block 1: get previous results from this race - start block.
	if len(laps) == 1 {
		//allready have one lap
		//get my last result:
		lastLap = laps[0]

	} else if len(laps) > 1 {
		//allready have more than one lap
		sort.Slice(laps, func(i, j int) bool {
			//sort ascending by DisoveryUnixTime
			return laps[i].DiscoveryUnixTime < laps[j].DiscoveryUnixTime
		})
		//get my last result (newest DisoveryUnixTime result)
		lastLap = laps[len(laps)-1]
	}
	return
	//block 1: get previous results from this race - finish block.
}


// Add new lap to laps buffer (private func)
func addNewLapToLapsBuffer(newLap Lap) {
<-lapsChannelLocker //grab the ticket via channel (lock)
	newLap.DiscoveryUnixTime = newLap.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond)

	if len(laps) == 0 {
		//empty create race and lap
		fmt.Println("Slice empty - adding new element with TagID = ", newLap.TagID)
		newLap.LapNumber=0;
		newLap.LapPosition=1;
		newLap.RaceID=1;
		newLap.CurrentRacePosition=1;
		newLap.DiscoveryAverageUnixTime = newLap.DiscoveryUnixTime
		laps = append(laps, newLap)
		log.Printf("SAVED %d TO BUFFER: laps: %d, raceid: %d, tag: %s\n", newLap.LapNumber, len(laps), newLap.RaceID,  newLap.TagID )
	} else {
		//get any previous lap data:
		lastLap := getLastLapFromBuffer()
		if lastLap != (Lap{}) {
			//lastLap not empty
			//get my previous lap data:
			gap := newLap.DiscoveryUnixTime - lastLap.DiscoveryUnixTime
			myLastLap := getMyLastLapFromBuffer(newLap)
			//my last lap not empty
			if myLastLap != (Lap{}) {
			  myGap := newLap.DiscoveryUnixTime - myLastLap.DiscoveryUnixTime
				fmt.Printf("gap: %d, myGap: %d \n", gap, myGap)

				if  myGap >= 0 && Config.RESULTS_PRECISION_SEC*1000 >= myGap  {
					//from 0 to 5 sec (RESULTS_PRECISION_SEC) = update DiscoveryAverageUnixTime data
					myLastLap.DiscoveryAverageUnixTime = (myLastLap.DiscoveryAverageUnixTime + newLap.DiscoveryUnixTime) / 2
					log.Printf("UPDATED BUFFER: laps: %d, raceid: %d, lap#: %d, avtime: %d, tag: %s\n", len(laps), myLastLap.RaceID, myLastLap.LapNumber, myLastLap.DiscoveryAverageUnixTime, myLastLap.TagID )
				} else if Config.RESULTS_PRECISION_SEC*1000 < myGap && myGap < Config.MINIMAL_LAP_TIME_SEC*1000 {
					//from 5 to 30 sec (RESULTS_PRECISION_SEC - MINIMAL_LAP_TIME_SEC) = discard data - ERROR DATA RECEIVED!
					log.Println("ERROR DATA RECEIVED: from 5 to 30 sec", newLap.TagID)
				} else if Config.MINIMAL_LAP_TIME_SEC*1000 <= myGap && gap < Config.RACE_TIMEOUT_SEC*1000 {
					//from 30 to 300 sec (MINIMAL_LAP_TIME_SEC - RACE_TIMEOUT_SEC) passed  = create new lap LapNumber++! 
					newLap.LapNumber = myLastLap.LapNumber + 1
					newLap.DiscoveryAverageUnixTime = newLap.DiscoveryUnixTime
					laps = append(laps, newLap)
					log.Printf("ADDED NEXT LAP TO BUFFER: laps: %d, raceid: %d, lap#: %d, tag: %s\n", len(laps), newLap.RaceID, newLap.LapNumber, newLap.TagID )
				} else if gap > Config.RACE_TIMEOUT_SEC*1000 {
					//> 300 sec (RACE_TIMEOUT_SEC) passed  = create new Race and First Lap: RaceID=lastLap.RaceID+1, LapNumber=0
					//but first - clean previous race data
					newLap.RaceID=lastLap.RaceID+1
					newLap.LapNumber=0;
					newLap.LapPosition=1;
					newLap.CurrentRacePosition=1;
					newLap.DiscoveryAverageUnixTime = newLap.DiscoveryUnixTime

					// Clear lap buffer, start with clean slice:
					var cL []Lap
					laps = cL
					laps = append(laps, newLap)
					log.Printf("SAVED NEXT RACE TO BUFFER: laps: %d, raceid: %d, lap#: %d, tag: %s\n", len(laps), newLap.RaceID, newLap.LapNumber, newLap.TagID )
				} else {

					log.Printf("STRANGE!: laps: %d, raceid: %d, lap#: %d, tag: %s\n", len(laps), newLap.RaceID, newLap.LapNumber, newLap.TagID )
				}
			} else {
				//no my results - create my first lap in same race
				newLap.LapNumber=0;
				newLap.RaceID=lastLap.RaceID;
				newLap.DiscoveryAverageUnixTime = newLap.DiscoveryUnixTime
				laps = append(laps, newLap)
				log.Printf("SAVED NEW PLAYER TO BUFFER: laps: %d, raceid: %d, lap#: %d, tag: %s\n", len(laps), newLap.RaceID, newLap.LapNumber, newLap.TagID )
			}
		} else {
			//SOME ERROR - lastLap EMPTY ?
			log.Println("SOME ERROR - lastLap EMPTY ?", lastLap, laps,)
		}
	}
	lapsChannelLocker <- 1 //give ticket back via channel (unlock)
}
// Set new expired date for rfid Tag
func setNewExpriredDataForRfidTag(tagID string) {

	newExpiredTime := time.Now().Add(time.Duration(Config.MINIMAL_LAP_TIME_SEC) * time.Second)
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
	var tempDelay time.Duration // how long to sleep on accept failure

	// Read connection in lap
	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				//log.Println("conn.Read(buf) error:", err)
				//log.Println("Message EOF detected - closing LAN connection.")
				break
			}

			if ne, ok := err.(*net.OpError); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			break
		}
		tempDelay = 0

		data := buf[:size]
		var lap Lap
		lap.AntennaIP = fmt.Sprintf("%s", conn.RemoteAddr().(*net.TCPAddr).IP)

		//various data formats processing (text csv, xml) start:
		if !IsValidXML(data) {
			// CSV data processing
			//log.Println("Received data is not XML, trying CSV text...")
			//received data of type TEXT (parse TEXT).
			r := csv.NewReader(bytes.NewReader(data))
			r.Comma = ','
			r.FieldsPerRecord = 3
			CSV, err := r.Read()
			if err != nil {
				log.Println("Recived incorrect CSV data", err)
				continue
			}

			// Prepare antenna position
			antennaPosition, err := strconv.ParseInt(strings.TrimSpace(CSV[2]), 10, 64)
			if err != nil {
				log.Println("Recived incorrect Antenna position CSV value:", err)
				continue
			}
			_, err = strconv.ParseInt(strings.TrimSpace(CSV[1]), 10, 64)
			if err != nil {
				log.Println("Recived incorrect discovery unix time CSV value:", err)
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
			//log.Println("Received data is valid XML")
			err := xml.Unmarshal(data, &lap)
			if err != nil {
				log.Println("xml.Unmarshal ERROR:", err)
				continue
			}
			//log.Println("TIME_ZONE=", Config.TIME_ZONE)
			loc, err := time.LoadLocation(Config.TIME_ZONE)
			if err != nil {
				log.Println(err)
				continue
			}
			xmlTimeFormat := `2006/01/02 15:04:05.000`
			discoveryTime, err := time.ParseInLocation(xmlTimeFormat, lap.DiscoveryTime, loc)
			if err != nil {
				log.Println("time.ParseInLocation ERROR:", err)
				continue
			}
			lap.DiscoveryTimePrepared = discoveryTime
			// Additional preparing for TagID
			lap.TagID = strings.ReplaceAll(lap.TagID, " ", "")
		}
		//various data formats processing (text csv, xml) end.

		//Debug all received data from RFID reader
		log.Printf("NEW: IP=%s, TAG=%s, TIME=%d, ANT=%d\n", lap.AntennaIP, lap.TagID, lap.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond), lap.Antenna)

		if Config.PROXY_ACTIVE == "true" {
			go Proxy.ProxyDataToMotosponder(lap.TagID, lap.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond), lap.Antenna)
		}
		// Add current Lap to Laps buffer
		go addNewLapToLapsBuffer(lap)
	}
}
