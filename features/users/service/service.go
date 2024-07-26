package service

import (
	"e-ticketing-gin/features/users"
	"e-ticketing-gin/helper/email"
	"e-ticketing-gin/helper/enkrip"
	"e-ticketing-gin/helper/jwt"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type UserService struct {
	data  users.UserDataInterface
	hash  enkrip.HashInterface
	jwt   jwt.JWTInterface
	email email.EmailInterface
}

func New(d users.UserDataInterface, e enkrip.HashInterface, j jwt.JWTInterface, em email.EmailInterface) *UserService {
	return &UserService{
		data:  d,
		hash:  e,
		jwt:   j,
		email: em,
	}
}

func (u *UserService) Register(newData users.User) (*users.User, error) {
	isAlready := u.data.CheckUsername(newData.Username)

	if !isAlready {
		logrus.Error("Service : Username already registered")
		return nil, errors.New("ERROR Username already registered")
	}

	hashPassword, err := u.hash.HashPassword(newData.Password)
	if err != nil {
		logrus.Error("Service : Error Hash Password : ", err.Error())
		return nil, errors.New("ERROR Error Hashing Password")
	}

	newData.Password = hashPassword
	newData.IsAdmin = false
	newData.Status = false

	result, err := u.data.Register(newData)
	if err != nil {
		logrus.Error("Service : Error Register : ", err.Error())
		return nil, errors.New("ERROR Error Register")
	}

	return result, nil
}
func (u *UserService) Login(username string, password string) (*users.UserCredential, error) {
	result, err := u.data.Login(username, password)

	if err != nil {
		if strings.Contains(err.Error(), "Incorrect Password") {
			return nil, errors.New("ERROR Incorrect Password")
		}
		if strings.Contains(err.Error(), "Not Found") {
			return nil, errors.New("ERROR Not Found")
		}
		return nil, errors.New("ERROR Process Failed")
	}

	role := "user"

	if result.IsAdmin {
		role = "admin"
	}

	tokenData := u.jwt.GenerateJWT(result.ID, result.Username, result.Email, result.PhoneNumber, role)

	if tokenData == nil {
		logrus.Error("Service : Error Generate JWT")
		return nil, errors.New("ERROR Generate JWT")
	}

	response := new(users.UserCredential)
	response.Access = tokenData
	response.Username = result.Username

	return response, nil
}

func (u *UserService) ForgetPasswordWeb(username string) error {
	user, err := u.data.GetByUsername(username)
	if err != nil {
		logrus.Error("Service : Error Get By Username : ", err.Error())
		return errors.New("ERROR Error Get By Username")
	}

	email := user.Email

	username = user.Username
	header, htmlBody, code := u.email.HTMLBodyReset(username)

	if err := u.data.InsertCodeReset(username, code); err != nil {
		logrus.Error("Service : Error Insert Code Reset User : ", err.Error())
		return errors.New("ERROR Error Insert Code Reset User")
	}

	errSend := u.email.SendEmail(email, header, htmlBody)
	if errSend != nil {
		logrus.Error("Service : Error Sending Email : ", err.Error())
		return errors.New("ERROR Sending Email")
	}

	return nil
}

func (u *UserService) TokenResetVerify(code string) (*users.UserResetPass, error) {
	result, err := u.data.GetByCodeReset(code)

	if err != nil {
		logrus.Error("Service : Error Get By Code : ", err.Error())
		return nil, errors.New("ERROR Error Get By Code")
	}

	if result.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("ERROR Token Expired")
	}

	return result, nil
}

func (u *UserService) ResetPassword(code, username, password string) error {
	hashPassword, err := u.hash.HashPassword(password)
	if err != nil {
		logrus.Error("Service : Error Hash Password : ", err.Error())
		return errors.New("ERROR Error Hashing Password")
	}

	password = hashPassword

	if err := u.data.ResetPassword(code, username, password); err != nil {
		logrus.Error("Service : Error Reset Password : ", err.Error())
		return errors.New("ERROR Error Reset Password")
	}

	return nil
}
func (u *UserService) UpdateProfile(id int, newData users.UpdateProfile) (bool, error) {
	res, err := u.data.UpdateProfile(id, newData)

	if err != nil {
		logrus.Error("Service : Error Update Profile : ", err.Error())
		return false, errors.New("ERROR Error Update Profile")
	}

	return res, nil
}
func (u *UserService) Profile(id int) (*users.User, error) {
	res, err := u.data.GetByID(id)
	if err != nil {
		logrus.Error("Service : Error Get ByID : ", err.Error())
		return nil, errors.New("ERROR Error Get ByID")
	}
	return &res, nil
}

func (u *UserService) GetAll() ([]users.User, error) {
	res, err := u.data.GetAll()
	if err != nil {
		logrus.Error("Service : Error GetAll : ", err.Error())
		return nil, errors.New("ERROR Error GetAll")
	}

	return res, nil
}
func (u *UserService) Activate(id int) (bool, error) {
	res, err := u.data.Activate(id)
	if err != nil {
		logrus.Error("Service : Error Activate : ", err.Error())
		return false, errors.New("ERROR Error Activate")
	}
	return res, nil
}
func (u *UserService) Deactivate(id int) (bool, error) {
	res, err := u.data.Deactivate(id)
	if err != nil {
		logrus.Error("Service : Error Deactivate : ", err.Error())
		return false, errors.New("ERROR Error Deactivate")
	}
	return res, nil
}

func (u *UserService) UserDashboard() (users.UserDashboard, error) {
	res, err := u.data.UserDashboard()
	if err != nil {
		logrus.Error("Service : Error User Dashboard : ", err.Error())
		return res, errors.New("ERROR Error User Dashboard")
	}

	return res, nil
}
func (u *UserService) UserVerificationCode(username, email string) error {
	header, htmlBody, code := u.email.HTMLBodyVerification(username)

	if err := u.data.InsertCodeVerification(username, code); err != nil {
		logrus.Error("Service : Error Insert Code Verification : ", err.Error())
		return errors.New("ERROR Error Insert Code Verification")
	}

	errSend := u.email.SendEmail(email, header, htmlBody)

	if errSend != nil {
		logrus.Error("Service : Error Send Email Verification : ", errSend.Error())
		return errors.New("ERROR Send Email Verification")
	}

	return nil
}
func (u *UserService) UserVerification(code, username string) error {
	if err := u.data.UserVerification(username, code); err != nil {
		logrus.Error("Service : Error User Verification : ", err.Error())
		return errors.New("ERROR Error User Verification")
	}

	return nil
}
func (u *UserService) TokenVerificationResetVerify(code string) (*users.UserVerification, error) {
	res, err := u.data.GetByCodeVerification(code)
	if err != nil {
		logrus.Error("Service : Error Get By Code : ", err.Error())
		return nil, errors.New("ERROR Error Get By Code")
	}

	if res.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("ERROR Token Expired")
	}

	return res, nil
}
