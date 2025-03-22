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

type DalleImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

type ImageData struct {
	URL string `json:"url"`
}
