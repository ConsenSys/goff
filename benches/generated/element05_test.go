// Code generated by goff DO NOT EDIT
package generated

import (
	"crypto/rand"
	"math/big"
	"math/bits"
	mrand "math/rand"
	"testing"
)

func TestELEMENT05CorrectnessAgainstBigInt(t *testing.T) {
	modulus, _ := new(big.Int).SetString("919386282782426590101321146749114408200961198002380330475565544407301787190520578580447098192129", 10)
	cmpEandB := func(e *Element05, b *big.Int, name string) {
		var _e big.Int
		if e.FromMont().ToBigInt(&_e).Cmp(b) != 0 {
			t.Fatal(name, "failed")
		}
	}
	var modulusMinusOne, one big.Int
	one.SetUint64(1)

	modulusMinusOne.Sub(modulus, &one)

	var n int
	if testing.Short() {
		n = 10
	} else {
		n = 500
	}

	for i := 0; i < n; i++ {

		// sample 2 random big int
		b1, _ := rand.Int(rand.Reader, modulus)
		b2, _ := rand.Int(rand.Reader, modulus)
		rExp := mrand.Uint64()

		// adding edge cases
		// TODO need more edge cases
		switch i {
		case 0:
			rExp = 0
			b1.SetUint64(0)
		case 1:
			b2.SetUint64(0)
		case 2:
			b1.SetUint64(0)
			b2.SetUint64(0)
		case 3:
			rExp = 0
		case 4:
			rExp = 1
		case 5:
			rExp = ^uint64(0) // max uint
		case 6:
			rExp = 2
			b1.Set(&modulusMinusOne)
		case 7:
			b2.Set(&modulusMinusOne)
		case 8:
			b1.Set(&modulusMinusOne)
			b2.Set(&modulusMinusOne)
		}

		rbExp := new(big.Int).SetUint64(rExp)

		var bMul, bAdd, bSub, bDiv, bNeg, bLsh, bInv, bExp, bExp2, bSquare big.Int

		// e1 = mont(b1), e2 = mont(b2)
		var e1, e2, eMul, eAdd, eSub, eDiv, eNeg, eLsh, eInv, eExp, eExp2, eSquare, eMulAssign, eSubAssign, eAddAssign Element05
		e1.SetBigInt(b1)
		e2.SetBigInt(b2)

		// (e1*e2).FromMont() === b1*b2 mod q ... etc
		eSquare.Square(&e1)
		eMul.Mul(&e1, &e2)
		eMulAssign.Set(&e1)
		eMulAssign.MulAssign(&e2)
		eAdd.Add(&e1, &e2)
		eAddAssign.Set(&e1)
		eAddAssign.AddAssign(&e2)
		eSub.Sub(&e1, &e2)
		eSubAssign.Set(&e1)
		eSubAssign.SubAssign(&e2)
		eDiv.Div(&e1, &e2)
		eNeg.Neg(&e1)
		eInv.Inverse(&e1)
		eExp.Exp(e1, rExp)
		bits := b2.Bits()
		exponent := make([]uint64, len(bits))
		for k := 0; k < len(bits); k++ {
			exponent[k] = uint64(bits[k])
		}
		eExp2.Exp(e1, exponent...)
		eLsh.Double(&e1)

		// same operations with big int
		bAdd.Add(b1, b2).Mod(&bAdd, modulus)
		bMul.Mul(b1, b2).Mod(&bMul, modulus)
		bSquare.Mul(b1, b1).Mod(&bSquare, modulus)
		bSub.Sub(b1, b2).Mod(&bSub, modulus)
		bDiv.ModInverse(b2, modulus)
		bDiv.Mul(&bDiv, b1).
			Mod(&bDiv, modulus)
		bNeg.Neg(b1).Mod(&bNeg, modulus)

		bInv.ModInverse(b1, modulus)
		bExp.Exp(b1, rbExp, modulus)
		bExp2.Exp(b1, b2, modulus)
		bLsh.Lsh(b1, 1).Mod(&bLsh, modulus)

		cmpEandB(&eSquare, &bSquare, "Square")
		cmpEandB(&eMul, &bMul, "Mul")
		cmpEandB(&eMulAssign, &bMul, "MulAssign")
		cmpEandB(&eAdd, &bAdd, "Add")
		cmpEandB(&eAddAssign, &bAdd, "AddAssign")
		cmpEandB(&eSub, &bSub, "Sub")
		cmpEandB(&eSubAssign, &bSub, "SubAssign")
		cmpEandB(&eDiv, &bDiv, "Div")
		cmpEandB(&eNeg, &bNeg, "Neg")
		cmpEandB(&eInv, &bInv, "Inv")
		cmpEandB(&eExp, &bExp, "Exp")
		cmpEandB(&eExp2, &bExp2, "Exp multi words")
		cmpEandB(&eLsh, &bLsh, "Lsh")

		// legendre symbol
		if e1.Legendre() != big.Jacobi(b1, modulus) {
			t.Fatal("legendre symbol computation failed")
		}
		if e2.Legendre() != big.Jacobi(b2, modulus) {
			t.Fatal("legendre symbol computation failed")
		}

		// sqrt
		var eSqrt Element05
		var bSqrt big.Int
		bSqrt.ModSqrt(b1, modulus)
		eSqrt.Sqrt(&e1)
		cmpEandB(&eSqrt, &bSqrt, "Sqrt")
	}
}

