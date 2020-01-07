package lru

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {
	m := NewMap(10)
	m.Set("key", 1, -1)
	fmt.Println(m.Get("key"))
}

func TestLRU(t *testing.T) {
	var capacity = 100
	m := NewMap(capacity)

	for i := 0; i < 200; i++ {
		m.Set(strconv.Itoa(i), i, -1)
	}
	for i := 0; i < 200; i++ {
		val, ok := m.Get(strconv.Itoa(i))
		if i <= capacity {
			if ok {
				t.Fatal("unexpected data",i)
			}
		} else {
			if !ok || val != i {
				t.Fatal("key",i,"wrong value",val)
			}
		}
	}
}

func TestLRU_Expire(t *testing.T){
	m := NewMap(10)
	k0,v0 := "k","v"

	m.Set(k0,v0,time.Millisecond)
	time.Sleep(time.Second)

	_,ok :=m.Get(k0)
	if ok{
		t.FailNow()
	}

}