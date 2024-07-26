package jwt

import (
	"e-ticketing-gin/configs"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type JWTInterface interface {
	GenerateJWT(id uint, username, email, phoneNumber, role string) map[string]any
	RefreshJWT(accessToken string, refreshToken *jwt.Token) (map[string]any, error)
	ExtractToken(g *gin.Context) (ExtractToken, error)
	ValidateRole(g *gin.Context) bool
	GetCurrentToken(g *gin.Context) *jwt.Token
}

type JWT struct {
	c *configs.ProgramConfig
}

type ExtractToken struct {
	ID          uint
	Username    string
	Email       string
	PhoneNumber string
	Role        string
}

func NewJWT(c *configs.ProgramConfig) JWTInterface {
	return &JWT{
		c: c,
	}
}

func (j *JWT) GenerateJWT(id uint, username, email, phoneNumber, role string) map[string]any {
	var result = map[string]any{}
	var accessToken = j.generateToken(id, username, email, phoneNumber, role)
	var refreshToken = j.generateRefreshToken()
	if accessToken == "" || refreshToken == "" {
		return nil
	}

	result["access_token"] = accessToken
	result["refresh_token"] = refreshToken

	return result
}

func (j *JWT) generateToken(id uint, username, email, phoneNumber, role string) string {
	var claims = jwt.MapClaims{}
	claims["id"] = id
	claims["username"] = username
	claims["email"] = email
	claims["phone_number"] = phoneNumber
	claims["role"] = role
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.c.Secret))

	if err != nil {
		return ""
	}

	return validToken
}

func (j *JWT) generateRefreshToken() string {
	var claims = jwt.MapClaims{}
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := sign.SignedString([]byte(j.c.RefSecret))

	if err != nil {
		return ""
	}
	return refreshToken
}

func (j *JWT) RefreshJWT(accessToken string, refreshToken *jwt.Token) (map[string]any, error) {
	var result = map[string]any{}
	expTime, err := refreshToken.Claims.GetExpirationTime()
	if err != nil {
		logrus.Error("Get Token Expiration Error : ", err.Error())
		return nil, errors.New("JWT : Token Expiration Error")
	}

	if refreshToken.Valid && expTime.Time.Compare(time.Now()) > 0 {
		var newClaim = jwt.MapClaims{}
		newToken, err := jwt.ParseWithClaims(accessToken, newClaim, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.c.Secret), nil
		})

		if err != nil {
			logrus.Error("Parse Token Error : ", err.Error())
			return nil, errors.New("JWT : Parse Token Error")
		}

		newClaim = newToken.Claims.(jwt.MapClaims)
		newClaim["iat"] = time.Now().Unix()
		newClaim["exp"] = time.Now().Add(time.Hour * 24).Unix()

		var newRefreshClaim = refreshToken.Claims.(jwt.MapClaims)
		newRefreshClaim["exp"] = time.Now().Add(time.Hour * 24).Unix()

		var newRefreshToken = jwt.NewWithClaims(refreshToken.Method, newRefreshClaim)
		newSignRefToken, err := newRefreshToken.SignedString(refreshToken.Signature)

		if err != nil {
			logrus.Error("Sign Refresh Token Error : ", err.Error())
			return nil, errors.New("JWT : Sign Refresh Token Error")
		}

		result["access_token"] = newToken.Raw
		result["refresh_token"] = newSignRefToken

		return result, nil
	}

	return nil, errors.New("JWT : Refresh Token Not Valid & Expired")
}

func (j *JWT) validateToken(token string) (*jwt.Token, error) {
	var authHeader = token[7:]
	parseToken, err := jwt.Parse(authHeader, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("JWT : Unexpected Signing Method : ", t.Header["alg"])
		}
		return []byte(j.c.Secret), nil
	})

	if err != nil {
		logrus.Error("Parse Token Error : ", err.Error())
		return nil, errors.New("JWT : Parse Token Error")
	}
	return parseToken, nil
}

func (j *JWT) ExtractToken(g *gin.Context) (ExtractToken, error) {
	var result = new(ExtractToken)
	authHeader := g.GetHeader("Authorization")

	token, err := j.validateToken(authHeader)

	if err != nil {
		logrus.Error("Validate Token Error : ", err.Error())
	}

	mapClaims := token.Claims.(jwt.MapClaims)
	idFloat, ok := mapClaims["id"].(float64)
	email := mapClaims["email"].(string)
	username := mapClaims["username"].(string)
	phoneNumber := mapClaims["phone_number"].(string)
	role := mapClaims["role"].(string)

	if !ok {
		return ExtractToken{}, errors.New("JWT : ID not found or not a valid number")
	}

	result.ID = uint(idFloat)
	result.Email = email
	result.Username = username
	result.PhoneNumber = phoneNumber
	result.Role = role

	return *result, nil
}

func (j *JWT) GetCurrentToken(g *gin.Context) *jwt.Token {
	authHeader := g.GetHeader("Authorization")

	token, err := j.validateToken(authHeader)

	if err != nil {
		logrus.Error("Validate Token Error : ", err.Error())
	}

	return token
}

func (j *JWT) ValidateRole(g *gin.Context) bool {
	ext, err := j.ExtractToken(g)
	if err != nil {
		logrus.Error("Validate Role Error : ", err.Error())
		return false
	}

	return ext.Role == "Admin"
}
