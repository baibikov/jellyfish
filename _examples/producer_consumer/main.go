package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"jellifish/pkg/consumer"
	"jellifish/pkg/producer"
)

func main() {
	brokerProducer, err := producer.New(&producer.Config{
		Addr: ":7654",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = brokerProducer.Close()
	}()

	brokerConsumer, err := consumer.New(&consumer.Config{
		Addr: ":7654",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = brokerConsumer.Close()
	}()

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for p := range brokerConsumer.Consume(context.Background(), "test1") {
			if err := p.Err(); err != nil {
				log.Fatalln(err)
			}

			fmt.Println(string(p.Message))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 4; i++ {
			time.Sleep(time.Second)
			err = brokerProducer.Push(context.Background(), &producer.Params{
				Topic:   "test1",
				Message: []byte(fmt.Sprintf("message [%d] to topic test", rand.Intn(100))),
			})
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	wg.Wait()
}
