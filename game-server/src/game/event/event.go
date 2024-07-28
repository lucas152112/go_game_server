package event

import (
	"github.com/golang/glog"
	"sync"
)

const (
	SyncEvent  = 1  //同步
	AsyncEvent = 2  //异步
	PipeEvent  = 3
)

type HandlerOne   func(data interface{} )
type HandlerTwo   func(data interface{} ) error
type HandlerThree func(data interface{} ) (interface{},error)


type Event struct {
	EventType  int
	handler    interface{}
	data       interface{}
}


func NewSyncEvent( data interface{} , handler interface{} ) Event  {
	return Event{
		EventType: SyncEvent,
		handler:   handler,
		data:      data,
	}
}

func NewAsyncEvent( data interface{} ,handler interface{} ) Event  {
	return Event{
		EventType: AsyncEvent,
		handler:   handler,
		data:      data,
	}
}


func NewEvents() *Events {
	return &Events{
		handlers: make(map[string][]Event),
		RWMutex:  sync.RWMutex{},
	}
}

type Events struct {
	handlers   map[string] []Event
	sync.RWMutex
}

func (this *Events)  EventRegister( eventName string , event ...Event )  {
	for i:=0;i<len(event);i++ {
		item:=event[i]
		if _, ok := this.handlers[eventName]; !ok {
			this.Lock()
			this.handlers[eventName] = []Event{}
			this.Unlock()
		}
		this.Lock()
		this.handlers[eventName] = append(this.handlers[eventName], item)
		this.Unlock()
	}
}

func (this *Events) EventTrigger( eventName string )  {
	this.RLock()
	eventList,ok := this.handlers[eventName]
	this.RUnlock()
	if !ok{
		return
	}
	glog.Info("Event Trigger EventName ",eventName," EventLen ",len(eventList))

	for i:=0;i<len(eventList);i++ {
		event:= eventList[i]
		switch event.EventType {
		case SyncEvent:
			this.syncEventTrigger(event)
		case AsyncEvent:
			glog.Info("Event Trigger Async")
			this.asyncEventTrigger(event)
		case PipeEvent:
			this.pipeEventTrigger(event)
		default:
			glog.Info("cant find EventType ",event.EventType)
		}
	}
	return
}

func (this *Events) syncEventTrigger( event Event ) {
	switch event.handler.(type) {
	case HandlerOne:
		event.handler.(HandlerOne)( event.data )
	}
}

func (this *Events) asyncEventTrigger( event Event ) {
	switch event.handler.(type) {
	case func( data interface{}):
		go event.handler.(func(data interface{}))( event.data )
		glog.Info("AsyncEventTrigger ... ")
	//case *PipeLine :
	//	go event.handler.(*PipeLine).Run(event.data)
	default:
		glog.Info("EventTrigger Cant Found type ")
	}
}

func (this *Events) pipeEventTrigger( event Event )  {
	switch event.handler.(type) {
	case *PipeLine:
		event.handler.(*PipeLine).Run(event.data)
	}
}
