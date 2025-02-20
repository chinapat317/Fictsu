package models

type OpenAIRequestBodyText struct {
    Message string `json:"message"`
}

type OpenAIRequestBodyTextToImage struct {
    Message string `json:"message"`
    Size    string `json:"size"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponseBody struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}
