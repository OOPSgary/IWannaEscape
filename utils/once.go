package utils

import "sync"

type Once struct {
	key map[string]struct{}
	m   sync.Mutex
}

func (o *Once) Init() {
	o.key = make(map[string]struct{})
}
func (o *Once) Do(id string, f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.key == nil {
		o.Init()
	}
	if _, ok := o.key[id]; ok {
		return
	}
	o.key[id] = struct{}{}
	f()
	return
}
func (o *Once) Clear(id string) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.key == nil {
		o.Init()
	}
	delete(o.key, id)
}
func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()
	if o.key == nil {
		o.Init()
	}
	clear(o.key)
}
