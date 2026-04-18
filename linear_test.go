package main

import "testing"

func byName(intervals []LiveInterval) map[string]LiveInterval {
	m := make(map[string]LiveInterval, len(intervals))
	for _, iv := range intervals {
		m[iv.symbol.name] = iv
	}
	return m
}

func TestPaperExample(t *testing.T) {
	intervals := []LiveInterval{
		{symbol: Symbol{name: "a", location: -1}, startPoint: 1, endPoint: 4},
		{symbol: Symbol{name: "b", location: -1}, startPoint: 2, endPoint: 6},
		{symbol: Symbol{name: "c", location: -1}, startPoint: 3, endPoint: 10},
		{symbol: Symbol{name: "d", location: -1}, startPoint: 5, endPoint: 9},
		{symbol: Symbol{name: "e", location: -1}, startPoint: 7, endPoint: 8},
	}

	alloc := NewAllocator(NewRegisterPool(EAX, EBX))
	alloc.Allocate(intervals)

	got := byName(intervals)
	wantReg := map[string]Register{
		"a": EAX,
		"b": EBX,
		"d": EAX,
		"e": EBX,
	}
	for name, want := range wantReg {
		if have := got[name].symbol.register; have != want {
			t.Errorf("%s: register = %s, want %s", name, have, want)
		}
	}
	if r := got["c"].symbol.register; r != NoRegister {
		t.Errorf("c: register = %s, want NoRegister (spilled)", r)
	}
	if loc := got["c"].symbol.location; loc != 0 {
		t.Errorf("c: stack location = %d, want 0", loc)
	}
}

func TestExpireFreesAllEligible(t *testing.T) {
	intervals := []LiveInterval{
		{symbol: Symbol{name: "a", location: -1}, startPoint: 1, endPoint: 2},
		{symbol: Symbol{name: "b", location: -1}, startPoint: 3, endPoint: 4},
		{symbol: Symbol{name: "c", location: -1}, startPoint: 5, endPoint: 6},
	}
	alloc := NewAllocator(NewRegisterPool(EAX))
	alloc.Allocate(intervals)

	for _, iv := range intervals {
		if iv.symbol.register != EAX {
			t.Errorf("%s: register = %s, want EAX", iv.symbol.name, iv.symbol.register)
		}
	}
}

func TestSpillActiveInFavorOfNewer(t *testing.T) {
	intervals := []LiveInterval{
		{symbol: Symbol{name: "long", location: -1}, startPoint: 1, endPoint: 100},
		{symbol: Symbol{name: "short", location: -1}, startPoint: 2, endPoint: 5},
	}
	alloc := NewAllocator(NewRegisterPool(EAX))
	alloc.Allocate(intervals)

	got := byName(intervals)
	if r := got["short"].symbol.register; r != EAX {
		t.Errorf("short: register = %s, want EAX", r)
	}
	if r := got["long"].symbol.register; r != NoRegister {
		t.Errorf("long: register = %s, want NoRegister (spilled)", r)
	}
	if loc := got["long"].symbol.location; loc != 0 {
		t.Errorf("long: stack location = %d, want 0", loc)
	}
}

func TestSpillIncomingInterval(t *testing.T) {
	intervals := []LiveInterval{
		{symbol: Symbol{name: "short", location: -1}, startPoint: 1, endPoint: 5},
		{symbol: Symbol{name: "long", location: -1}, startPoint: 2, endPoint: 100},
	}
	alloc := NewAllocator(NewRegisterPool(EAX))
	alloc.Allocate(intervals)

	got := byName(intervals)
	if r := got["short"].symbol.register; r != EAX {
		t.Errorf("short: register = %s, want EAX", r)
	}
	if r := got["long"].symbol.register; r != NoRegister {
		t.Errorf("long: register = %s, want NoRegister (spilled)", r)
	}
	if loc := got["long"].symbol.location; loc != 0 {
		t.Errorf("long: stack location = %d, want 0", loc)
	}
}

func TestStackSlotsAreDistinct(t *testing.T) {
	intervals := []LiveInterval{
		{symbol: Symbol{name: "a", location: -1}, startPoint: 1, endPoint: 100},
		{symbol: Symbol{name: "b", location: -1}, startPoint: 2, endPoint: 100},
		{symbol: Symbol{name: "c", location: -1}, startPoint: 3, endPoint: 100},
		{symbol: Symbol{name: "d", location: -1}, startPoint: 4, endPoint: 100},
	}
	alloc := NewAllocator(NewRegisterPool(EAX, EBX))
	alloc.Allocate(intervals)

	seen := map[int]string{}
	for _, iv := range intervals {
		if iv.symbol.register != NoRegister {
			continue
		}
		if other, dup := seen[iv.symbol.location]; dup {
			t.Errorf("%s and %s share stack slot %d",
				other, iv.symbol.name, iv.symbol.location)
		}
		seen[iv.symbol.location] = iv.symbol.name
	}
	if len(seen) != 2 {
		t.Errorf("expected 2 spills, got %d", len(seen))
	}
}

func TestInsertByEndPointKeepsOrder(t *testing.T) {
	mk := func(end int) *LiveInterval {
		return &LiveInterval{endPoint: end}
	}
	var active ActiveLiveIntervals
	for _, end := range []int{5, 1, 3, 4, 2} {
		active = insertByEndPoint(active, mk(end))
	}
	prev := -1
	for _, iv := range active {
		if iv.endPoint < prev {
			t.Fatalf("active not sorted: endpoints out of order around %d", iv.endPoint)
		}
		prev = iv.endPoint
	}
	if len(active) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(active))
	}
}

func TestRegisterPoolEmptyError(t *testing.T) {
	p := NewRegisterPool(EAX)
	if _, err := p.Acquire(); err != nil {
		t.Fatalf("first Acquire should succeed, got %v", err)
	}
	r, err := p.Acquire()
	if err == nil {
		t.Fatalf("second Acquire should fail on empty pool")
	}
	if r != NoRegister {
		t.Errorf("failed Acquire returned %s, want NoRegister", r)
	}
	if !p.Empty() {
		t.Errorf("pool should report Empty() == true")
	}
}
