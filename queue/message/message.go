package message

type Message struct {
	Content string
}

func (m *Message) GetContent() string {
	return m.Content
}

func NewMessage(content string) *Message {
	return &Message{Content: content}
}
