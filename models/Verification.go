package models

import (
	"os"
	"time"
)

type SmsCode struct {
	Number   string `json:"number"`
	Code     int    `json:"code"`
	Sent     bool   `json:"sent"`
	Attempts int    `json:"attempts"`
	SentTime string `json:"sent_time"`
	Created  string `json:"created"`
}

type EmailCode struct {
	Code      string `json:"code"`
	SendCount int    `json:"send_count"`
	Attempts  int    `json:"attempts"`
	Created   string `json:"created"`
}

type SmsBody struct {
	Number string `json:"number"`
	Code   int    `json:"code"`
}

func AttemptSubmitSms(number string, userCode int) (int, bool) {
	var code SmsCode
	DB.Model(&code).Where("number = ?", number).First(&code)

	if code.Attempts >= 3 {
		now := time.Now().UTC().Add(time.Second * -180).Format(os.Getenv("DATE_FORMAT"))
		if now <= code.Created {
			return 3, false
		}
		DB.Model(&code).Where("number = ?", number).Delete(&code)
		DB.Create(SmsCode{Code: userCode, Sent: true, Attempts: 0, Created: time.Now().UTC().Format(os.Getenv("DATE_FORMAT"))})
	}

	if code.Code != userCode {
		DB.Model(code).Where("number = ?", number).Update("attempts", code.Attempts+1)
		return code.Attempts + 1, false
	}
	return 0, true
}

func CheckAccessToSendSms(number string) bool {
	var userSms SmsCode

	DB.Model(&userSms).Where("number = ?", number).First(&userSms)

	if userSms.Sent == true {
		if userSms.SentTime == "" {
			panic("sent_time in empty")
		}
		now := time.Now().UTC().Add(time.Second * -600).Format(os.Getenv("DATE_FORMAT"))
		if userSms.SentTime >= now {
			return false
		}
		DB.Model(userSms).Where("number = ?", number).Delete(userSms)
		return true
	}

	return true
}
