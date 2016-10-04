package strutil

import "math/rand"

var a = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandASCII generates an alphabetical random string with length n.
func RandASCII(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = a[rand.Intn(len(a))]
	}
	return string(b)
}

// Join is alternative for http://golang.org/pkg/strings/#Join.
func Join(a ...string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	}

	var n int
	for _, v := range a {
		n += len(v)
	}

	b := make([]byte, n)

	n = 0
	for _, v := range a {
		n += copy(b[n:], v)
	}

	return string(b)
}

// First returns n first characters.
// Based on https://groups.google.com/forum/#!topic/golang-nuts/oPuBaYJ17t4.
func TrimRightN(s string, n int) string {
	if s == "" || n == 0 {
		return s
	}

	if n < 0 {
		n = -n
	}

	r := []rune(s)
	if len(r) > n {
		return string(r[0:n])
	}

	return string(r)
}
