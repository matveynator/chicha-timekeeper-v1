package Models

import ( 
	"strconv"
	"time"
	"../Packages/Race"
)

// Get laps by race ID
func GetAllLapsByRaceId(u *[]LapSmall, raceid_string string) (err error) {
	ChichaRace.FetchData()
	raceid_int, _ := strconv.Atoi (raceid_string)
	result := DB.Select("race_id", "lap_number", "discovery_time", "tag_id").Where("race_id = ?" , raceid_int).Order("discovery_time asc").Find(u)
	//result := DB.Where("race_id = ?" , raceid_int).Order("discovery_time asc").Find(u)
	return result.Error
}

// Return all laps in system order by date
func GetAllLaps(u *[]Lap) (err error) {

	result := DB.Order("discovery_time desc").Find(u)
	return result.Error
}

// Return all laps in system order by date
func GetLastLap(u *Lap) (err error) {

	result := DB.Order("discovery_time desc").First(u)
	return result.Error
}

// Return last known lap
func GetLastRaceIDandTime(u *Lap) (lastLapRaceID uint, lastLapTime time.Time) {
	if DB.Order("discovery_time desc").First(u).Error == nil {
		lastLapRaceID = u.RaceID
		lastLapTime = u.DiscoveryTimePrepared
	}
	return
}

func GetLastLapDataFromRaceByTagID(tagID string, raceID uint) (lastlapLapNumber uint, lastlapLapTime, lastlapDiscoveryUnixTime, lastlapRaceTotalTime int64) {
	var lapStructCopy Lap
	if DB.Table("laps").Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_time desc").First(&lapStructCopy).Error == nil {
		lastlapLapNumber = lapStructCopy.LapNumber
		lastlapLapTime = lapStructCopy.LapTime
		lastlapDiscoveryUnixTime = lapStructCopy.DiscoveryUnixTime
		lastlapRaceTotalTime = lapStructCopy.RaceTotalTime

	} else {
		lastlapLapNumber = 0
                lastlapLapTime = 0
                lastlapDiscoveryUnixTime = 0
                lastlapRaceTotalTime = 0
	}
	return
}

func GetMyLastLapDataFromCurrentRace(u *Lap)  (err error) {
	result := DB.Where("tag_id = ? AND race_id = ?", u.TagID, u.RaceID).Order("discovery_time desc").First(u)
	return result.Error
}


// Get laps by tag ID
func GetAllLapsByTagId(u *[]Lap, tag_id string) (err error) {
	result := DB.Where("tag_id = ?" , tag_id).Order("discovery_time desc").Find(u)
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
