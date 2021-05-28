package Models


// Return all race in system
func GetAllRaces(u *[]Race) (err error) {

	result := DB.Find(u)
	return result.Error
}

func GetAllUnfinishedRace(u *[]Race) (err error) {

	result := DB.Where("is_active = '1'").Find(u)
	return result.Error
}

func AddNewRace(u *Race) (err error) {
	if err = DB.Create(u).Error; err != nil {
		return err
	}

    return nil
}

func GetOneRace(u *Race, race_id string) (err error) {
	if err := DB.Where("id = ?", race_id).First(u).Error; err != nil {
		return err
	}

	return nil
}

func GetOneUnfinishedRace(u *Race) (err error) {
	if err := DB.Where("is_active = '1'").First(u).Error; err != nil {
		return err
	}

	return nil
}

func PutOneRace(u *Race) (err error) {
	DB.Save(u)
	return nil
}

func DeleteOneRace(u *Race, race_id string) (err error) {
	DB.Where("id = ?", race_id).Delete(u)
	return nil
}
