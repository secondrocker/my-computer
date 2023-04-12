%display_addr = 0x07

data r0, 0xf000
data r1, 0x7e40
store r0, r1

echo:
data r0, 0xf001
data r1,0x407c
store r0, r1

data r0, 0xf002
data r1, 0x4040
store r0, r1

data r0, 0xf003
data r1, 0x7e00
store r0, r1

data r0, %display_addr
out addr, r0

display:
  data r0, 0xf000
  load r0, r1
  out data, r1

  data r0, 0xf001
  load r0, r1
  out data, r1

  data r0, 0xf002
  load r0, r1
  out data, r1

  data r0, 0xf003
  load r0, r1
  out data, r1
  jmp display

