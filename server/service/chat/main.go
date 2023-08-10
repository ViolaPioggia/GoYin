package main

import (
	chat "GoYin/server/kitex_gen/chat/chatservice"
	"log"
)

func main() {
	svr := chat.NewServer(new(ChatServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
