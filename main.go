package main

import (
	// "flag"
	"fmt"
	"time"
	"log"
	"context"

	"github.com/segmentio/kafka-go"
)
func main(){
	topic := "test-topic"
	partition :=0
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	fmt.Println("1")
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
_, err = conn.WriteMessages(
    kafka.Message{Value: []byte("one!")},
    kafka.Message{Value: []byte("two!")},
    kafka.Message{Value: []byte("three!")},
)
if err != nil {
    log.Fatal("failed to write messages:", err)
}
fmt.Println("2")

if err := conn.Close(); err != nil {
    log.Fatal("failed to close writer:", err)
}
}