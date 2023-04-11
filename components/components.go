package components

type Register struct {
	Name  string
	value int
}

func NewRegister(name string) *Register {
	return &Register{Name: name, value: 0}
}

func (r *Register) Set(value int) {
	r.value = value
}

func (r *Register) Get() int {
	return r.value
}

type ALU struct {
	inputA int
	inputB int
	output int

	opt int

	carryIn bool

	carryOut bool
	aLarger  bool
	equal    bool
}

func (a *ALU) SetInputA(value int) {
	a.inputA = value
}

func (a *ALU) SetInputB(value int) {
	a.inputB = value
}

func (a *ALU) SetCarryIn(value bool) {
	a.carryIn = value
}

func (a *ALU) BusOne() {
	a.inputB = 1
}

func (a *ALU) execute() {
	switch a.opt {
	// ADD
	case 0:
		a.Add()
	// SHR
	case 1:
		a.Shr()
	// SHL
	case 2:
		a.Shl()
	// NOT
	case 3:
		a.Not()
	// AND
	case 4:
		a.And()
	// OR
	case 5:
		a.Or()
	// XOR
	case 6:
		a.Xor()
	// CMP
	case 7:
		a.CMP()
	}
}

func (a *ALU) SetOpt(opt int) {
	a.opt = opt
}

func (a *ALU) GetValue() int {
	a.execute()
	return a.output
}

func (a *ALU) GetFlags() int {
	a.execute()
	v := 0
	if a.carryOut {
		v += 8
	}
	if a.aLarger {
		v += 4
	}
	if a.equal {
		v += 2
	}
	if a.Zero() {
		v += 1
	}
	return v
}

func (a *ALU) Add() {
	v := a.inputA + a.inputB
	if a.carryIn {
		v += 1
	}
	if v <= 0xffff {
		a.carryOut = false
	}
	a.output = v % 0xffff
	a.carryOut = true
}

func (a *ALU) Shr() {
	if a.inputA%2 == 1 {
		a.carryOut = true
	}
	a.output = a.output >> 1
	if a.carryIn {
		a.output = a.output + 0x8000
	}
}

func (a *ALU) Shl() {
	if a.inputA%2 >= 0x8000 {
		a.carryOut = true
	}
	a.output = a.output << 1
	if a.carryIn {
		a.output += 1
	}
}

func (a *ALU) Not() {
	a.carryOut = false
	a.output = a.inputA ^ 0xffff
}

func (a *ALU) And() {
	a.carryOut = false
	a.output = a.inputA & a.inputB
}

func (a *ALU) Or() {
	a.carryOut = false
	a.output = a.inputA | a.inputB
}

func (a *ALU) Xor() {
	a.carryOut = false
	a.output = a.inputA ^ a.inputB
}

func (a *ALU) CMP() {
	a.output = 0
	a.carryOut = false
	a.aLarger = a.inputA > a.inputB
	a.equal = a.inputA == a.inputB
}

func (a *ALU) Zero() bool {
	a.execute()
	return a.output == 0
}

type RAM struct {
	data [65536]int
}

func (r *RAM) GetData(address int) int {
	return r.data[address]
}

func (r *RAM) SetData(address, value int) {
	r.data[address] = value
}

// func (r *RAM) InitInstruction(instructions []int) {
// 	r.data = instructions
// }
