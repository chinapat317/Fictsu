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

func OpenAIGetText(ctx *gin.Context) {
	var clReqBody models.RequestBody
	if err := ctx.ShouldBindJSON(&clReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	message_prompt := clReqBody.Message
	url := "https://api.openai.com/v1/chat/completions"
	var requestBody = map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "user", "content": message_prompt},
		}}

	jsonBody, err := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Add("Content-Type", secret.OPENAI_APP_HEADER)
	req.Header.Add("Authorization", secret.OPENAI_KEY)
	req.Header.Add("OpenAI-Organization", secret.OPENAI_ORG_ID)
	req.Header.Add("OpenAI-Project", secret.OPENAI_PROJ_ID)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	ctx.JSON(http.StatusOK, gin.H{"received_message": message_prompt})
}
