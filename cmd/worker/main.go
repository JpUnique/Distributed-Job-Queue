package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Worker started")

	for {
		log.Println("Worker heartbeat...")
		time.Sleep(5 * time.Second)
	}
}

