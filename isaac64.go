// Copyright 2021 skdltmxn. All rights reserved.
//
// isaac64.go
//
// ISAAC64 random generator by Bob Jenkins

package isaac

import "unsafe"

// Isaac64 represents ISAAC64 random generator
type Isaac64 struct {
	randrsl [256]uint64
	randmem [256]uint64
	randcnt uint64
	aa      uint64
	bb      uint64
	cc      uint64
}

// NewIsaac64 returns a new instance of ISAAC64.
func NewIsaac64() *Isaac64 {
	return &Isaac64{
		randmem: [256]uint64{},
		randrsl: [256]uint64{},
		randcnt: 0,
		aa:      0,
		bb:      0,
		cc:      0,
	}
}

// Seed initializes the state of ISAAC instance using given 64bit integer.
func (ctx *Isaac64) Seed(seed int64) {
	ctx.randrsl[0] = uint64(seed)
	ctx.randInit(true)
}

// SeedBytes initializes the state of ISAAC instance using given byte sequence.
func (ctx *Isaac64) SeedBytes(seed []byte) {
	if len(seed) > 2048 {
		seed = seed[:2048]
	}

	unsafeCopy(unsafe.Pointer(&ctx.randrsl[0]), unsafe.Pointer(&seed[0]), len(seed))
	ctx.randInit(true)
}

// SeedString initializes the state of ISAAC instance using given string.
func (ctx *Isaac64) SeedString(seed string) {
	ctx.SeedBytes([]byte(seed))
}

// Int63 returns a non-negative 63-bit integer as an int64.
func (ctx *Isaac64) Int63() int64 {
	return int64(ctx.Uint64() & uintMask)
}

// Uint32 returns a random 32-bit unsigned integer.
func (ctx *Isaac64) Uint32() uint32 {
	return uint32(ctx.next())
}

// Uint64 returns a random 64-bit unsigned integer.
func (ctx *Isaac64) Uint64() uint64 {
	return ctx.next()
}

// Int31 returns a non-negative 31-bit integer as an int32.
func (ctx *Isaac64) Int31() int32 {
	return int32(ctx.Int63() >> 32)
}

// Int returns a non-negative integer as an int
func (ctx *Isaac64) Int() int {
	u := uint(ctx.Uint64())
	return int(u << 1 >> 1)
}

func (ctx *Isaac64) isaac64() {
	var a, b, x uint64
	mm := ctx.randmem[:]
	r := ctx.randrsl[:]
	ctx.cc++
	a, b = ctx.aa, ctx.bb+ctx.cc

	for ii := 0; ii < 256; ii += 4 {
		var i uint8 = uint8(ii)

		x = mm[i]
		a = ^(a ^ (a << 21)) + mm[i+128]
		mm[i] = mm[(x>>3)&255] + a + b
		r[i] = mm[(mm[i]>>11)&255] + x
		b = r[i]

		x = mm[i+1]
		a = (a ^ (a >> 5)) + mm[i+129]
		mm[i+1] = mm[(x>>3)&255] + a + b
		r[i+1] = mm[(mm[i+1]>>11)&255] + x
		b = r[i+1]

		x = mm[i+2]
		a = (a ^ (a << 12)) + mm[i+130]
		mm[i+2] = mm[(x>>3)&255] + a + b
		r[i+2] = mm[(mm[i+2]>>11)&255] + x
		b = r[i+2]

		x = mm[i+3]
		a = (a ^ (a >> 33)) + mm[i+131]
		mm[i+3] = mm[(x>>3)&255] + a + b
		r[i+3] = mm[(mm[i+3]>>11)&255] + x
		b = r[i+3]
	}

	ctx.bb, ctx.aa = b, a
}

