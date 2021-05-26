package Models


// Return all laps in system
func GetAllLaps(u *[]Lap) (err error) {

	result := DB.Find(u)
	return result.Error
}

// Get laps by tag ID
func GetAllLapsByTagId(u *[]Lap, tag_id string) (err error) {
    result := DB.Where("tag_id = ?", tag_id).Find(u)
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
