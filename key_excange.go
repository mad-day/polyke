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

package polyke

import "github.com/mad-day/polyke/poly"
import "io"
import "golang.org/x/crypto/sha3"
import "bytes"

const PolyLength = 16 // Key-Size = 1024 bit

func genPoly(s []byte) poly.UPoly {
	sb := sha3.NewShake256()
	sb.Write(s)
	return make(poly.UPoly,PolyLength).SampleEp(sb,nil)
}

func KeyPair(r io.Reader) (pub []byte,priv []byte,err error){
	var seed [32]byte
	var pkb [64]byte
	_,err = io.ReadFull(r,seed[:])
	if err!=nil { return }
	_,err = io.ReadFull(r,pkb[:])
	if err!=nil { return }
	Gen := genPoly(seed[:])
	iPriv := genPoly(pkb[:])
	
	{
		buf := bytes.NewBuffer(seed[:])
		make(poly.UPoly,PolyLength).Mul(Gen,iPriv).Serialize(buf)
		pub = buf.Bytes()
	}
	priv = pkb[:]//append(seed[:],pkb[:]...)
	return
}

func Encrypt(pub []byte,r io.Reader) (enc []byte, sk []byte, err error){
	var pkb [64]byte
	Gen := genPoly(pub[:32])
	PKPoly := make(poly.UPoly,PolyLength).SampleEp(bytes.NewReader(pub[32:]),&err)
	if err!=nil { return }
	_,err = io.ReadFull(r,pkb[:])
	if err!=nil { return }
	iPriv := genPoly(pkb[:])
	
	{
		buf := new(bytes.Buffer)
		t := make(poly.UPoly,PolyLength)
		
		t.Mul(Gen,iPriv).Serialize(buf)
		enc = buf.Bytes()
		*buf = bytes.Buffer{}
		t.Mul(PKPoly,iPriv).Serialize(buf)
		sk = buf.Bytes()
	}
	for i := range iPriv { iPriv[i]=0 }
	return
}

func Decrypt(enc, priv []byte) (sk []byte, err error) {
	EPoly := make(poly.UPoly,PolyLength).SampleEp(bytes.NewReader(enc),&err)
	if err!=nil { return }
	iPriv := genPoly(priv)
	
	buf := new(bytes.Buffer)
	make(poly.UPoly,PolyLength).Mul(EPoly,iPriv).Serialize(buf)
	sk = buf.Bytes()
	return
}

