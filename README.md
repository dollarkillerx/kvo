# kvo 消息总线 Message BUS 
受Swift Kvo 启发  的Golang  消息总线

### install
```shell
go get github.com/dollarkillerx/kvo
```

### 🌰栗子 Example
``` 
	sub1, err := Kvo.Subscription("aa")  // 订阅消息
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		data := <-sub1.Channel   // 获取消息
		log.Println(data)
		sub1.Unsubscribe()          // 解除订阅
	}()

	Kvo.Publish("aa","this is aa")    // 发布消息

	Kvo.Unsubscribe("aa")         // 解除订阅

	time.Sleep(time.Second)
```



``` 
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
```
