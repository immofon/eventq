package eventq

import (
	"testing"
)

func Func(ud interface{}) {
	eq := ud.(*Queue)
	eq.Add(10, Func)
}

func TestEventQueue(t *testing.T) {
	eq := New()
	eq.Add(1, Func)
	for i := 0; i < 100; i++ {
		eq.Run(eq)
		eq.Print()
	}
}
