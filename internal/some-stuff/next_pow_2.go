package some_stuff

func NextPowOf2(n int64) int64 {
	var power int64
	power = 1

	for power < n {
		power <<= 1
	}

	return power
}
