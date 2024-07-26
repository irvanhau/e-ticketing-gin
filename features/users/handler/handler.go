package handler

import (
	"e-ticketing-gin/features/users"
	"e-ticketing-gin/helper"
	"e-ticketing-gin/helper/jwt"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	service users.UserServiceInterface
	jwt     jwt.JWTInterface
}

func NewHandler(jwt jwt.JWTInterface, service users.UserServiceInterface) *UserHandler {
	return &UserHandler{
		jwt:     jwt,
		service: service,
	}
}

func (u *UserHandler) Register(c *gin.Context) {
	var req RegisterInput

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error("Handler : Bind Input Error : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid User Input", nil))
		return
	}

	isValid, errors := helper.ValidateJSON(req)
	if !isValid {
		c.JSON(http.StatusBadRequest, helper.FormatResponseValidation("Invalid Format Request", errors))
		return
	}

	if !helper.ValidatePassword(req.Password) {
		errPass := []string{"Password must contain a combination letters, symbols, and numbers"}
		c.JSON(http.StatusBadRequest, helper.FormatResponseValidation("Invalid Password", errPass))
		return
	}

	var serviceInput = new(users.User)
	serviceInput.Email = req.Email
	serviceInput.Username = req.Username
	serviceInput.PhoneNumber = req.PhoneNumber
	serviceInput.Password = req.Password

	res, errData := u.service.Register(*serviceInput)
	if errData != nil {
		if strings.Contains(errData.Error(), "Username already registered") {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("Username Already Registered", nil))
			return
		}
		logrus.Error("Handler : Register Error : ", errData.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Register Process Failed", nil))
		return
	}

	verificationCode := u.service.UserVerificationCode(req.Username, req.Email)
	if verificationCode != nil {
		logrus.Error("Handler : Send Email Error : ", verificationCode.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Send Email Failed", nil))
		return
	}

	var response = new(RegisterResponse)
	response.Email = res.Email
	response.Username = res.Username
	response.PhoneNumber = res.PhoneNumber

	c.JSON(http.StatusCreated, helper.FormatResponse("Register Success, Please check your email for verification code", response))
}
func (u *UserHandler) Login(c *gin.Context) {
	var input = new(LoginInput)

	if err := c.ShouldBindJSON(input); err != nil {
		logrus.Error("Handler : Bind Input Error : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Input", nil))
		return
	}

	isValid, errors := helper.ValidateJSON(input)
	if !isValid {
		c.JSON(http.StatusBadRequest, helper.FormatResponseValidation("Invalid Format Request", errors))
		return
	}

	res, err := u.service.Login(input.Username, input.Password)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			c.JSON(http.StatusNotFound, helper.FormatResponse("User Not Found", nil))
			return
		}
		if strings.Contains(err.Error(), "Incorrect Password") {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("Incorrect Password", nil))
			return
		}
		logrus.Error("Handler : Login Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Login Process Failed", nil))
		return
	}

	var response = new(LoginResponse)
	response.Username = res.Username
	response.Token = res.Access

	c.JSON(http.StatusOK, helper.FormatResponse("Success Login", response))
	return
}
func (u *UserHandler) ForgetPasswordWeb(c *gin.Context) {
	var input = new(ForgetPasswordInput)

	if err := c.ShouldBindJSON(input); err != nil {
		logrus.Error("Handler : Bind Input Error : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Input", nil))
		return
	}

	isValid, errors := helper.ValidateJSON(input)
	if !isValid {
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Format Request", errors))
		return
	}

	res := u.service.ForgetPasswordWeb(input.Username)

	if res != nil {
		logrus.Error("Handler : Send Email Error")
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Forget Password Error", res))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Send Reset Code to Email", nil))
	return
}
func (u *UserHandler) ResetPassword(c *gin.Context) {
	var token = c.Query("token_reset_password")
	if token == "" {
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Token Not Found", nil))
		return
	}

	dataToken, err := u.service.TokenResetVerify(token)
	if err != nil {
		if strings.Contains(err.Error(), "Token Expired") {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("Token Expired", nil))
			return
		}
		logrus.Error("Handler : Token Reset Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Token Reset Failed", nil))
		return
	}

	var input = new(ResetPasswordInput)
	if err := c.ShouldBindJSON(input); err != nil {
		logrus.Error("Handler : Bind Input Error : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Input", nil))
		return
	}

	isValid, errors := helper.ValidateJSON(input)
	if !isValid {
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Format Request", errors))
		return
	}

	if input.Password != input.PasswordConfirm {
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Password Not Match", nil))
		return
	}

	if !helper.ValidatePassword(input.Password) {
		errPass := []string{"Password must contain a combination letters, symbols, and numbers"}
		c.JSON(http.StatusBadRequest, helper.FormatResponseValidation("Invalid Format Request", errPass))
		return
	}

	result := u.service.ResetPassword(dataToken.Code, dataToken.Username, input.Password)

	if result != nil {
		logrus.Error("Handler : Reset Password Error")
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Reset Password Error", result))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Reset Password", nil))
	return
}

