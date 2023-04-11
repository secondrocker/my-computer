package peripherals

type Adapter struct {
	addr       int
	equipments map[int]Equipment
}

func (a *Adapter) SetAddr(addr int) {
	a.addr = addr
}

func (a *Adapter) GetAddr() int {
	return a.addr
}

func NewAdapter() *Adapter {
	display := NewDisplay()
	return &Adapter{
		equipments: map[int]Equipment{
			7: display,
		},
	}
}

func (a *Adapter) Get() int {
	return a.equipments[a.addr].Get()
}

func (a *Adapter) Set(v int) {
	a.equipments[a.addr].Set(v)
}

type Equipment interface {
	Get() int
	Set(int)
}
