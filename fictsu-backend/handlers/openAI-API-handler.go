package handlers

import (
	"io"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"

	models "github.com/Fictsu/Fictsu/models"
	configs "github.com/Fictsu/Fictsu/configs"
)

func AddHeader(request *http.Request) {
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer " + configs.OpenAIKey)
	request.Header.Add("OpenAI-Organization", configs.OpenAIOrgID)
	request.Header.Add("OpenAI-Project", configs.OpenAIProjID)
}

func OpenAIGetText(ctx *gin.Context) {
	requestBody := models.OpenAIRequestBodyText{}
	if err_req := ctx.ShouldBindJSON(&requestBody); err_req != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	// Prepare OpenAI request payload
	URL := "https://api.openai.com/v1/chat/completions"
	openAIRequest := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{"role": "user", "content": requestBody.Message},
		},
	}

	// Convert request body to JSON
	JSONBody, err := json.Marshal(openAIRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to encode request"})
		return
	}

	// Create new HTTP request
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(JSONBody))
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create request"})
		return
	}

	AddHeader(request)

	// Send request to OpenAI
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to send request"})
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": fmt.Sprintf("API error: %s", string(body))})
		return
	}

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to read response"})
		return
	}

	// Unmarshal OpenAI response
	responseBody := models.OpenAIResponseBody{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to decode response"})
		return
	}

	// Check if the response has choices
	if len(responseBody.Choices) == 0 {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "No choices returned from OpenAI"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Received_Message": responseBody.Choices[0].Message.Content})
}

func OpenAIGetTextToImage(ctx *gin.Context) {
	requestBody := models.OpenAIRequestBodyTextToImage{}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	// OpenAI API endpoint
	URL := "https://api.openai.com/v1/images/generations"
	openAIRequest := map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": requestBody.Message,
		"n":      1,
		"size":   requestBody.Size,
	}

	// Convert request body to JSON
	JSONBody, err := json.Marshal(openAIRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to encode request"})
		return
	}

	// Create new HTTP request
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(JSONBody))
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create request"})
		return
	}

	AddHeader(request)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to send request"})
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": fmt.Sprintf("API error: %s", string(body))})
		return
	}

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to read response"})
		return
	}

	fmt.Println("Response Body: ", string(body))
}
