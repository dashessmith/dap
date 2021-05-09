package utils

func RuneIn(r rune, them ...rune) bool {
	for _, x := range them {
		if x == r {
			return true
		}
	}
	return false
}
