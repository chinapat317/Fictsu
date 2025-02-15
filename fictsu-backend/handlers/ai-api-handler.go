package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Fictsu/Fictsu/models"
	"github.com/Fictsu/Fictsu/secret"
	"github.com/gin-gonic/gin"
)

func AddHeader(req *http.Request) {
	req.Header.Add("Content-Type", secret.OPENAI_APP_HEADER)
	req.Header.Add("Authorization", secret.OPENAI_KEY)
	req.Header.Add("OpenAI-Organization", secret.OPENAI_ORG_ID)
	req.Header.Add("OpenAI-Project", secret.OPENAI_PROJ_ID)
}

func OpenAIGetText(ctx *gin.Context) {
	var reqBody models.OpenAIReqBodyText
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	url := "https://api.openai.com/v1/chat/completions"
	var requestBody = map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{"role": "user", "content": reqBody.Message},
		}}
	jsonBody, err := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	AddHeader(req)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	var resBody models.OpenAIResBody
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"received_message": resBody.Choices[0].Message.Content})
}

func OpenAITextToPic(ctx *gin.Context) {
	var reqBody models.OpenAIReqBodyTextToImg
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	url := "https://api.openai.com/v1/images/generations"
	var requestBody = map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": reqBody.Message,
		"n":      1,
		"size":   reqBody.Size,
	}
	jsonBody, err := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	AddHeader(req)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}
