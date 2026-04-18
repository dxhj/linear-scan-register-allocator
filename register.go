package main

import "errors"

type Register int

const (
	NoRegister Register = iota
	EAX
	EBX
	ECX
	EDX
	ESP
	EBP
	ESI
	EDI
)

var registerNames = [...]string{
	NoRegister: "no-register",
	EAX:        "eax",
	EBX:        "ebx",
	ECX:        "ecx",
	EDX:        "edx",
	ESP:        "esp",
	EBP:        "ebp",
	ESI:        "esi",
	EDI:        "edi",
}

func (r Register) String() string {
	if int(r) < 0 || int(r) >= len(registerNames) {
		return "unknown"
	}
	return registerNames[r]
}

var ErrNoRegister = errors.New("register pool: no free register")

type RegisterPool struct {
	free []Register
}

func NewRegisterPool(registers ...Register) *RegisterPool {
	return &RegisterPool{free: append([]Register(nil), registers...)}
}

func (p *RegisterPool) Acquire() (Register, error) {
	if len(p.free) == 0 {
		return NoRegister, ErrNoRegister
	}
	r := p.free[0]
	p.free = p.free[1:]
	return r, nil
}

func (p *RegisterPool) Release(r Register) {
	p.free = append(p.free, r)
}

func (p *RegisterPool) Empty() bool {
	return len(p.free) == 0
}
