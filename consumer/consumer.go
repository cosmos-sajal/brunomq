package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ListenToMessages(queueName string) {
	fmt.Println("Staring queue listener for:", queueName)

	for {
		resp, err := http.Get("http://localhost:8080/pop?queue_name=" + queueName)
		if err != nil {
			fmt.Println("Error getting message from queue:", err)
			continue
		}

		// defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			// fmt.Println("Queue is empty, waiting...")
			time.Sleep(1 * time.Second) // Queue is empty, wait before retrying
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response:", err)
			continue
		}

		fmt.Println("Consumed message:", string(body))
	}
}

func main() {
	fmt.Println("Initializing consumer...")
	done := make(chan bool)

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <queue_name> <num_of_consumers>")
		return
	}
	queueName := os.Args[1]
	numOfConsumersStr := os.Args[2]
	numOfConsumers, err := strconv.Atoi(numOfConsumersStr)
	if err != nil {
		log.Fatalf("Invalid number of consumers: %v\n", err)
	}

	for i := 0; i < numOfConsumers; i++ {
		go ListenToMessages(queueName)
	}

	<-done
}
