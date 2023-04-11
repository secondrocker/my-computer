package peripherals

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Display struct {
	inputMar  int
	outputMar int
	ram       [8]int
	count     int
}

func NewDisplay() *Display {
	d := &Display{
		ram:   [8]int{},
		count: 0,
	}
	d.Draw()
	return d
}

func (d *Display) SetInputMar(v int) {
	d.inputMar = v
}

func (d *Display) Set(v int) {
	d.ram[d.count] = v >> 8
	d.ram[d.count+1] = v / 0x0100
	d.count += 2
	if d.count == 8 {
		d.count = 0
	}
}

func (d *Display) Get() int {
	return 0
}

func (d *Display) Draw() {
	go func() {
		t := time.Tick(time.Millisecond * 10)
		for {
			<-t
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			for i := 0; i < 8; i++ {
				show := d.ram[i]
				fmt.Printf("%s\n", changeStr(show))
				fmt.Printf("\n")
			}
		}
	}()
}

func changeStr(v int) string {
	s := fmt.Sprintf("%08b", v)

	aa := strings.Split(s, "")

	bb := make([]byte, 0, 8)
	for _, b := range aa {
		if b == "0" {
			bb = append(bb, 32)
		} else {
			bb = append(bb, 46)
		}
	}
	return string(bb)
}
