package shelltoken_test

import (
	"testing"

	"github.com/sni/shelltoken"
)

func BenchmarkParseShort(b *testing.B) {
	tst := `"test" some more ' test test test 123'`
	for n := 0; n < b.N; n++ {
		shelltoken.SplitLinux(tst)
	}
}

func BenchmarkParseLong(b *testing.B) {
	tst := `"test" some more ' test test test 123'`
	for x := 0; x < 10; x++ {
		tst += tst
	}

	for n := 0; n < b.N; n++ {
		shelltoken.SplitLinux(tst)
	}
}
