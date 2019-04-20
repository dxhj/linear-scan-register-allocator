package main

import "errors"

type RegisterPool struct {
	head      int
	registers []Register
}

type Register int

const (
	NoRegister = iota
	EAX
	EBX
	ECX
	EDX
	ESP
	EBP
	ESI
	EDI
)

var registers = [...]string{
	NoRegister: "No Register",
	EAX:        "eax",
	EBX:        "ebx",
	ECX:        "ecx",
	EDX:        "edx",
	ESP:        "esp",
	EBP:        "ebp",
	ESI:        "esi",
	EDI:        "edi",
}

func (register Register) getRegisterName() string {
	return registers[register]
}

func (pool *RegisterPool) getRegister() (Register, error) {
	if len(pool.registers) == 0 {
		return -1, errors.New("register pool: no remaining register")
	}
	register := pool.registers[0]
	pool.registers = pool.registers[1:]
	return register, nil
}

func (pool *RegisterPool) freeRegister(register Register) {
	for _, r := range pool.registers {
		if r == register {
			return
		}
	}
	pool.registers = append(pool.registers, register)
}

func (pool *RegisterPool) isEmpty() bool {
	return len(pool.registers) == 0
}
