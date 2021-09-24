package Models

import ( 
	"strconv"
	"time"
	"fmt"
)

// Get laps by race ID
func GetAllLapsByRaceId(u *[]Lap, race_id_string string) (err error) {
	race_id_int, _ := strconv.Atoi (race_id_string)
	result := DB.Where("race_id = ?" , race_id_int).Order("lap_number desc").Order("race_total_time asc").Find(u)
	return result.Error
}

// Get results by race ID
func GetAllResultsByRaceId(u *[]Lap, race_id_string string) (err error) {
        race_id_int, _ := strconv.Atoi (race_id_string)
        result := DB.Where("race_id = ?" , race_id_int).Where("lap_is_current = ?" , 1).Order("lap_number desc").Order("race_total_time asc").Find(u)
        return result.Error
}

// Return current race position
func GetCurrentRacePosition(race_id uint, tag_id string) (currentRacePosition uint) {
        //var lapCopy Lap
        var lapCopy Lap
        var lapsCopy []Lap
	if  DB.Table("laps").Where("race_id = ?" , race_id).First(&lapCopy).Error == nil {
               if DB.Table("laps").Where("race_id = ?" , race_id).Where("lap_is_current = ?" , 1).Order("lap_number desc").Order("race_total_time asc").Find(&lapsCopy).Error == nil {
                        var position uint = 1
                        for _ , m := range lapsCopy {
                                if m.TagID == tag_id {
                                        break
                                } else {
                                        position = position + 1
                                }
                        }
                        currentRacePosition = position

                } else {
                        //first lap, first result
                        currentRacePosition = 1
                }
	} else {
                //no such race found (may be this is the first result).
                fmt.Println("curerntRacePosition: race data empty, is this the first result?")
                //no such race found (may be this is the first result).
                currentRacePosition = 1
	}
	return
}
// Return lap position 
func GetLapPosition(race_id uint, lap_number int, tag_id string) (lapPosition uint) {
	//var lapCopy Lap
	var lapCopy Lap
	var lapsCopy []Lap
	if  DB.Table("laps").Where("race_id = ?" , race_id).First(&lapCopy).Error == nil {
		if DB.Table("laps").Where("race_id = ?" , race_id).Where("lap_number = ?" , lap_number).Order("lap_time asc").Find(&lapsCopy).Error == nil {
			var position uint = 1
			for _ , m := range lapsCopy {
				if m.TagID == tag_id {
					break
				} else {
					position = position + 1
				}
			}
			lapPosition = position

		} else {
			//first lap, first result
			lapPosition = 1
		}
	}  else  {
		fmt.Println("lapPosition: race data empty.")
		//no such race found (may be this is the first result).
		if lap_number ==  0 {
			fmt.Println("lapPosition: This is the first result.")
			//no such race found (may be this is the first result).
			lapPosition = 1
		} else {
			//no such race found, and not a 0 lap requested - return zero.
			lapPosition = 0
		}
	}
	return
}

// Return LeaderFirstLapDiscoveryUnixTime
func GetLeaderFirstLapDiscoveryUnixTime(race_id_uint uint) (leaderFirstLapDiscoveryUnixTime int64) {
	var u Lap
	race_id_int := int(race_id_uint)
	result := DB.Table("laps").Where("race_id = ?" , race_id_int).Where("lap_number = ?" , 0).Where("lap_time = ?" , 0).First(&u)
	if result.Error != nil  {
		//fmt.Println("Error: GetLeaderFirstLapDiscoveryUnixTime", result.Error)
		leaderFirstLapDiscoveryUnixTime = 0
	} else {
		//fmt.Println("leaderFirstLapDiscoveryUnixTime =", u.DiscoveryUnixTime)	
		leaderFirstLapDiscoveryUnixTime =  u.DiscoveryUnixTime
	}
	return leaderFirstLapDiscoveryUnixTime 
}








// Return all laps in system order by date
func GetAllLaps(u *[]Lap) (err error) {

	result := DB.Order("race_id desc").Order("lap_number desc").Order("race_total_time asc").Find(u)
	return result.Error
}

// Return last lap
func GetLastLap(u *Lap) (err error) {

	result := DB.Order("discovery_unix_time desc").First(u)
	return result.Error
}

// Return last lap data by race and tag
func GetLastLapByRaceIdAndTagId(u *Lap, race_id uint, tag_id string) (err error) {
	result := DB.Where("race_id = ?" , race_id).Where("tag_id = ?" , tag_id).Order("discovery_unix_time desc").First(u)
	return result.Error
}

// Return last known lap
func GetLastRaceIDandTime(u *Lap) (lastLapRaceID uint, lastLapTime time.Time) {
	if DB.Order("discovery_unix_time desc").First(u).Error == nil {
		lastLapRaceID = u.RaceID
		lastLapTime = u.DiscoveryTimePrepared
	}
	return
}

func GetPreviousLapDataFromRaceByTagID(tagID string, raceID uint) (previousLapNumber int, previousLapTime, previousDiscoveryUnixTime, previousRaceTotalTime int64) {
	var lapStructCopy Lap
	if DB.Table("laps").Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_unix_time desc").First(&lapStructCopy).Error == nil {
		previousLapNumber = lapStructCopy.LapNumber
		previousLapTime = lapStructCopy.LapTime
		previousDiscoveryUnixTime = lapStructCopy.DiscoveryUnixTime
		previousRaceTotalTime = lapStructCopy.RaceTotalTime

	} else {
		previousLapNumber = -1
                previousLapTime = 0
                previousDiscoveryUnixTime = 0
                previousRaceTotalTime = 0
	}
	return
}

func ExpireMyPreviousLap(tagID string, raceID uint) {
	var lapStructCopy Lap
	if DB.Table("laps").Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_unix_time desc").First(&lapStructCopy).Error == nil {
		//previos lap found - update lapStructCopy.LapIsCurrent = 0
		if DB.Model(&lapStructCopy).UpdateColumn("LapIsCurrent", 0).Error != nil {
			fmt.Println("Previous lap found but update LapIsCurrent = 0 failed")
		}
	}
}

func GetMyLastLapDataFromCurrentRace(u *Lap)  (err error) {
	result := DB.Where("tag_id = ? AND race_id = ?", u.TagID, u.RaceID).Order("discovery_unix_time desc").First(u)
	return result.Error
}


// Get laps by tag ID
func GetAllLapsByTagId(u *[]Lap, tag_id string) (err error) {
	result := DB.Where("tag_id = ?" , tag_id).Order("discovery_unix_time desc").Find(u)
	return result.Error
}

func AddNewLap(u *Lap) (err error) {
	if err = DB.Create(u).Error; err != nil {
		return err
	}

	return nil
}

func GetOneLap(u *Lap, lap_id string) (err error) {
	if err := DB.Where("id = ?", lap_id).First(u).Error; err != nil {
		return err
	}

	return nil
}

func PutOneLap(u *Lap) (err error) {
	DB.Save(u)
	return nil
}


func SaveLap(u *Lap) (err error) {
	if err = DB.Save(u).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOneLap(u *Lap, lap_id string) (err error) {
	DB.Where("id = ?", lap_id).Delete(u)
	return nil
}
