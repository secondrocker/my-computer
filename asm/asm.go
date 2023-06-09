package asm

const (
	//逻辑
	ADD = 0b1000
	SHR = 0b1001
	SHL = 0b1010
	NOT = 0b1011
	AND = 0b1100
	OR  = 0b1101
	XOR = 0b1110
	CMP = 0b1111

	//数据
	LD   = 0b0000
	ST   = 0b0001
	DATA = 0b0010

	// 跳转
	JMPR = 0b0011
	JMP  = 0b0100

	JC    = 0b01011000
	JA    = 0b01010100
	JE    = 0b01010010
	JZ    = 0b01010001
	JCA   = 0b01011100
	JCE   = 0b01011010
	JCZ   = 0b01011001
	JAE   = 0b01010110
	JAZ   = 0b01010101
	JEZ   = 0b01010011
	JCAE  = 0b01011110
	JCAZ  = 0b01011101
	JCEZ  = 0b01011011
	JAEZ  = 0b01010111
	JCAEZ = 0b01011111

	// 清除
	CLF = 0b01100000

	// IN
	IN  = 0b111
	OUT = 0b111

	//寄存器
	REG0 = 0b00
	REG1 = 0b01
	REG2 = 0b10
	REG3 = 0b11
)
