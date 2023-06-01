package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	ID 		int 	`json:"id"`
	Sender 	string 	`json:"from"`
	Subject string 	`json:"subject"`
	Date 	string 	`json:"date"`
}

func GetEmail() string {
	url := "https://www.1secmail.com/api/v1/?action=genRandomMailbox"
	response, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	email := string(body)
	return email[2: len(email)-2]
}

func Credentials(email string) (string, string) {
	parts := strings.Split(email, "@")
	login := parts[0]
	domain := parts[1]
	return login, domain
}

func GetMessages(email string) []Message {
	login, domain := Credentials(email)

	url := fmt.Sprintf("https://www.1secmail.com/api/v1/?action=getMessages&login=%s&domain=%s", login, domain)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer response.Body.Close()

	var messages []Message
	err = json.NewDecoder(response.Body).Decode(&messages)
	if err != nil {
		return nil
	}

	return messages
}

func messageHelper(message string, attachments []interface{}) string{
	if len(attachments) == 1 {
		message += "\nAttachment info:\n"
	} else {
		return message + "\n"
	}
	for _, attachment := range attachments {
		if attachmentMap, ok := attachment.(map[string]interface{}); ok {
			for key, value := range attachmentMap {
				message += fmt.Sprintf("%s: %v, ", key, value)
			}
		}
		message += "\n"
	}
	return message
}


func ReadLatestMessage(email string) string{
	login, domain := Credentials(email)

	jsonData := GetMessages(email)
	if len(jsonData) == 0{
        return "No messages found"
	}

	id := jsonData[0].ID
	urlRead := fmt.Sprintf("https://www.1secmail.com/api/v1/?action=readMessage&login=%s&domain=%s&id=%d", login, domain, id)

	response, err := http.Get(urlRead)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()

	var jsonDataRead map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&jsonDataRead)
	if err != nil {
		return err.Error()
	}

	textBody := jsonDataRead["body"].(string)
	firstLine := strings.Split(textBody, "<p class=MsoNormal>")[1]
	message := firstLine[:len(firstLine)-4]

	attachments := jsonDataRead["attachments"].([]interface{})
	message = messageHelper(message, attachments)
	return message
}

func ReadMessageById(email string, id int) string{
	login, domain := Credentials(email)

	urlRead := fmt.Sprintf("https://www.1secmail.com/api/v1/?action=readMessage&login=%s&domain=%s&id=%d", login, domain, id)
	response, err := http.Get(urlRead)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()

	var jsonDataRead map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&jsonDataRead)
	if err != nil {
		return "No message found"
	}

	textBody := jsonDataRead["body"].(string)
	firstLine := strings.Split(textBody, "<p class=MsoNormal>")[1]
	message := firstLine[:len(firstLine)-4]

	attachments := jsonDataRead["attachments"].([]interface{})
	message = messageHelper(message, attachments)
	return message
}

func readInbox(email string) string{
	var str string

	jsonData := GetMessages(email)
	if len(jsonData) == 0{
        return "No messages found"
	}

	for _, value := range jsonData {
		id := value.ID
		date := value.Date
		message := ReadMessageById(email, id)
		str += fmt.Sprintf("Message sent at %s: %s\n", date, message)
	}

	if len(str) == 0 {
		return "No messages Received"
	}

	str += "End of Messages"
	return str
}

/*
Crate a function called GetAttachments
Access it using the link 
https://www.1secmail.com/api/v1/?action=download&login=demo&domain=1secmail.com&id=639&file=iometer.pdf
*/

func GetAttachments() {

}