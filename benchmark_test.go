package shelltoken_test

import (
	"testing"

	"github.com/sni/shelltoken"
)

func BenchmarkParseShort(b *testing.B) {
	tst := `"test" some more ' test test test 123'`
	for range b.N {
		shelltoken.SplitLinux(tst)
	}
}

func BenchmarkParseLong(b *testing.B) {
	tst := `"test" some more ' test test test 123'`
	for range 10 {
		tst += tst
	}

	for range b.N {
		shelltoken.SplitLinux(tst)
	}
}
