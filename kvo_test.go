package kvo

import (
	"fmt"
	"log"
	"testing"
)

func TestKvo(t *testing.T) {
	kvo := New()
	sub, err := kvo.Subscription("up")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for i := 0; i < 100000; i++ {
			//time.Sleep(time.Microsecond * 10)
			err := kvo.Publish("up", "hello world")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	for {
		select {
		case r, ex := <-sub.Channel:
			if !ex {
				fmt.Println("over")
				return
			}
			fmt.Println(r)
		}
	}
}

//func TestKvo(t *testing.T) {
//	sub1, err := Kvo.Subscription("aa")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	go test1(sub1)
//
//	sub2, err := Kvo.Subscription("bb")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	go test2(sub2)
//
//	go func() {
//		for {
//			time.Sleep(time.Microsecond * 200)
//			Kvo.Publish("aa","this is a")
//		}
//	}()
//
//	for {
//		time.Sleep(time.Microsecond * 200)
//		Kvo.Publish("bb","this is b")
//	}
//
//}
//
//func test1(kvo *Channel) {
//	for {
//		select {
//		case val := <- kvo.Chan():
//			log.Println(val)
//		}
//	}
//}
//
//func test2(kvo *Channel) {
//	for {
//		select {
//		case val := <- kvo.Chan():
//			log.Println(val)
//		}
//	}
//}
//
//func TestKvo2(t *testing.T) {
//	sub1, err := Kvo.Subscription("aa")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	go func() {
//		data := <-sub1.Channel
//		log.Println(data)
//		sub1.Unsubscribe()
//	}()
//
//	Kvo.Publish("aa","this is aa")
//
//	Kvo.Unsubscribe("aa")
//
//	time.Sleep(time.Second)
//}
