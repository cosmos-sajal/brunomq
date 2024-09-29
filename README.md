# BrunoMQ

### How to run?
#### Run Broker queue
- cd to `queue`
- Can change the queues in config.go file.
- Run command `go run queue.go`

#### Run Consumer
- cd to `consumer`
- Run command `go run consumer.go <queue_name> <num_of_consumers>`
- Example `go run consumer.go my-queue-1 3`

#### Run Producer
- cd to `producer`
- Run command `go run producer.go`
- Follow the prompt to push messages to queue.
