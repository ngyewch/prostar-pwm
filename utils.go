package prostar_pwm

func checkBit[V uint8 | uint16 | uint32 | uint64](v V, bitNo int) bool {
	return ((v >> bitNo) & 1) == 1
}

func getBits(v uint16, bitNo int, mask uint16) uint16 {
	return (v >> bitNo) & mask
}
