package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Fictsu/Fictsu/secret"
	"github.com/gin-gonic/gin"
)

func OpenAIGetText(ctx *gin.Context) {
	url := "https://api.openai.com/v1/chat/completions"
	var jsonBody = []byte(`{
        "model": "gpt-4o-mini",
        "messages": [
            {
                "role": "user",
                "content": "Yo dude."
            }
        ]
    }`)
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
}
