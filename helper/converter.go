package helper

func BytesToRunes(b []byte) []rune {
	return []rune(string(b))
}
func RunesToBytes(runes []rune) []byte {
	return []byte(string(runes))
}
