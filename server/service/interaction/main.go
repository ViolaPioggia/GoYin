package main

import (
	interaction "GoYin/server/kitex_gen/interaction/interactionserver"
	"log"
)

func main() {
	svr := interaction.NewServer(new(InteractionServerImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
