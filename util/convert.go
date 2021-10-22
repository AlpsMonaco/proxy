package util

func ByteToString(b byte) string {
	var result []byte = make([]byte, 3)
	var i int
	for i = 2; i >= 0; i-- {
		result[i] = (b % 10) + 48
		b = b / 10
		if b == 0 {
			break
		}
	}

	return string(result[i:])
}
