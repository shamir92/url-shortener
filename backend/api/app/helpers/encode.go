package helpers

func ShortUrlEncoder(counterDecimal uint64) string {
	var encodeChar string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~"
	var hashString string = ""
	for counterDecimal > 0 {
		hashString = string(encodeChar[counterDecimal%66]) + hashString
		counterDecimal = uint64(counterDecimal / 66)
	}
	var hashStringLength = len([]rune(hashString))

	if hashStringLength < 7 {
		for i := 0; i < 7-hashStringLength; i++ {
			hashString = "0" + hashString
		}
	}

	return hashString
}
