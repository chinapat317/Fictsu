package models

type OpenAIReqBodyText struct {
	Message string `json:"message"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResBody struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}

type OpenAIReqBodyTextToImg struct {
	Message string `json:"message"`
	Size    string `json:"size"`
}
