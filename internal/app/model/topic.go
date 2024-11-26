package model

type Topic struct {
	TopicName   string `json:"topicname"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}
