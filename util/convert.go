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

// convert only byte.
// 5 = 5 * 10 ^ 0
// 61 = 6 * 10 ^ 1 + 1 * 10 ^ 0
// 121 = 1 * 10 ^ 2 + 2 * 10 ^ 1 + 1 * 10 ^ 0
func StringToByte(s string, b *byte) {
	size := byte(len(s))
	switch size {
	case 1:
		*b = s[0] - 48
	case 2:
		*b = (s[0]-48)*10 + s[1] - 48
	case 3:
		*b = ((s[0] - 48) * 100) + ((s[1] - 48) * 10) + s[2] - 48
	default:
		return
	}
}

func IPV4AddrToByte(addr string, bPtr *[]byte) {
	var i, j, k int
	addr = addr + "."
	for i = 0; i < len(addr); i++ {
		if addr[i] == 46 {
			StringToByte(addr[k:i], &(*bPtr)[j])
			j++
			k = i + 1
		}
	}
}
