package md5crypt

import "testing"

func BenchmarkMD5Crypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var pw = "zzzzz"
		var salt = "hfT7jp2q"
		var magic = "$1$"
		MD5Crypt([]byte(pw), []byte(salt), []byte(magic))
	}
}
