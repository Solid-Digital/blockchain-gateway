package adapter

import (
	crand "crypto/rand"
	"math/big"
	"testing"
)

func BenchmarkBTNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMessage([]byte(`{
                    "date": "19-07-2017 09:45",
                    "sourceEan": "123456789325869503",
                    "sourceType": "Social",
                    "destinationEan": "987654321625183729",
                    "amountReceived": "1.42343587970964396819"
                }`))
	}
}

func BenchmarkNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMessage([]byte("{\n" +
			"\"date\": \"19-07-2017 09:45\",\n" +
			"\"sourceEan\": \"123456789325869503\",\n" +
			"\"sourceType\": \"Social\",\n" +
			"\"destinationEan\": \"987654321625183729\",\n" +
			"\"amountReceived\": \"1.42343587970964396819\"\n" +
			"}\n"))
	}
}

func BenchmarkNewTaggedMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTaggedMessage([]byte(`{
                    "date": "19-07-2017 09:45",
                    "sourceEan": "123456789325869503",
                    "sourceType": "Social",
                    "destinationEan": "987654321625183729",
                    "amountReceived": "1.42343587970964396819"
                }`))
	}
}

func BenchmarkNewTaggedMessageWithTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTaggedMessage([]byte(`{
                    "date": "19-07-2017 09:45",
                    "sourceEan": "123456789325869503",
                    "sourceType": "Social",
                    "destinationEan": "987654321625183729",
                    "amountReceived": "1.42343587970964396819"
                }`), WithTag(3))
	}
}

func BenchmarkNewTaggedMessageWithRandomTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTaggedMessage([]byte(`{
                    "date": "19-07-2017 09:45",
                    "sourceEan": "123456789325869503",
                    "sourceType": "Social",
                    "destinationEan": "987654321625183729",
                    "amountReceived": "1.42343587970964396819"
                }`), WithRandomTag())
	}
}

func BenchmarkMathRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomTag()
		//rand.Uint64()
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := crand.Int(crand.Reader, big.NewInt(27))
		if err != nil {
			panic(err)
		}
	}
}
