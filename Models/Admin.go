package Models


// Return all admins in system order by date
func GetAllAdmins(u *[]Admin) (err error) {

	result := DB.Find(u)
	return result.Error
}

func AddNewAdmin(u *Admin) (err error) {
	if err = DB.Create(u).Error; err != nil {
		return err
	}

	return nil
}

func GetOneAdmin(u *Admin, admin_id string) (err error) {
	if err := DB.Where("id = ?", admin_id).First(u).Error; err != nil {
		return err
	}

	return nil
}

func GetOneAdminByLogin(u *Admin, login string) (err error) {
    if err := DB.Where("login = ?", login).First(u).Error; err != nil {
		return err
	}

	return nil
}

func PutOneAdmin(u *Admin) (err error) {
	DB.Save(u)
	return nil
}

func DeleteOneAdmin(u *Admin, admin_id string) (err error) {
	DB.Where("id = ?", admin_id).Delete(u)
	return nil
}
