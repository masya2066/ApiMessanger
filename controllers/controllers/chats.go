package controllers

import (
	"ApiMessenger/consumers"
	"ApiMessenger/controllers"
	"ApiMessenger/language"
	"ApiMessenger/middlewares"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type ChatBody struct {
	Name    string `json:"name"`
	Members []int  `json:"members"`
}

func NewChat(c *gin.Context) {
	var chat models.Chat
	var user models.User

	isAuth, parse := middlewares.IsAuthorized(c)

	if !isAuth || parse.Subject == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
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

	models.DB.Model(&models.User{}).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	code := utils.GenerateId()
	chat.Name = body.Name
	chat.Owner = user.ID
	chat.ChatId = code

	var chatMembers models.ChatMembers
	var checkUser models.User
	var invalidUsers []int

	for i := 0; i < len(body.Members); i++ {
		models.DB.Model(&models.User{}).Where("ID = ?", body.Members[i]).First(&checkUser)

		if checkUser.ID == 0 {
			invalidUsers = append(invalidUsers, body.Members[i])
		}
		checkUser.ID = 0
	}

	if len(invalidUsers) != 0 {
		c.JSON(400, ErrorMsg(26, language.Language("error_invite_member_to_chat")+utils.IntSliceToString(invalidUsers)))
		return
	}

	for i := 0; i < len(body.Members); i++ {
		models.DB.Model(&models.User{}).Where("ID = ?", body.Members[i]).First(&checkUser)
		if chat.Owner == uint(body.Members[i]) {
			c.JSON(400, ErrorMsg(27, language.Language("owner_self_invite")))
			return
		}
	}

	chat.Created, chat.Updated = time.Now().UTC().Format(os.Getenv("DATE_FORMAT")), time.Now().UTC().Format(os.Getenv("DATE_FORMAT"))
	chat.Phrase = utils.GenerateRandomSecretPhrase()

	models.DB.Create(&chat)
	models.DB.Create(&models.ChatMembers{UserId: int(chat.Owner), ChatId: chat.ChatId, Owner: true, Role: parse.Role, DateCreated: time.Now().UTC().Format(os.Getenv("DATE_FORMAT")), DateUpdated: time.Now().UTC().Format(os.Getenv("DATE_FORMAT"))})
	consumers.SendJSON(models.RMQMessage{SessionLost: false, ChatId: chat.ChatId, ChatDeleted: false, ChatCreated: true})

	for i := 0; i < len(body.Members); i++ {
		chatMembers.ChatId = chat.ChatId
		chatMembers.Owner = false
		chatMembers.DateCreated, chatMembers.DateUpdated = time.Now().UTC().Format(os.Getenv("DATE_FORMAT")), time.Now().UTC().Format(os.Getenv("DATE_FORMAT"))
		chatMembers.UserId = body.Members[i]
		chatMembers.Role = checkUser.Role
		models.DB.Model(&chatMembers).Create(&chatMembers)
	}

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

	isAuth, parse := middlewares.IsAuthorized(c)

	if !isAuth || parse.Subject == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	_ = c.ShouldBindJSON(&body)
	if body.ChatId == "" {
		c.JSON(400, ErrorMsg(21, language.Language("invalid_chat_id")))
		return
	}

	models.DB.Model(&user).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

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

	isAuth, parse := middlewares.IsAuthorized(c)

	if !isAuth || parse.Subject == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	models.DB.Model(&user).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	models.DB.Model(&models.ChatMembers{}).Where("user_id = ?", user.ID).Order("date_updated DESC").Find(&chats)

	if len(chats) == 0 {
		c.JSON(200, ErrorMsg(24, language.Language("chats_empty")))
		return
	}

	var chatsId []string
	for _, chatMember := range chats {
		chatsId = append(chatsId, chatMember.ChatId)
	}

	var chatList []models.ChatInfo
	var chatInfo models.ChatInfo
	var chat models.Chat

	for i := 0; i < len(chatsId); i++ {
		models.DB.Where("chat_id = ?", chatsId[i]).First(&chat)
		mes, err := controllers.GetMessagesOfChat(chat, 1)
		if err != nil {
			fmt.Println()
		}
		chatInfo.Name = chat.Name
		chatInfo.ChatId = chatsId[i]
		chatInfo.Owner = int(user.ID)
		chatInfo.Members = models.UsersOfChat(chatsId[i])
		chatInfo.Messages = mes
		chatInfo.Created = chat.Created
		chatInfo.Updated = chat.Updated
		chatList = append(chatList, chatInfo)
		chatInfo.Messages = []models.Message{}
		chat.ID = 0
	}

	if len(chatList) == 0 {
		panic("Error work with chats")
		return
	}

	c.JSON(200, chatList)
}

func ChatInfo(c *gin.Context) {
	isAuth, parse := middlewares.IsAuthorized(c)

	if !isAuth || parse.Subject == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	var user models.User

	models.DB.Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	chatId := c.Param("id")

	var chatMembers models.ChatMembers
	var chat models.Chat

	models.DB.Where("user_id = ? AND chat_id = ?", user.ID, chatId).First(&chatMembers)

	if chatMembers.ChatId == "" || chatMembers.UserId == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}
	models.DB.Where("chat_id = ?", chatId).First(&chat)

	mes, err := controllers.GetMessagesOfChat(chat, 1)
	if err != nil {
		fmt.Println()
	}

	chatInfo := models.ChatInfo{Name: chat.Name, ChatId: chatId, Members: models.UsersOfChat(chatId), Messages: mes, Owner: int(chat.Owner), Created: chat.Created, Updated: chat.Updated}

	c.JSON(200, chatInfo)
}
