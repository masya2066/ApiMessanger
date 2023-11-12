package configs

import "ApiMessenger/models"

func System(param string) string {
	var config models.Configuration

	models.DB.Model(&config).Where("param = ?", param).First(&config)

	return config.Value
}
