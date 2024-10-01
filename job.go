package slowdo

import (
	"sync"
	"time"
)

type Job[item any] struct {
	maxmait time.Duration
	commit  func([]item)
	itemmu  sync.Mutex
	items   []item
	timer   *time.Timer
}

func NewJob[item any](
	maxwait time.Duration, commit func([]item),
) (*Job[item], error) {
	if maxwait <= time.Millisecond {
		return nil, ErrWaitTimeTooShort
	}
	return &Job[item]{
		maxmait: maxwait,
		commit:  commit,
	}, nil
}

func (jb *Job[item]) Add(it item) {
	jb.itemmu.Lock()
	defer jb.itemmu.Unlock()
	if len(jb.items) == 0 {
		defer jb.collect()
	}
	jb.items = append(jb.items, it)
}

func (jb *Job[item]) Commit() {
	jb.itemmu.Lock()
	if jb.timer != nil {
		jb.timer.Stop()
		jb.timer = nil
	}
	if len(jb.items) == 0 {
		jb.itemmu.Unlock()
		return
	}
	itemscp := make([]item, len(jb.items))
	copy(itemscp, jb.items)
	jb.items = jb.items[:0]
	jb.itemmu.Unlock()
	jb.commit(itemscp)
}

func (jb *Job[item]) collect() {
	jb.itemmu.Lock()
	defer jb.itemmu.Unlock()
	jb.timer = time.AfterFunc(jb.maxmait, jb.Commit)
}
