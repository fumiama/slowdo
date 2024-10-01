package slowdo

import (
	"sync"
	"time"
)

type Job[item, ctx any] struct {
	maxmait time.Duration
	context ctx
	commit  func(ctx, []item)
	itemmu  sync.Mutex
	items   []item
	timer   *time.Timer
}

func NewJob[item, ctx any](
	maxwait time.Duration, context ctx, commit func(ctx, []item),
) (*Job[item, ctx], error) {
	if maxwait <= time.Millisecond {
		return nil, ErrWaitTimeTooShort
	}
	return &Job[item, ctx]{
		maxmait: maxwait,
		context: context,
		commit:  commit,
	}, nil
}

func (jb *Job[item, ctx]) Add(it item) {
	jb.itemmu.Lock()
	defer jb.itemmu.Unlock()
	if len(jb.items) == 0 {
		defer jb.collect()
	}
	jb.items = append(jb.items, it)
}

func (jb *Job[item, ctx]) Commit() {
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
	jb.commit(jb.context, itemscp)
}

func (jb *Job[item, ctx]) collect() {
	jb.itemmu.Lock()
	defer jb.itemmu.Unlock()
	jb.timer = time.AfterFunc(jb.maxmait, jb.Commit)
}
