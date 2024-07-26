package data

import (
	"e-ticketing-gin/features/users"
	"e-ticketing-gin/helper/enkrip"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserData struct {
	db     *gorm.DB
	enkrip enkrip.HashInterface
}

func New(db *gorm.DB, e enkrip.HashInterface) *UserData {
	return &UserData{
		db:     db,
		enkrip: e,
	}
}

func (ud *UserData) Register(newData users.User) (*users.User, error) {
	var dbData = new(User)
	dbData.Username = newData.Username
	dbData.Email = newData.Email
	dbData.PhoneNumber = newData.PhoneNumber
	dbData.Password = newData.Password
	dbData.IsAdmin = newData.IsAdmin
	dbData.Status = newData.Status

	if err := ud.db.Create(dbData).Error; err != nil {
		logrus.Error("DATA : Register Error : ", err.Error())
		return nil, err
	}

	return &newData, nil
}

func (ud *UserData) Login(username, password string) (*users.User, error) {
	var dbdata = new(User)
	var dataCount int64
	dbdata.Username = username

	var qry = ud.db.Where("username = ? AND status = ?", username, true).First(dbdata)
	qry.Count(&dataCount)

	if dataCount == 0 {
		logrus.Error("DATA : Login Error : Data Not Found")
		return nil, errors.New("ERROR Data Not Found")
	}

	if err := qry.Error; err != nil {
		logrus.Error("DATA : Get Username Error : ", err.Error())
		return nil, err
	}

	if err := ud.enkrip.Compare(dbdata.Password, password); err != nil {
		logrus.Error("DATA : Incorrect Password")
		return nil, errors.New("ERROR Incorrect Password")
	}

	var result = new(users.User)
	result.ID = dbdata.ID
	result.Username = dbdata.Username
	result.Email = dbdata.Email
	result.PhoneNumber = dbdata.PhoneNumber
	result.IsAdmin = dbdata.IsAdmin
	result.Status = dbdata.Status

	return result, nil
}

func (ud *UserData) GetByID(id int) (users.User, error) {
	var listUser users.User
	var qry = ud.db.Where("id = ? ", id).Where("status = ?", true).First(&listUser)

	if err := qry.Error; err != nil {
		logrus.Error("DATA : Error Get By ID : ", err.Error())
		return listUser, err
	}

	return listUser, nil
}

func (ud *UserData) GetByUsername(username string) (*users.User, error) {
	var dbData = new(User)
	dbData.Username = username
	var qry = ud.db.Where("username = ? ", dbData.Username).Where("status = ?", true).First(dbData)

	if err := qry.Error; err != nil {
		logrus.Error("DATA : Error Get By ID : ", err.Error())
		return nil, err
	}

	var result = new(users.User)
	result.ID = dbData.ID
	result.Username = dbData.Username
	result.Email = dbData.Email
	result.PhoneNumber = dbData.PhoneNumber
	result.IsAdmin = dbData.IsAdmin
	result.Status = dbData.Status

	return result, nil
}

func (ud *UserData) CheckUsername(username string) bool {
	var count int64
	var qry = ud.db.Table("users").Where("username = ? ", username).Count(&count)

	if err := qry.Error; err != nil {
		logrus.Error("DATA : Check Username Error : ", err.Error())
		return false
	}

	if count == 0 {
		return true
	}

	return false
}

func (ud *UserData) InsertCodeReset(username, code string) error {
	var newData = new(UserResetPass)
	newData.Username = username
	newData.Code = code
	newData.ExpiredAt = time.Now().Add(time.Minute * 10)

	_, err := ud.GetByCodeReset(code)
	if err != nil {
		errDelete := ud.DeleteCodeReset(code)
		if errDelete != nil {
			return errors.New("ERROR Delete CodeReset Error : " + errDelete.Error())
		}
	}

	if err := ud.db.Table("user_reset_passes").Create(newData).Error; err != nil {
		logrus.Error("DATA : Insert Code Reset Pass Error : ", err.Error())
		return err
	}

	return nil
}

func (ud *UserData) DeleteCodeReset(code string) error {
	var deleteData = new(UserResetPass)

	if err := ud.db.Table("user_reset_passes").Where("code = ?", code).Delete(deleteData).Error; err != nil {
		logrus.Error("DATA : Delete Code Reset Pass Error : ", err.Error())
		return err
	}

	return nil
}

func (ud *UserData) GetByCodeReset(code string) (*users.UserResetPass, error) {
	var dbData = new(UserResetPass)
	dbData.Code = code

	if err := ud.db.Table("user_reset_passes").Where("code = ?", dbData.Code).First(&dbData).Error; err != nil {
		logrus.Error("DATA : Get By Code Reset Pass Error : ", err.Error())
		return nil, err
	}

	var result = new(users.UserResetPass)
	result.Code = dbData.Code
	result.ExpiredAt = dbData.ExpiredAt
	result.Username = dbData.Username

	return result, nil
}
func (ud *UserData) ResetPassword(code, username, password string) error {
	if err := ud.db.Where("username = ?", username).Update("password", password).Error; err != nil {
		logrus.Error("DATA : Reset Password Error : ", err.Error())
		return err
	}

	checkData, err := ud.GetByCodeReset(code)
	if err != nil {
		logrus.Error("DATA : Get By Code Reset Error : ", err.Error())
		return err
	}

	if checkData.Code != "" {
		err := ud.DeleteCodeReset(code)
		if err != nil {
			logrus.Error("DATA : Delete Code Reset Error : ", err.Error())
			return err
		}
	}

	return nil
}

func (ud *UserData) UpdateProfile(id int, newData users.UpdateProfile) (bool, error) {
	var qry = ud.db.Where("id = ? ", id).Updates(User{
		Username:    newData.Username,
		Email:       newData.Email,
		PhoneNumber: newData.PhoneNumber,
	})

	if err := qry.Error; err != nil {
		logrus.Error("DATA : Error Update Profile : ", err.Error())
		return false, err
	}

	if datacount := qry.RowsAffected; datacount < 1 {
		logrus.Error("DATA : Update Profile Error : No Row Affected")
		return false, errors.New("ERROR Update Profile Error : No Row Affected")
	}

	return true, nil
}

func (ud *UserData) GetAll() ([]users.User, error) {
	var listUser []users.User

	if err := ud.db.Find(&listUser).Error; err != nil {
		return nil, err
	}

	return listUser, nil
}

func (ud *UserData) Activate(id int) (bool, error) {
	var qry = ud.db.Where("id = ?", id).Updates(User{Status: true})

	if err := qry.Error; err != nil {
		return false, err
	}

	return true, nil
}

func (ud *UserData) Deactivate(id int) (bool, error) {
	var qry = ud.db.Model(User{}).Where("id = ?", id).Updates(map[string]interface{}{"status": false})

	if err := qry.Error; err != nil {
		return false, err
	}

	return true, nil
}

func (ud *UserData) UserDashboard() (users.UserDashboard, error) {
	var dashboardUser users.UserDashboard

	tUser, tNewUser, tUserActive, tUserInactive := ud.getTotalUser()

	dashboardUser.TotalUser = tUser
	dashboardUser.TotalNewUser = tNewUser
	dashboardUser.TotalUserActive = tUserActive
	dashboardUser.TotalUserInactive = tUserInactive

	return dashboardUser, nil
}

func (ud *UserData) getTotalUser() (int, int, int, int) {
	var totalUser int64
	var totalNewUser int64
	var totalUserActive int64
	var totalUserInactive int64

	var now = time.Now()
	var before = now.AddDate(0, 0, -7)

	var _ = ud.db.Count(&totalUser)
	var _ = ud.db.Where("created_at BETWEEN ? AND ?", before, now).Count(&totalNewUser)
	var _ = ud.db.Where("status = ?", true).Count(&totalUserActive)
	var _ = ud.db.Where("status = ?", false).Count(&totalUserInactive)

	totalUserInt := int(totalUser)
	totalNewUserInt := int(totalNewUser)
	totalUserActiveInt := int(totalUserActive)
	totalUserInactiveInt := int(totalUserInactive)

	return totalUserInt, totalNewUserInt, totalUserActiveInt, totalUserInactiveInt
}

func (ud *UserData) InsertCodeVerification(username, code string) error {
	var newData = new(UserVerification)
	newData.Username = username
	newData.Code = code
	newData.ExpiredAt = time.Now().Add(time.Minute * 10)

	_, err := ud.GetByCodeVerification(code)
	if err != nil {
		errDel := ud.DeleteCodeVerification(code)
		if errDel != nil {
			logrus.Error("DATA : Delete Code Verification Error : ", err.Error())
			return err
		}
	}

	if err := ud.db.Table("user_verifications").Create(newData).Error; err != nil {
		logrus.Error("DATA : Insert Code Verification Error : ", err.Error())
		return err
	}

	return nil
}
func (ud *UserData) DeleteCodeVerification(code string) error {
	var deleteData = new(UserVerification)

	if err := ud.db.Table("user_verifications").Where("code = ?", code).Delete(deleteData).Error; err != nil {
		logrus.Error("DATA : Delete Code Verification Error : ", err.Error())
		return err
	}

	return nil
}
func (ud *UserData) GetByCodeVerification(code string) (*users.UserVerification, error) {
	var dbData = new(UserVerification)
	dbData.Code = code

	if err := ud.db.Table("user_verifications").Where("code = ?", code).First(&dbData).Error; err != nil {
		logrus.Error("DATA : Get By Code Verification Error : ", err.Error())
		return nil, err
	}

	var result = new(users.UserVerification)
	result.Code = dbData.Code
	result.ExpiredAt = dbData.ExpiredAt
	result.Username = dbData.Username

	return result, nil
}
func (ud *UserData) UserVerification(code, username string) error {
	if err := ud.db.Where("username = ?", username).Update("status", true).Error; err != nil {
		logrus.Error("DATA : Update User Verification Error : ", err.Error())
		return err
	}

	checkData, _ := ud.GetByCodeVerification(code)
	if checkData.Code != "" {
		errDel := ud.DeleteCodeVerification(code)
		if errDel != nil {
			logrus.Error("DATA : Delete Code Verification Error : ", errDel.Error())
		}
	}

	return nil
}
