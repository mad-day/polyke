/*
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org>
*/

package poly

import "io"
import "encoding/binary"

type UPoly []uint64

// Calculates (a*b) mod Φ(x) where Φ(x) = x^len(u) - 1.
func (u UPoly) Mul(a, b UPoly) UPoly {
	for i := range u { u[i]=0 }
	l := len(u)
	for i,ai := range a {
		for j,bj := range b {
			// The result is mod Φ(x) as a side-effect of wrapping
			// the index mod len(u).
			u[(i+j)%l] += ai*bj
		}
	}
	return u
}

// Calculates (a*b) mod Φ(x) where Φ(x) = x^len(u) - 1.
// Every coefficient is mod q.
func (u UPoly) MulQ(a, b UPoly, q uint64) UPoly {
	for i := range u { u[i]=0 }
	l := len(u)
	for i,ai := range a {
		for j,bj := range b {
			// The result is mod Φ(x) as a side-effect of wrapping
			// the index mod len(u).
			k := (i+j)%l
			u[k] = (u[k] + ((ai*bj)%q))%q
		}
	}
	return u
}

func (u UPoly) Add(a, b UPoly) UPoly {
	for i := range u { u[i] = a[i]+b[i] }
	return u
}
func (u UPoly) AddQ(a, b UPoly, q uint64) UPoly {
	for i := range u { u[i] = (a[i]+b[i])%q }
	return u
}
func (u UPoly) Sample(r io.Reader) error {
	return binary.Read(r,binary.LittleEndian,u)
}
func (u UPoly) SampleEp(r io.Reader, e *error) UPoly {
	d := binary.Read(r,binary.LittleEndian,u)
	if e!=nil { *e = d }
	return u
}
func (u UPoly) Serialize(w io.Writer) error {
	return binary.Write(w,binary.LittleEndian,u)
}
func (u UPoly) Mods(q uint64) UPoly {
	for i,b := range u { u[i] = b%q}
	return u
}

func (u UPoly) Muls(q uint64) UPoly {
	for i,b := range u { u[i] = b*q}
	return u
}
func (u UPoly) Clone() UPoly {
	v := make(UPoly,len(u))
	copy(v,u)
	return v
}

