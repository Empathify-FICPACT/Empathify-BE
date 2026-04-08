package request

type StartConversationRequest struct {
	TopicID string `json:"topic_id"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}