package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	done := make(chan bool)

	for {
		var userInput string
		fmt.Println("##########################")
		fmt.Println("Choose option: ")
		fmt.Println("1. Push to queue")
		fmt.Println("2. Exit")
		fmt.Println("##########################")
		fmt.Scanln(&userInput)

		switch userInput {
		case "1":
			fmt.Println("Enter message to push to queue:")
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			msg = msg[:len(msg)-1]
			fmt.Println("Enter the queue to push to:")
			var queueName string
			fmt.Scanln(&queueName)
			resp, err := http.Get("http://localhost:8080/push?queue_name=" + queueName + "&message_content=" + msg)
			if err != nil {
				fmt.Println("Error pushing message to queue:", err)
				continue
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error pushing message to queue:", string(respBody))
				continue
			}
		}

		if userInput == "2" {
			fmt.Println("Exiting...")
			done <- true
			break
		}
	}

	<-done
}
