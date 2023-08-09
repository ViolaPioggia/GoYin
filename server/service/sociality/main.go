package main

import (
	sociality "GoYin/server/kitex_gen/sociality/socialityservice"
	"log"
)

func main() {
	svr := sociality.NewServer(new(SocialityServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
