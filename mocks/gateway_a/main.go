package main

import (
	"log"
)

func main() {
	log.Println("Starting Gateway A...")
	StartKafkaConsumer("transactions.json")
	select {}
}
