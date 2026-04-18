package main

import (
	"fmt"
	"sort"
)

type Symbol struct {
	name     string
	register Register
	location int
}

type Allocator struct {
	pool          *RegisterPool
	active        ActiveLiveIntervals
	nextStackSlot int
	stackSlotSize int
	trace         bool
}

func NewAllocator(pool *RegisterPool) *Allocator {
	return &Allocator{
		pool:          pool,
		stackSlotSize: 4,
	}
}

func (a *Allocator) Allocate(intervals []LiveInterval) {
	sort.Sort(ByStartPoint(intervals))

	for i := range intervals {
		iv := &intervals[i]
		a.expireOldIntervals(iv)

		if a.pool.Empty() {
			a.spillAtInterval(iv)
			continue
		}

		r, err := a.pool.Acquire()
		if err != nil {
			continue
		}
		iv.symbol.register = r
		a.logf("ALLOCATE %s -> %s", iv.symbol.name, r)
		a.active = insertByEndPoint(a.active, iv)
	}
}

func (a *Allocator) expireOldIntervals(iv *LiveInterval) {
	i := 0
	for ; i < len(a.active); i++ {
		if a.active[i].endPoint >= iv.startPoint {
			break
		}
		expired := a.active[i]
		a.pool.Release(expired.symbol.register)
		a.logf("EXPIRE  %s (free %s)", expired.symbol.name, expired.symbol.register)
	}
	a.active = a.active[i:]
}

func (a *Allocator) spillAtInterval(iv *LiveInterval) {
	spill := a.active[len(a.active)-1]
	if spill.endPoint > iv.endPoint {
		a.logf("SPILL   %s (reassign %s -> %s)",
			spill.symbol.name, spill.symbol.register, iv.symbol.name)
		iv.symbol.register = spill.symbol.register
		spill.symbol.register = NoRegister
		spill.symbol.location = a.nextStackSlot
		a.active = a.active[:len(a.active)-1]
		a.active = insertByEndPoint(a.active, iv)
	} else {
		a.logf("SPILL   %s (kept on stack)", iv.symbol.name)
		iv.symbol.location = a.nextStackSlot
	}
	a.nextStackSlot += a.stackSlotSize
}

func (a *Allocator) logf(format string, args ...any) {
	if a.trace {
		fmt.Printf(format+"\n", args...)
	}
}

func main() {
	intervals := []LiveInterval{
		{symbol: Symbol{name: "a", location: -1}, startPoint: 1, endPoint: 4},
		{symbol: Symbol{name: "b", location: -1}, startPoint: 2, endPoint: 6},
		{symbol: Symbol{name: "c", location: -1}, startPoint: 3, endPoint: 10},
		{symbol: Symbol{name: "d", location: -1}, startPoint: 5, endPoint: 9},
		{symbol: Symbol{name: "e", location: -1}, startPoint: 7, endPoint: 8},
	}

	alloc := NewAllocator(NewRegisterPool(EAX, EBX))
	alloc.trace = true
	alloc.Allocate(intervals)

	fmt.Println("--- assignments ---")
	for _, iv := range intervals {
		if iv.symbol.register != NoRegister {
			fmt.Printf("%s -> %s\n", iv.symbol.name, iv.symbol.register)
		} else {
			fmt.Printf("%s -> stack[%d]\n", iv.symbol.name, iv.symbol.location)
		}
	}
}
