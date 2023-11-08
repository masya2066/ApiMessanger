package models

type Message struct {
	UserId    int    `json:"user_id"`
	ChatId    string `json:"chat_id"`
	MessageId string `json:"message_id"`
	Message   string `json:"message"`
	Owner     int    `json:"owner"`
	Changed   bool   `json:"changed"`
	Deleted   bool   `json:"deleted"`
	Created   string `json:"created"` //os.Getenv("DATE_FORMAT")
	Update    string `json:"update"`  //os.Getenv("DATE_FORMAT")
}
