/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

//Source: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang/31832326#31832326

package internal

import (
	"math/rand"
	"sync"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var mutex sync.Mutex

func int63() int64 {
	mutex.Lock()
	v := src.Int63()
	mutex.Unlock()
	return v
}

// RandomString generates a random string with the specified length
func RandomString(n int) string {
	b := make([]byte, n)
	// A int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
