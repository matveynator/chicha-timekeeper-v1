package Models


// Return all checkin in system
func GetAllCheckins(u *[]Checkin) (err error) {

	result := DB.Find(u)
	return result.Error
}

func AddNewCheckin(u *Checkin) (err error) {
	if err = DB.Create(u).Error; err != nil {
		return err
	}

    return nil
}

func GetOneCheckin(u *Checkin, checkin_id string) (err error) {
	if err := DB.Where("id = ?", checkin_id).First(u).Error; err != nil {
		return err
	}

	return nil
}

func PutOneCheckin(u *Checkin) (err error) {
	DB.Save(u)
	return nil
}

func DeleteOneCheckin(u *Checkin, checkin_id string) (err error) {
	DB.Where("id = ?", checkin_id).Delete(u)
	return nil
}
