package Models

import ( 
	"fmt"
	"time" 
)


// Return all laps in system
func GetAllLaps(u *[]Lap) (err error) {

	result := DB.Order("discovery_time desc").Find(u)
	return result.Error
}

// current/raceid
func GetLastLap(u *Lap) (err error) {
	result := DB.Order("discovery_time desc").First(u)

	fmt.Println(time.Now().UnixNano()/int64(time.Millisecond)-300000);
	fmt.Println(u.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond))

	if (time.Now().UnixNano()/int64(time.Millisecond)-300000 > u.DiscoveryTimePrepared.UnixNano()/int64(time.Millisecond)) {
		//last lap data was created more than 300 seconds ago
		//RaceID++ (create new race)
		//LapNumber=0 (first lap)
		u.RaceID = (u.RaceID+1)
		u.LapNumber = 0

	} else {
		//last lap data was created less than 300 seconds ago
		//RaceID==RaceID (use same race)
		//LapNumber=LapNumber+1
		if DB.Where("tag_id = ? AND race_id = ?", u.TagID, u.RaceID).Order("discovery_time desc").First(u).Error != nil { 
			u.LapNumber = 0
		} else {
			u.LapNumber = u.LapNumber+1
		}	
	}

	fmt.Printf("Current RaceId = %d\n", u.RaceID)
	fmt.Printf("Current LapNumber = %d\n", u.LapNumber)
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
