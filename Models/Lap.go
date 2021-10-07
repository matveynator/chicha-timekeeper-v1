package Models

import (
	"fmt"
	"strconv"
	"time"
)

// Get laps by race ID
func GetAllLapsByRaceId(u *[]Lap, race_id_string string) (err error) {
	race_id_int, _ := strconv.ParseInt(race_id_string, 10, 64)
	result := DB.Where("race_id = ?", race_id_int).Order("lap_number desc").Order("race_total_time asc").Find(u)
	return result.Error
}

// Get results by race ID
func GetAllResultsByRaceId(u *[]Lap, race_id_string string) (err error) {
	race_id_int, _ := strconv.ParseInt(race_id_string, 10, 64)
	result := DB.Where("race_id = ?", race_id_int).Where("lap_is_current = ?", 1).Order("lap_number desc").Order("race_total_time asc").Find(u)
	return result.Error
}

// Update current start positions by race ID
func UpdateAllStageNotYetFinishedByRaceId(race_id uint) (err error) {
	var laps []Lap
	err = DB.Where("race_id = ?", race_id).Where("lap_is_current = ?", 1).Find(&laps).Error
	if err == nil {
		for _, lap := range laps {
			lap.StageFinished = 0
			err := DB.Save(lap).Error
			if err != nil {
				fmt.Println("UpdateAllStageNotYetFinishedByRaceId DB.Save(lap) Error:", err)
			}
		}
	}
	return
}

// Update current start positions by race ID
func UpdateCurrentStartPositionsByRaceId(race_id uint) (err error) {
	var laps []Lap
	err = DB.Where("race_id = ?", race_id).Where("lap_number != ?", 0).Where("lap_is_current = ?", 1).Order("best_lap_time asc").Find(&laps).Error
	if err == nil {
		var position uint = 1
		for _, lap := range laps {
			lap.BestLapPosition = position
			err := DB.Save(lap).Error
			if err != nil {
				fmt.Println("UpdateCurrentStartPositionsByRaceId DB.Save(lap) Error:", err)
			}
			position = position + 1
		}
	}
	return
}

// Update current results by race ID
func UpdateCurrentResultsByRaceId(race_id uint) (err error) {
	var laps []Lap
	err = DB.Where("race_id = ?", race_id).Where("lap_is_current = ?", 1).Order("lap_number desc").Order("race_total_time asc").Find(&laps).Error
	if err == nil {
		var position uint = 1
		var LeaderDiscoveryUnixTime int64
		for _, lap := range laps {
			if position == 1 {
				LeaderDiscoveryUnixTime = lap.DiscoveryUnixTime
			}
			if lap.TimeBehindTheLeader == 0 {
				//update previous leader results
				lap.TimeBehindTheLeader = LeaderDiscoveryUnixTime - lap.DiscoveryUnixTime
			}
			lap.CurrentRacePosition = position
			err := DB.Save(lap).Error
			if err != nil {
				fmt.Println("UpdateCurrentResultsByRaceId Error:", err)
			}
			position = position + 1
		}
	}
	return
}

//Print Current Results
func PrintCurrentResultsByRaceId(race_id uint) (err error) {
	var laps []Lap
	err = DB.Where("race_id = ?", race_id).Where("lap_is_current = ?", 1).Order("lap_number desc").Order("race_total_time asc").Find(&laps).Error
	if err == nil {
		for _, lap := range laps {
			fmt.Printf("lap: %d, tag: %s, position: %d, start#: %d, time: %d, gap: %d, best lap: %d, finish?: %d, strange?: %d\n", lap.LapNumber, lap.TagID, lap.CurrentRacePosition, lap.BestLapPosition, lap.RaceTotalTime, lap.TimeBehindTheLeader, lap.BestLapTime, lap.StageFinished, lap.LapIsStrange)
		}
	}
	return
}

// Return leader race_total_time by race_id
func GetLeaderRaceTotalTimeByRaceIdAndLapNumber(race_id uint, lap_number int) (leaderRaceTotalTime int64) {
	var lap Lap
	if DB.Where("race_id = ?", race_id).Where("lap_number = ?", lap_number).Where("lap_position = ?", 1).First(&lap).Error == nil {
		//fmt.Println("lap found - time:", lap.RaceTotalTime)
		leaderRaceTotalTime = lap.RaceTotalTime
	} else {
		//fmt.Println("lap not found", race_id, lap_number)
		leaderRaceTotalTime = 0
	}
	return
}

