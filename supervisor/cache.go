package supervisor

import (
	"sync"
)

type cache struct {
	sync.Map
}

var c cache

func (c *cache) get(name string) *Program {
	if p, ok := c.Load(name); ok {
		return p.(*Program)
	}
	return nil
}

func (c *cache) set(name, last string, state State) {
	if p := c.get(name); p != nil {
		p.State = state
		p.Last = last
	} else {
		p = new(Program)
		p.Name = name
		p.State = state
		p.Last = last

		c.Store(name, p)
	}
}

func (c *cache) del(name string) {
	c.Delete(name)
}

func (c *cache) list() (temp []Program) {
	c.Range(func(k, v interface{}) bool {
		t := v.(*Program)
		temp = append(temp, Program{Name: t.Name, State: t.State, Last: t.Last})
		return true
	})
	return
}
