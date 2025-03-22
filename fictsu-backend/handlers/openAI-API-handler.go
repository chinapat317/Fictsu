package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	configs "github.com/Fictsu/Fictsu/configs"
	db "github.com/Fictsu/Fictsu/database"
	models "github.com/Fictsu/Fictsu/models"
)

// prompt
const (
	INTRO_TEXT string = "Please generate story about: "
	OUTRO_TEXT string = " Only provide structure of story not whole story."
	//Generate character pic prompt
	INTRO_CH string = "Please generate character follow this prompt in T-pose so image can be use as reference for future generation. The prompt is: '"
)

func AddHeader(request *http.Request) {
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+configs.OpenAIKey)
	request.Header.Add("OpenAI-Organization", configs.OpenAIOrgID)
	request.Header.Add("OpenAI-Project", configs.OpenAIProjID)
}

func OpenAIGenStruc(ctx *gin.Context) {
	requestBody := models.OpenAIRequestBodyText{}
	if err_req := ctx.ShouldBindJSON(&requestBody); err_req != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	// Prepare OpenAI request payload
	var prompt_message = INTRO_TEXT + "'" + requestBody.Message + "'" + OUTRO_TEXT
	URL := "https://api.openai.com/v1/chat/completions"
	openAIRequest := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{"role": "user", "content": prompt_message},
		},
	}

	// Convert request body to JSON
	JSONBody, err_mar := json.Marshal(openAIRequest)
	if err_mar != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to encode request"})
		return
	}

	// Create new HTTP request
	request, err_new_req := http.NewRequest("POST", URL, bytes.NewBuffer(JSONBody))
	if err_new_req != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create request"})
		return
	}

	AddHeader(request)

	// Send request to OpenAI
	client := &http.Client{}
	response, err_res := client.Do(request)
	if err_res != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to send request"})
		return
	}

	defer response.Body.Close()

	// Read response body
	body, err_res_body := io.ReadAll(response.Body)
	if err_res_body != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to read response"})
		return
	}

	fmt.Println("Response Body: ", string(body))

	// Unmarshal OpenAI response
	responseBody := models.OpenAIResponseBody{}
	if err_unmar := json.Unmarshal(body, &responseBody); err_unmar != nil {
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

func OpenAICreateChar(ctx *gin.Context) {
	requestBody := models.OpenAIRequestBodyTextToImage{}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	URL := "https://api.openai.com/v1/images/generations"
	var prompt_message = INTRO_CH + "'" + requestBody.Message + "'"
	openAIRequest := map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": prompt_message,
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

	body, err_res_body := io.ReadAll(response.Body)
	if err_res_body != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to read response"})
		return
	}
	fmt.Println("Response Body: ", string(body))

	responseBody := models.DalleImageResponse{}
	if err_unmar := json.Unmarshal(body, &responseBody); err_unmar != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Failed to decode response"})
		return
	}

	// Check if the response has Data
	if len(responseBody.Data) == 0 {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "No choices returned from OpenAI"})
		return
	}
	var count int
	db.DB.QueryRow("SELECT COUNT(*) FROM Character;").Scan(&count)
	imageURL := responseBody.Data[0].URL
	filePath := configs.CharPath + strconv.Itoa(count+1) + ".png"
	err = downloadImage(imageURL, filePath)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}
	fmt.Println("Image saved successfully to", filePath)
}

func downloadImage(url, filePath string) error {
	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the image data to the file
	_, err = io.Copy(out, resp.Body)
	return err
}
