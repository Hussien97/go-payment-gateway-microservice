package main

import (
	"log"
)

func main() {
	log.Println("Starting Gateway B...")
	StartKafkaConsumer("transactions.soap")
	select {}
}
