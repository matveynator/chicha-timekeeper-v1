package Models

import ( 
	"strconv"
	"time"
)

// Get laps by race ID
func GetAllLapsByRaceId(u *[]Lap, race_id_string string) (err error) {
	race_id_int, _ := strconv.Atoi (race_id_string)
	result := DB.Where("race_id = ?" , race_id_int).Order("lap_number desc").Order("race_total_time asc").Find(u)
	return result.Error
}


// Return all laps in system order by date
func GetAllLaps(u *[]Lap) (err error) {

	result := DB.Order("race_id desc").Order("lap_number desc").Order("race_total_time asc").Find(u)
	return result.Error
}

// Return all laps in system order by date
func GetLastLap(u *Lap) (err error) {

	result := DB.Order("discovery_unix_time desc").First(u)
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

func GetLastLapDataFromRaceByTagID(tagID string, raceID uint) (previousLapID uint, previousLapNumber int, previousLapTime, previousDiscoveryUnixTime, previousRaceTotalTime int64) {
	var lapStructCopy Lap
	if DB.Table("laps").Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_unix_time desc").First(&lapStructCopy).Error == nil {
		previousLapID = lapStructCopy.ID
		previousLapNumber = lapStructCopy.LapNumber
		previousLapTime = lapStructCopy.LapTime
		previousDiscoveryUnixTime = lapStructCopy.DiscoveryUnixTime
		previousRaceTotalTime = lapStructCopy.RaceTotalTime

	} else {
		previousLapID = 0
		previousLapNumber = -1
                previousLapTime = 0
                previousDiscoveryUnixTime = 0
                previousRaceTotalTime = 0
	}
	return
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

func DeleteOneLap(u *Lap, lap_id string) (err error) {
	DB.Where("id = ?", lap_id).Delete(u)
	return nil
}
