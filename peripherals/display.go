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
	ram       [256]int
}

func NewDisplay() *Display {
	d := &Display{
		ram: [256]int{},
	}
	d.Draw()
	return d
}

func (d *Display) SetInputMar(v int) {
	d.inputMar = v
}

func (d *Display) Set(v int) {
	d.ram[d.inputMar] = v
}

func (d *Display) Get() int {
	return 0
}

func (d *Display) Draw() {
	t := time.Tick(time.Millisecond * 10)

	for {
		<-t
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		for i := 0; i < 32; i++ {
			for j := 0; j < 8; j++ {
				show := d.ram[i*8+j]
				fmt.Printf("%s\n", changeStr(show))
			}
			fmt.Printf("\n")
		}
	}
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
