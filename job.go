package slowdo

import (
	"sync"
	"time"
)

type Job[ctx, item any] struct {
	maxmait time.Duration
	context ctx
	commit  func(ctx, []item)
	itemmu  sync.Mutex
	items   []item
	timer   *time.Timer
}

func NewJob[ctx, item any](
	maxwait time.Duration, context ctx, commit func(ctx, []item),
) (*Job[ctx, item], error) {
	if maxwait <= time.Millisecond {
		return nil, ErrWaitTimeTooShort
	}
	return &Job[ctx, item]{
		maxmait: maxwait,
		context: context,
		commit:  commit,
	}, nil
}

func (jb *Job[ctx, item]) Add(it item) {
	jb.itemmu.Lock()
	defer jb.itemmu.Unlock()
	if len(jb.items) == 0 {
		jb.timer = time.AfterFunc(jb.maxmait, jb.Commit)
	}
	jb.items = append(jb.items, it)
}

func (jb *Job[ctx, item]) Commit() {
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