func TestELEMENT05IsRandom(t *testing.T) {
	for i := 0; i < 50; i++ {
		var x, y Element05
		x.SetRandom()
		y.SetRandom()
		if x.Equal(&y) {
			t.Fatal("2 random numbers are unlikely to be equal")
		}
	}
}

// -------------------------------------------------------------------------------------------------
// benchmarks
// most benchmarks are rudimentary and should sample a large number of random inputs
// or be run multiple times to ensure it didn't measure the fastest path of the function

var benchResElement05 Element05

func BenchmarkInverseELEMENT05(b *testing.B) {
	var x Element05
	x.SetRandom()
	benchResElement05.SetRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchResElement05.Inverse(&x)
	}

}
func BenchmarkExpELEMENT05(b *testing.B) {
	var x Element05
	x.SetRandom()
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Exp(x, mrand.Uint64())
	}
}

func BenchmarkDoubleELEMENT05(b *testing.B) {
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Double(&benchResElement05)
	}
}

func BenchmarkAddELEMENT05(b *testing.B) {
	var x Element05
	x.SetRandom()
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Add(&x, &benchResElement05)
	}
}

func BenchmarkSubELEMENT05(b *testing.B) {
	var x Element05
	x.SetRandom()
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Sub(&x, &benchResElement05)
	}
}

func BenchmarkNegELEMENT05(b *testing.B) {
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Neg(&benchResElement05)
	}
}

func BenchmarkDivELEMENT05(b *testing.B) {
	var x Element05
	x.SetRandom()
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Div(&x, &benchResElement05)
	}
}

func BenchmarkFromMontELEMENT05(b *testing.B) {
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.FromMont()
	}
}

func BenchmarkToMontELEMENT05(b *testing.B) {
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.ToMont()
	}
}
func BenchmarkSquareELEMENT05(b *testing.B) {
	benchResElement05.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Square(&benchResElement05)
	}
}

func BenchmarkSqrtELEMENT05(b *testing.B) {
	var a Element05
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.Sqrt(&a)
	}
}

func BenchmarkMulAssignELEMENT05(b *testing.B) {
	x := Element05{
		6305763461833391547,
		1900997390964946128,
		17928209215758611772,
		6473583891360551725,
		2186045224142037132,
	}
	benchResElement05.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.MulAssign(&x)
	}
}

