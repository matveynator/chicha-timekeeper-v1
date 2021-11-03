package Models


// Return all user in system
func GetAllUsers(u *[]User) (err error) {

	result := DB.Find(u)
	return result.Error
}

func AddNewUser(u *User) (err error) {
	if err = DB.Create(u).Error; err != nil {
		return err
	}

    return nil
}

func GetOneUser(u *User, user_id string) (err error) {
	if err := DB.Where("id = ?", user_id).First(u).Error; err != nil {
		return err
	}

	return nil
}

func PutOneUser(u *User) (err error) {
	DB.Save(u)
	return nil
}

func DeleteOneUser(u *User, user_id string) (err error) {
	DB.Where("id = ?", user_id).Delete(u)
	return nil
}
