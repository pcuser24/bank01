package util

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

var ALPHANUMERIC string

var Currencies = [...]string{
	"USD", "EUR", "JPY",
	"VND", "CNY", "RUB",
}

func init() {
	rand.Seed(time.Now().UnixNano())

	const ALPHANUM = ALPHABET + "0123456789"
	chars := []byte(ALPHANUM)
	for i := len(ALPHANUM) - 1; i > 0; i-- {
		j := RandomInt(0, int64(i))
		temp := chars[i]
		chars[i] = chars[j]
		chars[j] = temp
	}
	ALPHANUMERIC = string(chars)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomStr(n int, alphabet string) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomAlphabetStr(n int) string {
	return randomStr(n, ALPHABET)
}

func RandomAlphanumericStr(n int) string {
	return randomStr(n, ALPHANUMERIC)
}

func RandomOwner() string {
	return RandomAlphabetStr(6)
}

func RandomMoney() int64 {
	return rand.Int63n(1000)
}

func RandomCurrency() string {
	return Currencies[rand.Intn(len(Currencies))]
}

func RandomCountryCode() sql.NullInt32 {
	return sql.NullInt32{Int32: rand.Int31n(999), Valid: true}
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomAlphabetStr(8))
}
