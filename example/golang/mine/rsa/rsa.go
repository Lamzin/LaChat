package rsa

import (
	"math/rand"
	"time"
)

type RSA struct {
	numb      int64
	encodeExp int64
	decodeExp int64
}

func NewRSA() (r RSA) {
	rand.Seed(time.Now().UTC().UnixNano())

	var p, q int64 = getTwoSimple()
	var fi int64 = (p - 1) * (q - 1)

	r.numb = p * q
	for {
		if r.encodeExp = rand.Int63n(fi) + 1; gcd(fi, r.encodeExp) == 1 {
			break
		}
	}

	for _, r.decodeExp, _ = gcdExtended(r.encodeExp, fi); r.decodeExp < 0; r.decodeExp += fi {
	}

	return r
}

func binpow(numb, pow, mod int64) int64 {
	res := int64(1)
	for pow > 0 {
		if pow%2 == 1 {
			res = (res * numb) % mod
		}
		numb = (numb * numb) % mod
		pow /= 2
	}
	return res
}

func (r RSA) encode(x int64) int64 {
	return binpow(x, r.encodeExp, r.numb)
}

func (r RSA) decode(x int64) int64 {
	return binpow(x, r.decodeExp, r.numb)
}

func isSimple(n int64) bool {
	if n%2 == 0 || n < 2 {
		return false
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func getTwoSimple() (a, b int64) {
	a = 10000 + rand.Int63n(1000)
	for ; ; a++ {
		if isSimple(a) {
			break
		}
	}
	b = a + 1
	for ; ; b++ {
		if isSimple(b) {
			break
		}
	}
	return
}

func gcdExtended(a, b int64) (gcd, x, y int64) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x, y = gcdExtended(b%a, a)
	return gcd, y - (b/a)*x, x
}

func gcd(a, b int64) int64 {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}
