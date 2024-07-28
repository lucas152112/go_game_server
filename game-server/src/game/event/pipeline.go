package event



/**
 使用PipeLine
 */

func NewPipeEvent( data interface{},handler ...HandlerThree) Event  {
	return Event{
		EventType : PipeEvent,
		handler   : NewPipeLine(handler...),
		data:data,
	}
}

func NewPipeLine(handler ...HandlerThree ) *PipeLine  {
	pipe := &PipeLine{
		handlers: []HandlerThree{},
	}
	for i:=0;i<=len(handler);i++{
		pipe.handlers = append(pipe.handlers,handler[i])
	}
	return pipe
}

type PipeLine struct {
	handlers  []HandlerThree    //使用第三种
}

func (this *PipeLine) Run( data interface{} ) ( interface{} , error )  {
	result := data
	var err error
	for i:=0;i<len(this.handlers);i++{
		handler := this.handlers[i]
		result,err =handler( result )
		if err != nil{
			return result,err
		}
	}
	return result,err
}

