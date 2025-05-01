package main

import (
	server "consumer/internal/nats_server"
	"consumer/internal/store/sqlite"
	"log"
)

func main() {
	store := sqlite.NewStore("./storage/log.db")
	s := server.NewNatsServer(store, 50052)

	_, err := s.NC.Subscribe("order.created", s.HandleCreateOrder)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Subscribed to order.created events")
	select {} // keep running
}
