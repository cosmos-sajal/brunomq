package main

import (
	"fmt"
	"log"
	"net/http"
	"queue/config"
	"queue/message"
	"sync"
)

type Queue struct {
	messages           []*message.Message
	retryMessages      []*message.Message
	deadLetterMessages []*message.Message
	consumerOffset     map[string]int
	name               string
	mu                 sync.Mutex
}

func NewQueue(name string) *Queue {
	return &Queue{
		name: name,
	}
}

var queuesMap = make(map[string]*Queue)

func (q *Queue) AddMessage(m *message.Message) {
	q.messages = append(q.messages, m)
}

func (q *Queue) AddRetryMessage(m *message.Message) {
	q.retryMessages = append(q.retryMessages, m)
}

func (q *Queue) AddDeadLetterMessage(m *message.Message) {
	q.deadLetterMessages = append(q.deadLetterMessages, m)
}

func (q *Queue) GetMessage(offset int) *message.Message {
	if offset < len(q.messages) {
		return q.messages[offset]
	}

	return nil
}

func (q *Queue) SetOffset(consumer string, offset int) {
	if q.consumerOffset == nil {
		q.consumerOffset = make(map[string]int)
	}

	q.consumerOffset[consumer] = offset
}

func (q *Queue) GetOffset(consumer string) int {
	return q.consumerOffset[consumer]
}

func (q *Queue) GetName() string {
	return q.name
}

func handlePush(w http.ResponseWriter, r *http.Request) {
	messageContent := r.URL.Query().Get("message_content")
	queueName := r.URL.Query().Get("queue_name")
	if messageContent == "" {
		http.Error(w, "Missing message", http.StatusBadRequest)
		return
	} else if queueName == "" {
		http.Error(w, "Missing queue name", http.StatusBadRequest)
		return
	}

	if q, ok := queuesMap[queueName]; ok {
		newMessage := message.NewMessage(messageContent)
		q.AddMessage(newMessage)
		fmt.Fprintf(w, "Message added to queue %s", queueName)
		fmt.Println("Message added to queue:", queueName)
		return
	} else {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}
}

func handlePop(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Query().Get("queue_name")
	if queueName == "" {
		http.Error(w, "Missing queue name", http.StatusBadRequest)
		return
	}

	if q, ok := queuesMap[queueName]; ok {
		q.mu.Lock()
		defer q.mu.Unlock()

		message := q.GetMessage(0)
		if message == nil {
			http.Error(w, "No messages in queue", http.StatusNotFound)
			return
		}

		q.messages = q.messages[1:]
		fmt.Fprintf(w, "Message: %s", message.GetContent())
		return
	} else {
		http.Error(w, "Queue not found", http.StatusBadRequest)
		return
	}
}

func init() {
	fmt.Println("Queue config:", config.QueueConfig)
	for _, q := range config.QueueConfig["queues"].([]string) {
		createdQueue := NewQueue(q)
		queuesMap[q] = createdQueue
	}

	fmt.Println("Queues created:", config.QueueConfig["queues"])
}

func main() {
	http.HandleFunc("/push", handlePush)
	http.HandleFunc("/pop", handlePop)

	log.Println("Starting queue server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