// Montgomery multiplication benchmarks
func (z *Element05) mulCIOS(x *Element05) *Element05 {

	var t [6]uint64
	var D uint64
	var m, C uint64
	// -----------------------------------
	// First loop

	C, t[0] = bits.Mul64(x[0], z[0])
	C, t[1] = madd1(x[0], z[1], C)
	C, t[2] = madd1(x[0], z[2], C)
	C, t[3] = madd1(x[0], z[3], C)
	C, t[4] = madd1(x[0], z[4], C)

	D = C

	// m = t[0]n'[0] mod W
	m = t[0] * 17076148859205665023

	// -----------------------------------
	// Second loop
	C = madd0(m, 16086029110469869825, t[0])

	C, t[0] = madd2(m, 2669506345493802173, t[1], C)

	C, t[1] = madd2(m, 17981300449585932002, t[2], C)

	C, t[2] = madd2(m, 1690474207027858972, t[3], C)

	C, t[3] = madd3(m, 7939974905350761517, t[4], C, t[5])

	t[4], t[5] = bits.Add64(D, C, 0)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[1], z[0], t[0])
	C, t[1] = madd2(x[1], z[1], t[1], C)
	C, t[2] = madd2(x[1], z[2], t[2], C)
	C, t[3] = madd2(x[1], z[3], t[3], C)
	C, t[4] = madd2(x[1], z[4], t[4], C)

	D = C

	// m = t[0]n'[0] mod W
	m = t[0] * 17076148859205665023

	// -----------------------------------
	// Second loop
	C = madd0(m, 16086029110469869825, t[0])

	C, t[0] = madd2(m, 2669506345493802173, t[1], C)

	C, t[1] = madd2(m, 17981300449585932002, t[2], C)

	C, t[2] = madd2(m, 1690474207027858972, t[3], C)

	C, t[3] = madd3(m, 7939974905350761517, t[4], C, t[5])

	t[4], t[5] = bits.Add64(D, C, 0)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[2], z[0], t[0])
	C, t[1] = madd2(x[2], z[1], t[1], C)
	C, t[2] = madd2(x[2], z[2], t[2], C)
	C, t[3] = madd2(x[2], z[3], t[3], C)
	C, t[4] = madd2(x[2], z[4], t[4], C)

	D = C

	// m = t[0]n'[0] mod W
	m = t[0] * 17076148859205665023

	// -----------------------------------
	// Second loop
	C = madd0(m, 16086029110469869825, t[0])

	C, t[0] = madd2(m, 2669506345493802173, t[1], C)

	C, t[1] = madd2(m, 17981300449585932002, t[2], C)

	C, t[2] = madd2(m, 1690474207027858972, t[3], C)

	C, t[3] = madd3(m, 7939974905350761517, t[4], C, t[5])

	t[4], t[5] = bits.Add64(D, C, 0)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[3], z[0], t[0])
	C, t[1] = madd2(x[3], z[1], t[1], C)
	C, t[2] = madd2(x[3], z[2], t[2], C)
	C, t[3] = madd2(x[3], z[3], t[3], C)
	C, t[4] = madd2(x[3], z[4], t[4], C)

	D = C

	// m = t[0]n'[0] mod W
	m = t[0] * 17076148859205665023

	// -----------------------------------
	// Second loop
	C = madd0(m, 16086029110469869825, t[0])

	C, t[0] = madd2(m, 2669506345493802173, t[1], C)

	C, t[1] = madd2(m, 17981300449585932002, t[2], C)

	C, t[2] = madd2(m, 1690474207027858972, t[3], C)

	C, t[3] = madd3(m, 7939974905350761517, t[4], C, t[5])

	t[4], t[5] = bits.Add64(D, C, 0)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[4], z[0], t[0])
	C, t[1] = madd2(x[4], z[1], t[1], C)
	C, t[2] = madd2(x[4], z[2], t[2], C)
	C, t[3] = madd2(x[4], z[3], t[3], C)
	C, t[4] = madd2(x[4], z[4], t[4], C)

	D = C

	// m = t[0]n'[0] mod W
	m = t[0] * 17076148859205665023

	// -----------------------------------
	// Second loop
	C = madd0(m, 16086029110469869825, t[0])

	C, t[0] = madd2(m, 2669506345493802173, t[1], C)

	C, t[1] = madd2(m, 17981300449585932002, t[2], C)

	C, t[2] = madd2(m, 1690474207027858972, t[3], C)

	C, t[3] = madd3(m, 7939974905350761517, t[4], C, t[5])

	t[4], t[5] = bits.Add64(D, C, 0)

	if t[5] != 0 {
		// we need to reduce, we have a result on 6 words
		var b uint64
		z[0], b = bits.Sub64(t[0], 16086029110469869825, 0)
		z[1], b = bits.Sub64(t[1], 2669506345493802173, b)
		z[2], b = bits.Sub64(t[2], 17981300449585932002, b)
		z[3], b = bits.Sub64(t[3], 1690474207027858972, b)
		z[4], _ = bits.Sub64(t[4], 7939974905350761517, b)
		return z
	}

	// copy t into z
	z[0] = t[0]
	z[1] = t[1]
	z[2] = t[2]
	z[3] = t[3]
	z[4] = t[4]

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[4] < 7939974905350761517 || (z[4] == 7939974905350761517 && (z[3] < 1690474207027858972 || (z[3] == 1690474207027858972 && (z[2] < 17981300449585932002 || (z[2] == 17981300449585932002 && (z[1] < 2669506345493802173 || (z[1] == 2669506345493802173 && (z[0] < 16086029110469869825))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 16086029110469869825, 0)
		z[1], b = bits.Sub64(z[1], 2669506345493802173, b)
		z[2], b = bits.Sub64(z[2], 17981300449585932002, b)
		z[3], b = bits.Sub64(z[3], 1690474207027858972, b)
		z[4], _ = bits.Sub64(z[4], 7939974905350761517, b)
	}
	return z
}

func (z *Element05) mulNoCarry(x *Element05) *Element05 {

	var t [5]uint64
	var c [3]uint64
	{
		// round 0
		v := z[0]
		c[1], c[0] = bits.Mul64(v, x[0])
		m := c[0] * 17076148859205665023
		c[2] = madd0(m, 16086029110469869825, c[0])
		c[1], c[0] = madd1(v, x[1], c[1])
		c[2], t[0] = madd2(m, 2669506345493802173, c[2], c[0])
		c[1], c[0] = madd1(v, x[2], c[1])
		c[2], t[1] = madd2(m, 17981300449585932002, c[2], c[0])
		c[1], c[0] = madd1(v, x[3], c[1])
		c[2], t[2] = madd2(m, 1690474207027858972, c[2], c[0])
		c[1], c[0] = madd1(v, x[4], c[1])
		t[4], t[3] = madd3(m, 7939974905350761517, c[0], c[2], c[1])
	}
	{
		// round 1
		v := z[1]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 17076148859205665023
		c[2] = madd0(m, 16086029110469869825, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2669506345493802173, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 17981300449585932002, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, 1690474207027858972, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		t[4], t[3] = madd3(m, 7939974905350761517, c[0], c[2], c[1])
	}
	{
		// round 2
		v := z[2]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 17076148859205665023
		c[2] = madd0(m, 16086029110469869825, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2669506345493802173, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 17981300449585932002, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, 1690474207027858972, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		t[4], t[3] = madd3(m, 7939974905350761517, c[0], c[2], c[1])
	}
	{
		// round 3
		v := z[3]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 17076148859205665023
		c[2] = madd0(m, 16086029110469869825, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 2669506345493802173, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 17981300449585932002, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], t[2] = madd2(m, 1690474207027858972, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		t[4], t[3] = madd3(m, 7939974905350761517, c[0], c[2], c[1])
	}
	{
		// round 4
		v := z[4]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 17076148859205665023
		c[2] = madd0(m, 16086029110469869825, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], z[0] = madd2(m, 2669506345493802173, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], z[1] = madd2(m, 17981300449585932002, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		c[2], z[2] = madd2(m, 1690474207027858972, c[2], c[0])
		c[1], c[0] = madd2(v, x[4], c[1], t[4])
		z[4], z[3] = madd3(m, 7939974905350761517, c[0], c[2], c[1])
	}

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[4] < 7939974905350761517 || (z[4] == 7939974905350761517 && (z[3] < 1690474207027858972 || (z[3] == 1690474207027858972 && (z[2] < 17981300449585932002 || (z[2] == 17981300449585932002 && (z[1] < 2669506345493802173 || (z[1] == 2669506345493802173 && (z[0] < 16086029110469869825))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 16086029110469869825, 0)
		z[1], b = bits.Sub64(z[1], 2669506345493802173, b)
		z[2], b = bits.Sub64(z[2], 17981300449585932002, b)
		z[3], b = bits.Sub64(z[3], 1690474207027858972, b)
		z[4], _ = bits.Sub64(z[4], 7939974905350761517, b)
	}
	return z
}

func (z *Element05) mulFIPS(x *Element05) *Element05 {

	var p [5]uint64
	var t, u, v uint64
	u, v = bits.Mul64(z[0], x[0])
	p[0] = v * 17076148859205665023
	u, v, _ = madd(p[0], 16086029110469869825, 0, u, v)
	t, u, v = madd(z[0], x[1], 0, u, v)
	t, u, v = madd(p[0], 2669506345493802173, t, u, v)
	t, u, v = madd(z[1], x[0], t, u, v)
	p[1] = v * 17076148859205665023
	u, v, _ = madd(p[1], 16086029110469869825, t, u, v)
	t, u, v = madd(z[0], x[2], 0, u, v)
	t, u, v = madd(p[0], 17981300449585932002, t, u, v)
	t, u, v = madd(z[1], x[1], t, u, v)
	t, u, v = madd(p[1], 2669506345493802173, t, u, v)
	t, u, v = madd(z[2], x[0], t, u, v)
	p[2] = v * 17076148859205665023
	u, v, _ = madd(p[2], 16086029110469869825, t, u, v)
	t, u, v = madd(z[0], x[3], 0, u, v)
	t, u, v = madd(p[0], 1690474207027858972, t, u, v)
	t, u, v = madd(z[1], x[2], t, u, v)
	t, u, v = madd(p[1], 17981300449585932002, t, u, v)
	t, u, v = madd(z[2], x[1], t, u, v)
	t, u, v = madd(p[2], 2669506345493802173, t, u, v)
	t, u, v = madd(z[3], x[0], t, u, v)
	p[3] = v * 17076148859205665023
	u, v, _ = madd(p[3], 16086029110469869825, t, u, v)
	t, u, v = madd(z[0], x[4], 0, u, v)
	t, u, v = madd(p[0], 7939974905350761517, t, u, v)
	t, u, v = madd(z[1], x[3], t, u, v)
	t, u, v = madd(p[1], 1690474207027858972, t, u, v)
	t, u, v = madd(z[2], x[2], t, u, v)
	t, u, v = madd(p[2], 17981300449585932002, t, u, v)
	t, u, v = madd(z[3], x[1], t, u, v)
	t, u, v = madd(p[3], 2669506345493802173, t, u, v)
	t, u, v = madd(z[4], x[0], t, u, v)
	p[4] = v * 17076148859205665023
	u, v, _ = madd(p[4], 16086029110469869825, t, u, v)
	t, u, v = madd(z[1], x[4], 0, u, v)
	t, u, v = madd(p[1], 7939974905350761517, t, u, v)
	t, u, v = madd(z[2], x[3], t, u, v)
	t, u, v = madd(p[2], 1690474207027858972, t, u, v)
	t, u, v = madd(z[3], x[2], t, u, v)
	t, u, v = madd(p[3], 17981300449585932002, t, u, v)
	t, u, v = madd(z[4], x[1], t, u, v)
	u, v, p[0] = madd(p[4], 2669506345493802173, t, u, v)
	t, u, v = madd(z[2], x[4], 0, u, v)
	t, u, v = madd(p[2], 7939974905350761517, t, u, v)
	t, u, v = madd(z[3], x[3], t, u, v)
	t, u, v = madd(p[3], 1690474207027858972, t, u, v)
	t, u, v = madd(z[4], x[2], t, u, v)
	u, v, p[1] = madd(p[4], 17981300449585932002, t, u, v)
	t, u, v = madd(z[3], x[4], 0, u, v)
	t, u, v = madd(p[3], 7939974905350761517, t, u, v)
	t, u, v = madd(z[4], x[3], t, u, v)
	u, v, p[2] = madd(p[4], 1690474207027858972, t, u, v)
	t, u, v = madd(z[4], x[4], t, u, v)
	u, v, p[3] = madd(p[4], 7939974905350761517, t, u, v)

	p[4] = v
	z[4] = p[4]
	z[3] = p[3]
	z[2] = p[2]
	z[1] = p[1]
	z[0] = p[0]
	// copy(z[:], p[:])

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[4] < 7939974905350761517 || (z[4] == 7939974905350761517 && (z[3] < 1690474207027858972 || (z[3] == 1690474207027858972 && (z[2] < 17981300449585932002 || (z[2] == 17981300449585932002 && (z[1] < 2669506345493802173 || (z[1] == 2669506345493802173 && (z[0] < 16086029110469869825))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 16086029110469869825, 0)
		z[1], b = bits.Sub64(z[1], 2669506345493802173, b)
		z[2], b = bits.Sub64(z[2], 17981300449585932002, b)
		z[3], b = bits.Sub64(z[3], 1690474207027858972, b)
		z[4], _ = bits.Sub64(z[4], 7939974905350761517, b)
	}
	return z
}

func BenchmarkMulCIOSELEMENT05(b *testing.B) {
	x := Element05{
		6305763461833391547,
		1900997390964946128,
		17928209215758611772,
		6473583891360551725,
		2186045224142037132,
	}
	benchResElement05.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.mulCIOS(&x)
	}
}

func BenchmarkMulFIPSELEMENT05(b *testing.B) {
	x := Element05{
		6305763461833391547,
		1900997390964946128,
		17928209215758611772,
		6473583891360551725,
		2186045224142037132,
	}
	benchResElement05.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.mulFIPS(&x)
	}
}

func BenchmarkMulNoCarryELEMENT05(b *testing.B) {
	x := Element05{
		6305763461833391547,
		1900997390964946128,
		17928209215758611772,
		6473583891360551725,
		2186045224142037132,
	}
	benchResElement05.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement05.mulNoCarry(&x)
	}
}
