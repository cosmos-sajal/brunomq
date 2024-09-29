package config

var QueueConfig = map[string]interface{}{
	"max_retry": 3,
	"queues": []string{
		"my-queue-1",
		"my-queue-2",
		"my-queue-3",
	},
}