func (ctx *Isaac64) randInit(flag bool) {
	var a, b, c, d, e, f, g, h uint64
	a, b, c, d, e, f, g, h = 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13, 0x9e3779b97f4a7c13

	// scramble
	for i := 0; i < 4; i++ {
		a -= e
		f ^= h >> 9
		h += a
		b -= f
		g ^= a << 9
		a += b
		c -= g
		h ^= b >> 23
		b += c
		d -= h
		a ^= c << 15
		c += d
		e -= a
		b ^= d >> 14
		d += e
		f -= b
		c ^= e << 20
		e += f
		g -= c
		d ^= f >> 17
		f += g
		h -= d
		e ^= g << 14
		g += h
	}

	if flag {
		// initialize using seed
		for i := 0; i < 256; i += 8 {
			a += ctx.randrsl[i]
			b += ctx.randrsl[i+1]
			c += ctx.randrsl[i+2]
			d += ctx.randrsl[i+3]
			e += ctx.randrsl[i+4]
			f += ctx.randrsl[i+5]
			g += ctx.randrsl[i+6]
			h += ctx.randrsl[i+7]

			// mix
			a -= e
			f ^= h >> 9
			h += a
			b -= f
			g ^= a << 9
			a += b
			c -= g
			h ^= b >> 23
			b += c
			d -= h
			a ^= c << 15
			c += d
			e -= a
			b ^= d >> 14
			d += e
			f -= b
			c ^= e << 20
			e += f
			g -= c
			d ^= f >> 17
			f += g
			h -= d
			e ^= g << 14
			g += h

			ctx.randmem[i] = a
			ctx.randmem[i+1] = b
			ctx.randmem[i+2] = c
			ctx.randmem[i+3] = d
			ctx.randmem[i+4] = e
			ctx.randmem[i+5] = f
			ctx.randmem[i+6] = g
			ctx.randmem[i+7] = h
		}

		// second pass
		for i := 0; i < 256; i += 8 {
			a += ctx.randmem[i]
			b += ctx.randmem[i+1]
			c += ctx.randmem[i+2]
			d += ctx.randmem[i+3]
			e += ctx.randmem[i+4]
			f += ctx.randmem[i+5]
			g += ctx.randmem[i+6]
			h += ctx.randmem[i+7]

			// mix
			a -= e
			f ^= h >> 9
			h += a
			b -= f
			g ^= a << 9
			a += b
			c -= g
			h ^= b >> 23
			b += c
			d -= h
			a ^= c << 15
			c += d
			e -= a
			b ^= d >> 14
			d += e
			f -= b
			c ^= e << 20
			e += f
			g -= c
			d ^= f >> 17
			f += g
			h -= d
			e ^= g << 14
			g += h

			ctx.randmem[i] = a
			ctx.randmem[i+1] = b
			ctx.randmem[i+2] = c
			ctx.randmem[i+3] = d
			ctx.randmem[i+4] = e
			ctx.randmem[i+5] = f
			ctx.randmem[i+6] = g
			ctx.randmem[i+7] = h
		}
	} else {
		for i := 0; i < 256; i += 8 {
			// mix
			a -= e
			f ^= h >> 9
			h += a
			b -= f
			g ^= a << 9
			a += b
			c -= g
			h ^= b >> 23
			b += c
			d -= h
			a ^= c << 15
			c += d
			e -= a
			b ^= d >> 14
			d += e
			f -= b
			c ^= e << 20
			e += f
			g -= c
			d ^= f >> 17
			f += g
			h -= d
			e ^= g << 14
			g += h

			ctx.randmem[i] = a
			ctx.randmem[i+1] = b
			ctx.randmem[i+2] = c
			ctx.randmem[i+3] = d
			ctx.randmem[i+4] = e
			ctx.randmem[i+5] = f
			ctx.randmem[i+6] = g
			ctx.randmem[i+7] = h
		}
	}

	ctx.isaac64()
	ctx.randcnt = 256
}

func (ctx *Isaac64) next() uint64 {
	if ctx.randcnt == 0 {
		ctx.isaac64()
		ctx.randcnt = 255
		return ctx.randrsl[255]
	}

	ctx.randcnt--
	return ctx.randrsl[ctx.randcnt]
}
