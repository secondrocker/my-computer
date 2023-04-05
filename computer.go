package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	cm "mc/components"
	per "mc/peripherals"
)

var (
	binPath string
)

func init() {
	flag.StringVar(&binPath, "bin", "EOR", "file path")
}
func main() {
	flag.Parse()
	cmp := initComputer()
	cmp.loadInstruction()
	for {
		cmp.mar.Set(cmp.iar.Get())
		instruction := cmp.ram.GetData(cmp.mar.Get())

		cmp.alu.SetOpt(0)
		cmp.alu.SetInputA(cmp.iar.Get())
		cmp.alu.BusOne()
		cmp.acc.Set(cmp.alu.GetValue())
		cmp.iar.Set(cmp.acc.Get())

		cmp.ir.Set(instruction)

		cmp.processInstruction()
	}
}

type Computer struct {
	// alu
	alu *cm.ALU

	// 指令寄存器
	ir *cm.Register
	// 指令地址寄存器
	iar *cm.Register

	// 内存地址寄存器
	mar *cm.Register

	//
	acc *cm.Register

	// 状态寄存器
	flags *cm.Register
	// alu寄存器
	tmp *cm.Register
	// 内存
	ram *cm.RAM
	// 通用寄存器
	registers [4]*cm.Register
	// 外围设备寄存器
	adapter *per.Adapter
}

func initComputer() *Computer {
	cmp := &Computer{}
	cmp.alu = &cm.ALU{}
	cmp.ram = &cm.RAM{}
	cmp.ir = cm.NewRegister("ir")
	cmp.iar = cm.NewRegister("iar")
	cmp.mar = cm.NewRegister("mar")
	cmp.acc = cm.NewRegister("acc")
	cmp.flags = cm.NewRegister("flags")
	cmp.tmp = cm.NewRegister("tmp")
	cmp.registers = [4]*cm.Register{
		cm.NewRegister("r0"),
		cm.NewRegister("r1"),
		cm.NewRegister("r2"),
		cm.NewRegister("r3"),
	}
	cmp.adapter = per.NewAdapter()
	return cmp
}

func (c *Computer) loadInstruction() bool {
	data, err := ioutil.ReadFile(binPath)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(data); i += 2 {
		bt := data[i : i+2]
		c.ram.SetData(i/2, int(binary.BigEndian.Uint16(bt)))
	}
}

func (c *Computer) processInstruction() bool {
	num := c.ir.Get()
	leftPart := num >> 4
	rightPart := num << 4
	//算数计算
	if leftPart >= 8 {
		c.arithmeticProcess(leftPart, rightPart)
	} else {
		// 加载与存储指令 + 数据指令
		// 仅执行一个
		c.LD_ST(leftPart, rightPart)
		c.DATA(leftPart, rightPart)
		// 跳转指令
		c.JMPR(leftPart, rightPart)
		c.JMP(leftPart)
		c.JMPIF(leftPart, rightPart)
		//清除
		c.CLF(leftPart)
		// io
		c.IN_OUT(leftPart, rightPart)
	}

	return true
}

// 算数计算
func (c *Computer) arithmeticProcess(left, right int) {
	left = left << 1
	rL := c.registers[right>>2]
	rR := c.registers[right<<2]
	c.alu.SetOpt(left)
	c.alu.SetInputA(rL.Get())
	c.tmp.Set(rR.Get())
	c.alu.SetInputB(c.tmp.Get())
	// add shr shl not and or xor
	if left != 7 {
		c.acc.Set(c.alu.GetValue())
		rR.Set(c.acc.Get())
		c.flags.Set(c.alu.GetFlags())
	} else { //cmp
		c.flags.Set(c.alu.GetFlags())
	}
}

// 000 0/1
// 加载与存储指令
func (c *Computer) LD_ST(left, right int) {
	if left > 1 {
		return
	}
	rL := c.registers[right>>2]
	rR := c.registers[right<<2]
	// load
	if left == 0 {
		c.mar.Set(rL.Get())
		rR.Set(c.mar.Get())
	} else { // store
		c.mar.Set(rL.Get())
		c.ram.SetData(c.mar.Get(), rR.Get())
	}
}

// DATA 0010
// 数据指令
func (c *Computer) DATA(left, right int) {
	if left != 2 {
		return
	}
	rL := c.registers[right>>2]
	c.mar.Set(c.iar.Get())
	c.alu.SetInputA(c.iar.Get())
	c.alu.BusOne()
	c.acc.Set(c.alu.GetValue())
	rL.Set(c.ram.GetData(c.mar.Get()))
	c.iar.Set(c.acc.Get())
}

// JMPR jump to the address in rb
func (c *Computer) JMPR(left, right int) {
	if left != 3 {
		return
	}
	rR := c.registers[right<<2]

	c.iar.Set(rR.Get())
}

// JMP jump to the address in next byte
func (c *Computer) JMP(left int) {
	if left != 4 {
		return
	}
	c.mar.Set(c.iar.Get())
	c.iar.Set(c.ram.GetData(c.mar.Get()))
}

// JMPIF
func (c *Computer) JMPIF(left, right int) {
	if left != 5 {
		return
	}
	c.mar.Set(c.iar.Get())
	caez := c.flags.Get()
	if (caez & right) > 0 {
		c.iar.Set(c.ram.GetData(c.mar.Get()))
	}
}

// CLF clear all flags
func (c *Computer) CLF(left int) {
	if left != 6 {
		return
	}
	c.alu.SetInputA(0)
	c.alu.BusOne()
	c.alu.SetOpt(0)
	c.alu.GetValue()
	c.flags.Set(c.alu.GetFlags())
}

// IN/OUT
func (c *Computer) IN_OUT(left, right int) {
	if left != 7 {
		return
	}
	rg := c.registers[right<<2]
	codes := []byte(fmt.Sprintf("%04b", right))
	if codes[0] == 48 { //in
		if codes[1] == 48 { //data
			rg.Set(c.adapter.Get())
		} else { //address
			rg.Set(c.adapter.GetAddr())
		}
	} else { //out
		if codes[1] == 48 { //data
			c.adapter.Set(rg.Get())
		} else { //address
			c.adapter.SetAddr(rg.Get())
		}
	}
}
