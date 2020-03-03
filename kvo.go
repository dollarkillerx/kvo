package kvo

import (
	"errors"
	"sync"
)

type KvoInterface interface {
	Subscription(subjectName string) (*KvoChannel, error)
	Publish(subjectName string, msg interface{}) error
	Unsubscribe(subjectName string) error
}

type kvo struct {
	sync.Mutex
	kvoMap map[string][]*KvoChannel
}

type KvoChannel struct {
	kvo         *kvo
	Channel     chan interface{}
	id          int
	subjectName string
}

var Kvo *kvo

func init() {
	Kvo = &kvo{
		kvoMap: map[string][]*KvoChannel{},
	}
}

// 订阅
// params: 主题名称
// returns: KvoChannel,error
func (k *kvo) Subscription(subjectName string) (*KvoChannel, error) {
	k.Lock()
	defer k.Unlock()
	c := make(chan interface{}, 100)
	kvoChan := &KvoChannel{
		Channel:     c,
		id:          0,
		kvo:         k,
		subjectName: subjectName,
	}
	channels, bool := k.kvoMap[subjectName]
	if !bool {
		k.kvoMap[subjectName] = make([]*KvoChannel, 0)
		k.kvoMap[subjectName] = append(k.kvoMap[subjectName], kvoChan)
	} else {
		kvoChan.id = len(channels)
		k.kvoMap[subjectName] = append(k.kvoMap[subjectName], kvoChan)
	}
	return kvoChan, nil
}

// 退订
// params: 主题名称
// returns: error
func (k *kvo) Unsubscribe(subjectName string) error {
	k.Lock()
	defer k.Unlock()
	_, bool := k.kvoMap[subjectName]
	if !bool {
		return nil
	}
	delete(k.kvoMap, subjectName)
	return nil
}

// 发布
// params: 主题名称,message
// returns: error
func (k *kvo) Publish(subjectName string, msg interface{}) error {
	k.Lock()
	defer k.Unlock()
	chans, bool := k.kvoMap[subjectName]
	if !bool {
		return errors.New("does not exist")
	}
	for _, v := range chans {
		v.Channel <- msg
	}
	return nil
}

// 退订
// returns: error
func (k *KvoChannel) Unsubscribe() error {
	k.kvo.Lock()
	defer k.kvo.Unlock()
	chans, bool := k.kvo.kvoMap[k.subjectName]
	if !bool {
		return errors.New("not exist")
	}

	if len(chans)-1 == k.id {
		chans = append(chans[:k.id])
	} else {
		chans = append(chans[:k.id], chans[k.id+1:]...)
	}
	k.kvo.kvoMap[k.subjectName] = chans
	return nil
}

func (k *KvoChannel) Chan() chan interface{} {
	return k.Channel
}