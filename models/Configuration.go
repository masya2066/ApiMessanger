package models

type Configuration struct {
	Param  string `json:"param"`
	Value  string `json:"value"`
	Active int    `json:"active"`
}

func InitDefaultConfiguration() {
	var config Configuration
	DB.Model(&Configuration{}).Where("param = ?", "TOKEN_LIFE_TIME").First(&config)

	if config.Value == "" || config.Param == "" {
		DB.Model(&Configuration{}).Create(&Configuration{Param: "TOKEN_LIFE_TIME", Value: "86400", Active: 1})
		DB.Model(&Configuration{}).Create(&Configuration{Param: "VERIFICATION_ATTEMPTS", Value: "3", Active: 1})
		DB.Model(&Configuration{}).Create(&Configuration{Param: "LANGUAGE", Value: "en", Active: 1})
	}
}
