package chat

type MessageType int

const (
	Default MessageType = iota
	TypingNotification
)

// Message represents a WebSocket message
type Message struct {
	Type MessageType `json:"type"`
	User string      `json:"user"`
	Text string      `json:"text"`
}
