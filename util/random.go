package util

import (
	"math/rand" 
	"strings"
	"time"
	// "fmt"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomInt(min , max int32) int32  {
	return min + rand.Int31n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k:= len(alphabet)
	for i := 0; i < n; i++ {
        c := alphabet[rand.Intn(k)]
        sb.WriteByte(c)
    }

	return sb.String()
}

func RandomCustomerId() int32 {
	return RandomInt(0,20);
}

func RandomProductName() string {
	return RandomString(6);
}

func RandomProductPrice() int32 {
    return RandomInt(0,1000);
}

func RandomProductQuantity() int32 {
    return RandomInt(0,1000);
}

func RandomMoney() int32 {
    return RandomInt(0,1000);
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

