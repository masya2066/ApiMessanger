package controllers

import (
	"ApiMessenger/consumers"
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type ChatBody struct {
	Name    string `json:"name"`
	Members []int  `json:"members"`
}

func NewChat(c *gin.Context) {
	var chat models.Chat
	var user models.User

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if cookie == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	var body ChatBody

	_ = c.ShouldBindJSON(&body)

	if utils.IsArray(body.Members) == false {
		c.JSON(400, ErrorMsg(25, language.Language("members_is_not_array")))
		return
	}

	if utils.HasDuplicatesInArray(body.Members) == true {
		c.JSON(400, ErrorMsg(26, language.Language("members_has_duplicated")))
		return
	}

	if body.Name == "" {
		c.JSON(403, ErrorMsg(20, language.Language("invalid_chat_name")))
		return
	}

	models.DB.Model(&models.User{}).Where("email = ?", parse.Subject).First(&user)

	code := utils.GenerateId()
	chat.Name = body.Name
	chat.Owner = user.ID
	chat.ChatId = code

	var chatMembers models.ChatMembers
	var checkUser models.User
	var users []int

	for i := 0; i < len(body.Members); i++ {
		models.DB.Model(&models.User{}).Where("ID = ?", body.Members[i]).First(&checkUser)

		if checkUser.ID == 0 {
			users = append(users, body.Members[i])
		}
		checkUser.ID = 0
	}

	if len(users) != 0 {
		c.JSON(400, ErrorMsg(26, language.Language("error_invite_member_to_chat")+utils.IntSliceToString(users)))
		return
	}

	chat.Created = time.Now()
	chat.Updated = time.Now()
	chat.Phrase = utils.GenerateRandomSecretPhrase()

	models.DB.Create(&chat)

	models.DB.Create(&models.ChatMembers{UserId: int(chat.Owner), ChatId: chat.ChatId, Owner: 1, Role: parse.Role, DateCreated: time.Now(), DateUpdated: time.Now()})

	for i := 0; i < len(body.Members); i++ {
		models.DB.Model(&models.User{}).Where("ID = ?", body.Members[i]).First(&checkUser)

		if chat.Owner == uint(body.Members[i]) {
			c.JSON(400, ErrorMsg(27, language.Language("owner_self_invite")))
			return
		}
		chatMembers.ChatId = chat.ChatId
		chatMembers.Owner = 0
		chatMembers.DateCreated = time.Now()
		chatMembers.DateUpdated = time.Now()
		chatMembers.UserId = body.Members[i]
		chatMembers.Role = checkUser.Role
		models.DB.Model(&chatMembers).Create(&chatMembers)
	}
	consumers.SendJSON(models.RMQMessage{SessionLost: false, ChatId: chat.ChatId, ChatDeleted: false, ChatCreated: true})

	c.JSON(200, gin.H{
		"success": true,
		"chat_id": code,
		"owner":   user.ID,
		"name":    body.Name,
		"members": body.Members,
	})
}

func DeleteChat(c *gin.Context) {
	var body models.Chat
	var chat models.Chat
	var user models.User

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if cookie == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	err = c.ShouldBindJSON(&body)
	if body.ChatId == "" {
		c.JSON(400, ErrorMsg(21, language.Language("invalid_chat_id")))
		return
	}

	models.DB.Model(&user).Where("email = ?", parse.Subject).First(&user)
	models.DB.Model(&chat).Where("chat_id = ?", body.ChatId).First(&chat)

	if chat.ChatId == "" {
		c.JSON(400, ErrorMsg(25, language.Language("delete_chat_impossible")))
		return
	}

	if chat.Owner != user.ID {
		c.JSON(400, ErrorMsg(25, language.Language("delete_chat_impossible")))
		return
	}

	models.DB.Delete(chat)
	consumers.SendJSON(models.RMQMessage{SessionLost: false, ChatId: chat.ChatId, ChatDeleted: true, ChatCreated: false})
	c.JSON(200, gin.H{
		"success": true,
		"message": language.Language("delete_chat_successful"),
	})

}

func ListChat(c *gin.Context) {
	var user models.User
	var chats []models.ChatMembers

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if cookie == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	models.DB.Model(&user).Where("email = ?", parse.Subject).First(&user)
	models.DB.Model(&models.ChatMembers{}).Where("user_id = ?", user.ID).Order("date_updated DESC").Find(&chats)

	var chatsId []string
	for _, chatMember := range chats {
		chatsId = append(chatsId, chatMember.ChatId)
	}

	var chatList []models.ChatInfo
	var chatInfo models.ChatInfo
	var chat models.Chat
	for i := 0; i < len(chatsId); i++ {
		models.DB.Where("chat_id = ?", chatsId[i]).First(&chat)
		chatInfo.Name = chat.Name
		chatInfo.ChatId = chatsId[i]
		chatInfo.Owner = int(user.ID)
		chatInfo.Members = models.UsersOfChat(chatsId[i])
		chatInfo.Created = chat.Created
		chatInfo.Updated = chat.Updated

		chatList = append(chatList, chatInfo)
		chat.ID = 0
	}

	c.JSON(200, chatList)
}

func ChatInfo(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if cookie == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	var user models.User

	models.DB.Where("email = ?", parse.Subject).First(&user)

	chatId := c.Param("id")

	var chatMembers models.ChatMembers
	var chat models.Chat

	models.DB.Where("user_id = ? AND chat_id = ?", user.ID, chatId).First(&chatMembers)

	if chatMembers.ChatId == "" || chatMembers.UserId == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}
	models.DB.Where("chat_id = ?", chatId).First(&chat)

	chatInfo := models.ChatInfo{Name: chat.Name, ChatId: chatId, Members: models.UsersOfChat(chatId), Owner: int(chat.Owner), Created: chat.Created, Updated: chat.Updated}

	c.JSON(200, chatInfo)

}
