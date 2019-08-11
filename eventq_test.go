package eventq

import (
	"fmt"
	"testing"
)

func Func(tick int, ud interface{}) {
	eq := ud.(*EventQueue)
	eq.Add(10, Func)
}

type AIInitInfo struct {
	RoleID int
	BF     *BattleField
}

type Role struct {
	Name  string
	HP    int
	Atk   int
	Speed int
	AI    EventFunc // ud::AIInitInfo
}

type BattleField struct {
	Roles [2]Role
	EQ    *EventQueue
}

func TestEventQueue(t *testing.T) {
	eq := New()
	eq.Add(1, Func)
	for i := 0; i < 100; i++ {
		eq.Run(eq)
		eq.Print()
	}
}

func target_role_id(id int) int {
	if id == 0 {
		return 1
	} else {
		return 0
	}
}

func check_battle_end(bf *BattleField) bool {
	a := bf.Roles[0]
	b := bf.Roles[1]

	if a.HP > 0 && b.HP > 0 {
		return false
	}

	tick := bf.EQ.tick
	if a.HP > 0 {
		fmt.Printf("[%d] %s win\n", tick, a.Name)
	} else if b.HP > 0 {
		fmt.Printf("[%d] %s win\n", tick, b.Name)
	} else {
		fmt.Printf("[%d] 平局\n", tick)
	}

	return true
}

func AIFunc(id int) EventFunc {
	tid := target_role_id(id)
	var ai_func EventFunc
	ai_func = func(tick int, ud interface{}) {
		bf := ud.(*BattleField)
		attack := bf.Roles[id].Atk
		bf.Roles[tid].HP -= attack
		if bf.Roles[tid].HP <= 0 {
			bf.Roles[tid].HP = 0
		}
		fmt.Printf("[%d] %s攻击%s造成%d点伤害\n", tick, bf.Roles[id].Name, bf.Roles[tid].Name, attack)
		fmt.Printf("[%d] %s:HP:%d\n", tick, bf.Roles[id].Name, bf.Roles[id].HP)
		fmt.Printf("[%d] %s:HP:%d\n", tick, bf.Roles[tid].Name, bf.Roles[tid].HP)

		if check_battle_end(bf) {
			bf.EQ.Clean()
			return
		}

		bf.EQ.Add(bf.Roles[id].Speed, ai_func)
	}
	return ai_func
}

func TestBattle(t *testing.T) {
	bf := &BattleField{
		Roles: [2]Role{},
		EQ:    New(),
	}
	bf.Roles[0] = Role{
		Name:  "mofon",
		HP:    1000,
		Atk:   10,
		Speed: 1000,
		AI:    AIFunc(0),
	}
	bf.Roles[1] = Role{
		Name:  "beast",
		HP:    1000,
		Atk:   11,
		Speed: 1500,
		AI:    AIFunc(1),
	}
	bf.EQ.Add(1, bf.Roles[0].AI)
	bf.EQ.Add(1, bf.Roles[1].AI)

	bf.EQ.RunUntilEmpty(bf)
}