// Return current race position
func GetCurrentRacePosition(race_id uint, tag_id string) (currentRacePosition uint) {
	//var lapCopy Lap
	var lapCopy Lap
	var lapsCopy []Lap
	if DB.Table("laps").Where("race_id = ?", race_id).First(&lapCopy).Error == nil {
		if DB.Table("laps").Where("race_id = ?", race_id).Where("lap_is_current = ?", 1).Order("lap_number desc").Order("race_total_time asc").Find(&lapsCopy).Error == nil {
			var position uint = 1
			for _, m := range lapsCopy {
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
	if DB.Table("laps").Where("race_id = ?", race_id).First(&lapCopy).Error == nil {
		if DB.Table("laps").Where("race_id = ?", race_id).Where("lap_number = ?", lap_number).Order("lap_time asc").Find(&lapsCopy).Error == nil {
			var position uint = 1
			for _, m := range lapsCopy {
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
	} else {
		fmt.Println("lapPosition: race data empty.")
		//no such race found (may be this is the first result).
		if lap_number == 0 {
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
func GetLeaderFirstLapDiscoveryUnixTime(race_id_uint uint) (leaderFirstLapDiscoveryUnixTime int64, err error) {
	var u Lap
	race_id_int := int(race_id_uint)
	err = DB.Table("laps").Where("race_id = ?", race_id_int).Where("lap_number = ?", 0).Where("lap_time = ?", 0).First(&u).Error
	if err == nil {
		//fmt.Println("leaderFirstLapDiscoveryUnixTime =", u.DiscoveryUnixTime)
		leaderFirstLapDiscoveryUnixTime = u.DiscoveryUnixTime
	}
	return
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
	result := DB.Where("race_id = ?", race_id).Where("tag_id = ?", tag_id).Order("discovery_unix_time desc").First(u)
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

//best lap time
func GetBestLapTimeFromRace(raceID uint) (bestLapTime int64, err error) {
	var lap Lap
	err = DB.Table("laps").Where("race_id = ?", raceID).Where("lap_number != ?", 0).Order("lap_time asc").First(&lap).Error
	if err == nil {
		bestLapTime = lap.LapTime
	}
	return
}

//best personal lap time
func GetBestLapTimeFromRaceByTagID(tagID string, raceID uint) (bestLapTime int64, err error) {
	var lap Lap
	err = DB.Table("laps").Where("tag_id = ?", tagID).Where("race_id = ?", raceID).Where("lap_number != ?", 0).Order("lap_time asc").First(&lap).Error
	if err == nil {
		bestLapTime = lap.LapTime
	}
	return
}

//best track lap time from all time all types
func GetBestLapTimeFromAllTime() (bestLapTime int64, err error) {
	var lapStructCopy Lap
	err = DB.Table("laps").Where("lap_number != ?", 0).Order("lap_time asc").First(&lapStructCopy).Error
	if err == nil {
		bestLapTime = lapStructCopy.LapTime
	}
	return
}

func GetPreviousLapDataFromRaceByTagID(tagID string, raceID uint) (previousLapNumber int, previousDiscoveryUnixTime, previousRaceTotalTime int64) {
	var lapStructCopy Lap
	if DB.Table("laps").Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_unix_time desc").First(&lapStructCopy).Error == nil {
		previousLapNumber = lapStructCopy.LapNumber
		previousDiscoveryUnixTime = lapStructCopy.DiscoveryUnixTime
		previousRaceTotalTime = lapStructCopy.RaceTotalTime

	} else {
		previousLapNumber = -1
		previousDiscoveryUnixTime = 0
		previousRaceTotalTime = 0
	}
	return
}

func ExpireMyPreviousLap(tagID string, raceID uint) {
	var lapStructCopy Lap
	if DB.Table("laps").Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_unix_time desc").First(&lapStructCopy).Error == nil {
		//previos lap found - update lapStructCopy.LapIsCurrent = 0
		err := DB.Model(&lapStructCopy).Update("LapIsCurrent", 0).Error
		if err != nil {
			fmt.Println("Previous lap found but update LapIsCurrent = 0 failed:", err)
		}
	}
}

func GetMyLastLapDataFromCurrentRace(u *Lap) (err error) {
	result := DB.Where("tag_id = ? AND race_id = ?", u.TagID, u.RaceID).Order("discovery_unix_time desc").First(u)
	return result.Error
}

// Get laps by tag ID
func GetAllLapsByTagId(u *[]Lap, tag_id string) (err error) {
	result := DB.Where("tag_id = ?", tag_id).Order("discovery_unix_time desc").Find(u)
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
