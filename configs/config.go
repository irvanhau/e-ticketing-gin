package configs

import (
	"errors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type ProgramConfig struct {
	Server       int
	DBPort       int
	DBHost       string
	DBUser       string
	DBPass       string
	DBName       string
	Email        string
	Password     string
	Secret       string
	RefSecret    string
	CloudURL     string
	MidServerKey string
	MidEnv       string
}

func InitConfig() *ProgramConfig {
	var res = new(ProgramConfig)
	res, errorRes := loadConfig()

	logrus.Error(errorRes)
	if errorRes != nil {
		logrus.Error("Config : Cannot start program, Failed to load configuration")
		return nil
	}
	return res
}

func loadConfig() (*ProgramConfig, error) {
	var errorLoad error
	var res = new(ProgramConfig)
	var permit = true

	if val, found := os.LookupEnv("SERVER"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : Invalid Port Value, ", err.Error())
			permit = false
		}
		res.Server = port
	} else {
		permit = false
		errorLoad = errors.New("SERVER PORT UNDEFINED")
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		res.DBHost = val
	} else {
		permit = false
		errorLoad = errors.New("DBHOST UNDEFINED")
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : Invalid DB Port Value, ", err.Error())
			permit = false
		}
		res.DBPort = port
	} else {
		permit = false
		errorLoad = errors.New("DBPORT UNDEFINED")
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		res.DBName = val
	} else {
		permit = false
		errorLoad = errors.New("DBNAME UNDEFINED")
	}

	if val, found := os.LookupEnv("DBUSER"); found {
		res.DBUser = val
	} else {
		permit = false
		errorLoad = errors.New("DBUSER UNDEFINED")
	}

	if val, found := os.LookupEnv("DBPASSWORD"); found {
		res.DBPass = val
	} else {
		permit = false
		errorLoad = errors.New("DBPASSWORD UNDEFINED")
	}

	if val, found := os.LookupEnv("EMAIL"); found {
		res.Email = val
	} else {
		permit = false
		errorLoad = errors.New("EMAIL UNDEFINED")
	}

	if val, found := os.LookupEnv("PASSWORD"); found {
		res.Password = val
	} else {
		permit = false
		errorLoad = errors.New("PASSWORD UNDEFINED")
	}

	if val, found := os.LookupEnv("SECRET"); found {
		res.Secret = val
	} else {
		permit = false
		errorLoad = errors.New("SECRET UNDEFINED")
	}

	if val, found := os.LookupEnv("REFSECRET"); found {
		res.RefSecret = val
	} else {
		permit = false
		errorLoad = errors.New("REFSECRET UNDEFINED")
	}

	if val, found := os.LookupEnv("CLOUDURL"); found {
		res.CloudURL = val
	} else {
		permit = false
		errorLoad = errors.New("CLOUDURL UNDEFINED")
	}

	if val, found := os.LookupEnv("MT_SERVER_KEY"); found {
		res.MidServerKey = val
	} else {
		permit = false
		errorLoad = errors.New("MT_SERVER_KEY UNDEFINED")
	}

	if val, found := os.LookupEnv("MT_ENV"); found {
		res.MidEnv = val
	} else {
		permit = false
		errorLoad = errors.New("MT_ENV UNDEFINED")
	}

	if !permit {
		return nil, errorLoad
	}

	return res, nil
}
