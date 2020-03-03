package kvo

import (
	"log"
	"testing"
	"time"
)

func TestKvo(t *testing.T) {
	sub1, err := Kvo.Subscription("aa")
	if err != nil {
		log.Fatalln(err)
	}

	go test1(sub1)

	sub2, err := Kvo.Subscription("bb")
	if err != nil {
		log.Fatalln(err)
	}

	go test2(sub2)

	go func() {
		for {
			time.Sleep(time.Microsecond * 200)
			Kvo.Publish("aa","this is a")
		}
	}()

	for {
		time.Sleep(time.Microsecond * 200)
		Kvo.Publish("bb","this is b")
	}

}

func test1(kvo *KvoChannel) {
	for {
		select {
		case val := <- kvo.Chan():
			log.Println(val)
		}
	}
}

func test2(kvo *KvoChannel) {
	for {
		select {
		case val := <- kvo.Chan():
			log.Println(val)
		}
	}
}

func TestKvo2(t *testing.T) {
	sub1, err := Kvo.Subscription("aa")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		data := <-sub1.Channel
		log.Println(data)
		sub1.Unsubscribe()
	}()

	Kvo.Publish("aa","this is aa")

	Kvo.Unsubscribe("aa")

	time.Sleep(time.Second)
}