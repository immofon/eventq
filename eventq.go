package eventq

import (
	"container/list"
	"fmt"
)

type EventFunc func(ud interface{})

type Event struct {
	Tick int
	Func EventFunc
}

type Queue struct {
	tick int
	data *list.List
}

func New() *Queue {
	return &Queue{
		tick: 0,
		data: list.New(),
	}
}

func (eq *Queue) should_insert_front(tick int) bool {
	v := eq.data.Front()
	if v == nil {
		return true
	}
	return v.Value.(Event).Tick > tick
}

func (eq *Queue) find_last(tick int) *list.Element {
	ret := eq.data.Front()
	if ret.Next() == nil {
		return ret
	}

	for {
		v := ret.Next()
		if v == nil {
			return ret
		}
		if v.Value.(Event).Tick > tick {
			return ret
		}
		ret = ret.Next()
	}
	return nil
}

func (eq *Queue) Tick() int {
	return eq.tick
}
func (eq *Queue) Add(afterTick int, fn EventFunc) {
	if afterTick <= 0 {
		panic("afterTick MUST great than 0")
	}

	e := Event{
		Tick: eq.tick + afterTick,
		Func: fn,
	}
	if eq.should_insert_front(e.Tick) {
		eq.data.PushFront(e)
	} else {
		eq.data.InsertAfter(e, eq.find_last(e.Tick))
	}
}

func (eq *Queue) Clean() {
	eq.data = list.New()
}

func (eq *Queue) IsEmpty() bool {
	return eq.data.Front() == nil
}

func (eq *Queue) NextEventTick() int {
	mark := eq.data.Front()
	if mark == nil {
		return 0
	}
	return mark.Value.(Event).Tick
}

func (eq *Queue) Run(ud interface{}) {
	mark := eq.data.Front()
	if mark == nil {
		return
	}
	eq.tick = eq.NextEventTick()

	for mark != nil {
		if e := mark.Value.(Event); e.Tick == eq.tick {
			if e.Func != nil {
				e.Func(ud)
			}
			eq.data.Remove(mark)
		} else {
			return
		}
		mark = mark.Next()
	}
}

func (eq *Queue) RunUntilEmpty(ud interface{}) {
	for !eq.IsEmpty() {
		eq.Run(ud)
	}
}

func (eq *Queue) Print() {
	fmt.Printf("%d: EventQueue [", eq.tick)
	for e := eq.data.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		t := e.Value.(Event).Tick
		fmt.Printf("%d", t)
		if e.Next() != nil {
			fmt.Print(", ")
		}
	}
	fmt.Println("]")
}
