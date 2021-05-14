package main

import (
	"fmt"
	"github.com/loophole-labs/frisbee"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

const PING = uint32(1)
const PONG = uint32(2)

func handlePong(incomingMessage frisbee.Message, incomingContent []byte) (outgoingMessage *frisbee.Message, outgoingContent []byte, action frisbee.Action) {
	if incomingMessage.ContentLength > 0 {
		log.Printf("Client Received Message: %s", string(incomingContent))
	}
	return
}

func main() {
	router := make(frisbee.ClientRouter)
	router[PONG] = handlePong
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	c := frisbee.NewClient("127.0.0.1:8192", router)
	err := c.Connect()
	if err != nil {
		panic(err)
	}

	go func() {
		i := 0
		for {
			message := []byte(fmt.Sprintf("ECHO MESSAGE: %d", i))
			err := c.Write(&frisbee.Message{
				To:            0,
				From:          0,
				Id:            uint32(i),
				Operation:     PING,
				ContentLength: uint64(len(message)),
			}, &message)
			if err != nil {
				panic(err)
			}
			i++
			time.Sleep(time.Second)
		}
	}()

	<-exit
	err = c.Close()
	if err != nil {
		panic(err)
	}
}