func (u *UserHandler) UpdateProfile(c *gin.Context) {
	ext, err := u.jwt.ExtractToken(c)

	if err != nil {
		logrus.Error("Handler : Extract Token Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Extract Token Error", nil))
		return
	}

	id := ext.ID
	var input = new(UpdateProfile)
	if err := c.ShouldBindJSON(input); err != nil {
		logrus.Error("Handler : Bind Input Error : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Input", nil))
		return
	}

	var serviceUpdate = new(users.UpdateProfile)
	serviceUpdate.Email = input.Email
	serviceUpdate.Username = input.Username
	serviceUpdate.PhoneNumber = input.PhoneNumber

	res, err := u.service.UpdateProfile(int(id), *serviceUpdate)
	if err != nil {
		logrus.Error("Handler : Update Profile Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Update Profile Error", res))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Update Profile", res))
	return
}

func (u *UserHandler) RefreshToken(c *gin.Context) {
	var input = new(RefreshTokenInput)
	if err := c.ShouldBindJSON(input); err != nil {
		logrus.Error("Handler : Bind Input Error : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid Input", nil))
		return
	}

	var currentToken = u.jwt.GetCurrentToken(c)

	res, err := u.jwt.RefreshJWT(input.Token, currentToken)

	if err != nil {
		logrus.Error("Handler : Refresh Token Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Refresh Token Error", nil))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Refresh Token", res))
}
func (u *UserHandler) Profile(c *gin.Context) {
	ext, err := u.jwt.ExtractToken(c)
	if err != nil {
		logrus.Error("Handler : Extract Token Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Extract Token Error", nil))
		return
	}

	id := int(ext.ID)

	res, err := u.service.Profile(id)
	if err != nil {
		logrus.Error("Handler : Profile Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Profile Error", res))
		return
	}

	var response = new(UserInfo)
	response.Username = res.Username
	response.PhoneNumber = res.PhoneNumber
	response.Email = res.Email
	response.Role = ext.Role

	c.JSON(http.StatusOK, helper.FormatResponse("Success Get Profile", response))
}

func (u *UserHandler) GetUsers(c *gin.Context) {
	if err := u.jwt.ValidateRole(c); !err {
		logrus.Error("Handler : Unauthorized Access : ", errors.New("you have no permission to access this feature"))
		c.JSON(http.StatusUnauthorized, helper.FormatResponse("Restricted Access", nil))
		return
	}

	res, err := u.service.GetAll()

	if err != nil {
		logrus.Error("Handler : Get Users Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Get Users Error", res))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Get Users", res))
}
func (u *UserHandler) ActivateUser(c *gin.Context) {
	if err := u.jwt.ValidateRole(c); !err {
		logrus.Error("Handler : Unauthorized Access : ", errors.New("you have no permission to activate this feature"))
		c.JSON(http.StatusUnauthorized, helper.FormatResponse("Restricted Access", nil))
		return
	}

	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error("Handler : Invalid ID : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid User ID", nil))
		return
	}

	res, err := u.service.Activate(userId)
	if err != nil {
		logrus.Error("Handler : Activate User Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Activate User Error", res))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Activate User", res))
	return
}
func (u *UserHandler) DeactivateUser(c *gin.Context) {
	if err := u.jwt.ValidateRole(c); !err {
		logrus.Error("Handler : Unauthorized Access : ", errors.New("you have no permission to deactivate this feature"))
		c.JSON(http.StatusUnauthorized, helper.FormatResponse("Restricted Access", nil))
		return
	}

	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error("Handler : Invalid ID : ", err.Error())
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid User ID", nil))
		return
	}

	res, err := u.service.Deactivate(userId)
	if err != nil {
		logrus.Error("Handler : Deactivate User Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Deactivate User Error", res))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success Deactivate User", res))
	return
}

func (u *UserHandler) UserDashboard(c *gin.Context) {
	if err := u.jwt.ValidateRole(c); !err {
		logrus.Error("Handler : Unauthorized Access : ", errors.New("you have no permission to access this feature"))
		c.JSON(http.StatusUnauthorized, helper.FormatResponse("Restricted Access", nil))
		return
	}

	res, err := u.service.UserDashboard()
	if err != nil {
		logrus.Error("Handler : User Dashboard : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("User Dashboard", res))
	}

	var response = new(DashboardResponse)
	response.TotalUserBaru = res.TotalNewUser
	response.TotalUser = res.TotalUser
	response.TotalUserActive = res.TotalUserActive
	response.TotalUserInactive = res.TotalUserInactive

	c.JSON(http.StatusOK, helper.FormatResponse("Success Get User Dashboard", response))
	return
}
func (u *UserHandler) UserVerification(c *gin.Context) {
	var token = c.Query("token_verification")
	if token == "" {
		c.JSON(http.StatusBadRequest, helper.FormatResponse("Token Not Found", nil))
	}

	dataToken, err := u.service.TokenVerificationResetVerify(token)
	if err != nil {
		logrus.Error("Handler : Token Reset Error : ", err.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("Token Reset Error", nil))
		return
	}

	res := u.service.UserVerification(dataToken.Code, dataToken.Username)
	if res != nil {
		logrus.Error("Handler : User Verification Error : ", res.Error())
		c.JSON(http.StatusInternalServerError, helper.FormatResponse("User Verification Error", nil))
		return
	}

	c.JSON(http.StatusOK, helper.FormatResponse("Success to verification, enable to login", res))
	return
}
