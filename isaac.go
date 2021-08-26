// Copyright 2021 skdltmxn. All rights reserved.
//
// isaac.go
//
// ISAAC random generator by Bob Jenkins

package isaac

import "unsafe"

// Isaac represents ISAAC random generator
type Isaac struct {
	randrsl [256]uint32
	randmem [256]uint32
	randcnt uint32
	aa      uint32
	bb      uint32
	cc      uint32
}

// NewIsaac returns a new instance of ISAAC.
func NewIsaac() *Isaac {
	return &Isaac{
		randmem: [256]uint32{},
		randrsl: [256]uint32{},
		randcnt: 0,
		aa:      0,
		bb:      0,
		cc:      0,
	}
}

// Seed initializes the state of ISAAC instance using given 64bit integer.
func (ctx *Isaac) Seed(seed int64) {
	ctx.randrsl[0] = uint32(seed)
	ctx.randrsl[1] = uint32(seed >> 32)
	ctx.randInit(true)
}

// SeedBytes initializes the state of ISAAC instance using given byte sequence.
func (ctx *Isaac) SeedBytes(seed []byte) {
	if len(seed) > 1024 {
		seed = seed[:1024]
	}

	unsafeCopy(unsafe.Pointer(&ctx.randrsl[0]), unsafe.Pointer(&seed[0]), len(seed))
	ctx.randInit(true)
}

// SeedString initializes the state of ISAAC instance using given string.
func (ctx *Isaac) SeedString(seed string) {
	ctx.SeedBytes([]byte(seed))
}

// Int63 returns a non-negative 63-bit integer as an int64.
func (ctx *Isaac) Int63() int64 {
	return int64(ctx.Uint64() & uintMask)
}

// Uint32 returns a random 32-bit unsigned integer.
func (ctx *Isaac) Uint32() uint32 {
	return ctx.next()
}

// Uint64 returns a random 64-bit unsigned integer.
func (ctx *Isaac) Uint64() uint64 {
	return (uint64(ctx.next()) << 32) | uint64(ctx.next())
}

// Int31 returns a non-negative 31-bit integer as an int32.
func (ctx *Isaac) Int31() int32 {
	return int32(ctx.Int63() >> 32)
}

// Int returns a non-negative integer as an int
func (ctx *Isaac) Int() int {
	u := uint(ctx.Uint64())
	return int(u << 1 >> 1)
}

func (ctx *Isaac) isaac() {
	var a, b, x uint32
	mm := ctx.randmem[:]
	r := ctx.randrsl[:]
	ctx.cc++
	a, b = ctx.aa, ctx.bb+ctx.cc

	for ii := 0; ii < 256; ii += 4 {
		var i uint8 = uint8(ii)

		x = mm[i]
		a = (a ^ (a << 13)) + mm[i+128]
		mm[i] = mm[(x>>2)&255] + a + b
		r[i] = mm[(mm[i]>>10)&255] + x
		b = r[i]

		x = mm[i+1]
		a = (a ^ (a >> 6)) + mm[i+129]
		mm[i+1] = mm[(x>>2)&255] + a + b
		r[i+1] = mm[(mm[i+1]>>10)&255] + x
		b = r[i+1]

		x = mm[i+2]
		a = (a ^ (a << 2)) + mm[i+130]
		mm[i+2] = mm[(x>>2)&255] + a + b
		r[i+2] = mm[(mm[i+2]>>10)&255] + x
		b = r[i+2]

		x = mm[i+3]
		a = (a ^ (a >> 16)) + mm[i+131]
		mm[i+3] = mm[(x>>2)&255] + a + b
		r[i+3] = mm[(mm[i+3]>>10)&255] + x
		b = r[i+3]
	}

	ctx.bb, ctx.aa = b, a
}

func (ctx *Isaac) randInit(flag bool) {
	var a, b, c, d, e, f, g, h uint32
	a, b, c, d, e, f, g, h = 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9

	// scramble
	for i := 0; i < 4; i++ {
		a ^= b << 11
		d += a
		b += c
		b ^= c >> 2
		e += b
		c += d
		c ^= d << 8
		f += c
		d += e
		d ^= e >> 16
		g += d
		e += f
		e ^= f << 10
		h += e
		f += g
		f ^= g >> 4
		a += f
		g += h
		g ^= h << 8
		b += g
		h += a
		h ^= a >> 9
		c += h
		a += b
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
			a ^= b << 11
			d += a
			b += c
			b ^= c >> 2
			e += b
			c += d
			c ^= d << 8
			f += c
			d += e
			d ^= e >> 16
			g += d
			e += f
			e ^= f << 10
			h += e
			f += g
			f ^= g >> 4
			a += f
			g += h
			g ^= h << 8
			b += g
			h += a
			h ^= a >> 9
			c += h
			a += b

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
			a ^= b << 11
			d += a
			b += c
			b ^= c >> 2
			e += b
			c += d
			c ^= d << 8
			f += c
			d += e
			d ^= e >> 16
			g += d
			e += f
			e ^= f << 10
			h += e
			f += g
			f ^= g >> 4
			a += f
			g += h
			g ^= h << 8
			b += g
			h += a
			h ^= a >> 9
			c += h
			a += b

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
			a ^= b << 11
			d += a
			b += c
			b ^= c >> 2
			e += b
			c += d
			c ^= d << 8
			f += c
			d += e
			d ^= e >> 16
			g += d
			e += f
			e ^= f << 10
			h += e
			f += g
			f ^= g >> 4
			a += f
			g += h
			g ^= h << 8
			b += g
			h += a
			h ^= a >> 9
			c += h
			a += b

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

	ctx.isaac()
	ctx.randcnt = 256
}

func (ctx *Isaac) next() uint32 {
	if ctx.randcnt == 0 {
		ctx.isaac()
		ctx.randcnt = 255
		return ctx.randrsl[255]
	}

	ctx.randcnt--
	return ctx.randrsl[ctx.randcnt]
}
