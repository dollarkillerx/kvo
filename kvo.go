package kvo

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Kvo struct {
	sync.Mutex
	kvoMap map[string][]*Channel
}

func New() *Kvo {
	return &Kvo{
		kvoMap: map[string][]*Channel{},
	}
}

type Channel struct {
	kvo         *Kvo
	Channel     chan interface{}
	id          int    // 当前ID
	subjectName string // 订阅主题
}

// Subscription 订阅
// params: 主题名称
// returns: KvoChannel,error
func (k *Kvo) Subscription(subjectName string) (*Channel, error) {
	k.Lock()
	defer k.Unlock()
	c := make(chan interface{}, 10)
	kvoChan := &Channel{
		Channel:     c,
		id:          0,
		kvo:         k,
		subjectName: subjectName,
	}
	channels := k.kvoMap[subjectName]
	kvoChan.id = len(channels)
	k.kvoMap[subjectName] = append(k.kvoMap[subjectName], kvoChan)
	return kvoChan, nil
}

// Unsubscribe 退订
// params: 主题名称
// returns: error
func (k *Kvo) Unsubscribe(subjectName string) error {
	k.Lock()
	defer k.Unlock()

	delete(k.kvoMap, subjectName)
	return nil
}

// Publish 发布
// params: 主题名称,message
// returns: error
func (k *Kvo) Publish(subjectName string, msg interface{}) error {
	for i := 0; ; i++ {
		// 当主题不存在是尝试插入
		err := k.publish(subjectName, msg)
		if err == nil {
			return nil
		}

		if i > 100 {
			return err
		}
		time.Sleep(time.Second * 3)
	}
}

func (k *Kvo) publish(subjectName string, msg interface{}) error {
	// fix: fatal error: all goroutines are asleep - deadlock!
	k.Lock()
	channels, ex := k.kvoMap[subjectName]
	if !ex {
		k.Unlock()
		return errors.New("not subject")
	}
	k.Unlock()

	for _, v := range channels {
		fmt.Println(len(v.Channel))
		v.Channel <- msg
	}
	return nil
}

// Unsubscribe 退订
// returns: error
func (k *Channel) Unsubscribe() error {
	k.kvo.Lock()
	defer k.kvo.Unlock()

	channels, ex := k.kvo.kvoMap[k.subjectName]
	if !ex {
		return errors.New("not exist")
	}

	if len(channels)-1 == k.id {
		channels = append(channels[:k.id])
	} else {
		channels = append(channels[:k.id], channels[k.id+1:]...)
	}
	k.kvo.kvoMap[k.subjectName] = channels
	return nil
}

func (k *Channel) Chan() chan interface{} {
	return k.Channel
}
