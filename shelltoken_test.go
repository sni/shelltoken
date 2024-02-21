package shelltoken_test

import (
	"testing"

	"github.com/sni/shelltoken"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellToken(t *testing.T) {
	tests := []struct {
		in  string
		res []string
	}{
		{"", []string{""}},
		{" a", []string{"a"}},
		{" a ", []string{"a"}},
		{"a bc d", []string{"a", "bc", "d"}},
		{"a 'bc' d", []string{"a", "bc", "d"}},
		{"a 'b c' d", []string{"a", "b c", "d"}},
		{`a "b'c" d`, []string{"a", `b'c`, "d"}},
		{`a 'b""c' d`, []string{"a", `b""c`, "d"}},
		{`a  """b""" '' ''c'' ''d'' ee""ee f' 'f '" "' "' ''"`, []string{`a`, `b`, ``, `c`, `d`, `eeee`, `f f`, `" "`, `' ''`}},
		{`"\'"`, []string{`\'`}},
		{`"\"'"`, []string{`"'`}},
		{`\'`, []string{`'`}},
		{`\"`, []string{`"`}},
		{`'"\a"'`, []string{`"\a"`}},
		{`\ a`, []string{` a`}},
		{`\\ a`, []string{`\`, `a`}},
		{`\\\ a`, []string{`\ a`}},
		{`\\\\ a`, []string{`\\`, `a`}},
		{`"\\\\ a"`, []string{`\\ a`}},
		{`'\\\\ a'`, []string{`\\\\ a`}},
		{"/bin/sh -c 'echo a b c '", []string{"/bin/sh", "-c", "echo a b c "}},
		{"cmd.exe /c '`foo`'", []string{"cmd.exe", "/c", "`foo`"}},
		// {`c:\"Program Files"\/crap\bs.exe`, []string{`c:\Program Files\/crap\bs.exe`}},
	}

	for _, tst := range tests {
		env, argv, _, err := shelltoken.Parse(tst.in, false)
		require.NoErrorf(t, err, "error while parsing: %s", tst.in)
		assert.Equalf(t, tst.res, argv, "Tokenize: %v -> %v", tst.in, argv)
		assert.Emptyf(t, env, "no env")
	}
}

func TestShellEnv(t *testing.T) {
	tests := []struct {
		in  string
		env []string
		arg []string
	}{
		{"test", []string{}, []string{"test"}},
		{"test arg1 arg2", []string{}, []string{"test", "arg1", "arg2"}},
		{"./test", []string{}, []string{"./test"}},
		{"./blah/test", []string{}, []string{"./blah/test"}},
		{"/blah/test arg1 arg2", []string{}, []string{"/blah/test", "arg1", "arg2"}},
		{"./test arg1 arg2", []string{}, []string{"./test", "arg1", "arg2"}},
		{"ENV1=1 ENV2=2 /blah/test", []string{"ENV1=1", "ENV2=2"}, []string{"/blah/test"}},
		{"ENV1=1 ENV2=2 ./test", []string{"ENV1=1", "ENV2=2"}, []string{"./test"}},
		{"ENV1=1 ENV2=2 ./test arg1 arg2", []string{"ENV1=1", "ENV2=2"}, []string{"./test", "arg1", "arg2"}},
		{`ENV1="1 2 3" ENV2='2' ./test arg1 arg2`, []string{"ENV1=1 2 3", "ENV2=2"}, []string{"./test", "arg1", "arg2"}},
		{`PATH=test:$PATH LD_LIB=... $(pwd)/test`, []string{"PATH=test:$PATH", "LD_LIB=..."}, []string{"$(pwd)/test"}},
		{"/python /tmp/file1 args1", []string{}, []string{"/python", "/tmp/file1", "args1"}},
		{"lib/negate /bin/python3 /tmp/file1 args1", []string{}, []string{"lib/negate", "/bin/python3", "/tmp/file1", "args1"}},
		{`ENV1="1 2 3" ENV2='2' ./test arg1 -P 'm1|m2';`, []string{"ENV1=1 2 3", "ENV2=2"}, []string{"./test", "arg1", "-P", "m1|m2;"}},
	}

	for _, tst := range tests {
		env, argv, _, err := shelltoken.Parse(tst.in, false)
		require.NoErrorf(t, err, "error while parsing: %s", tst.in)
		assert.Equalf(t, tst.arg, argv, "Tokenize: %v -> %v", tst.in, argv)
		assert.Equalf(t, tst.env, env, "Tokenize env: %v -> %v", tst.in, env)
	}
}

func TestShellErrors(t *testing.T) {
	tests := []struct {
		in  string
		err string
	}{
		{"test 'arg1 arg2", "unbalanced quotes"},
	}

	for _, tst := range tests {
		env, argv, _, err := shelltoken.Parse(tst.in, false)
		require.Errorf(t, err, "expected error for %s: %s", tst.in, tst.err)
		assert.Contains(t, err.Error(), tst.err)
		assert.Nil(t, argv, "argv is nil")
		assert.Nil(t, env, "argv is nil")
	}
}
