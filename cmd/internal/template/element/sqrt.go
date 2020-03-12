package element

const Sqrt = `
func (z *{{.ElementName}}) Legendre() int {
	var l {{.ElementName}}
	// z^((q-1)/2)
	l.Exp(*z, {{range $i := .LegendreExponent}}
		{{$i}},{{end}}
	)
	
	if l.IsZero() {
		return 0
	} 

	// if l == 1
	if {{- range $i :=  reverse .NbWordsIndexesNoZero}}(l[{{$i}}] == {{index $.One $i}}) &&{{end}}(l[0] == {{index $.One 0}})  {
		return 1
	}
	return -1
}

// Sqrt z = √x mod q
// if the square root doesn't exist (x is not a square mod q)
// Sqrt leaves z unchanged and returns nil
func (z *{{.ElementName}}) Sqrt(x *{{.ElementName}}) *{{.ElementName}} {
	{{- if .SqrtQ3Mod4}}
		// q ≡ 3 (mod 4)
		// using  z ≡ ± x^((p+1)/4) (mod q)
		var y, square {{.ElementName}}
		y.Exp(*x, {{range $i := .SqrtQ3Mod4Exponent}}
			{{$i}},{{end}}
		)
		// TODO is this needed? seems cheaper than computing the Legendre symbol 
		square.Square(&y)
		if square.Equal(x) {
			return z.Set(&y)
		} else {
			return nil
		}
	{{- else if .SqrtTonelliShanks}}
		// q ≡ 1 (mod 4)
		// see modSqrtTonelliShanks in math/big/int.go
		// using https://www.maa.org/sites/default/files/pdf/upload_library/22/Polya/07468342.di020786.02p0470a.pdf

		var y, b,t, w  {{.ElementName}}
		// w = x^((s-1)/2))
		w.Exp(*x, {{range $i := .SqrtSMinusOneOver2}}
			{{$i}},{{end}}
		)

		// y = x^((s+1)/2)) = w * x
		y.Mul(x, &w)

		// b = x^s = w * w * x = y * x
		b.Mul(&w, &y)

		// g = nonResidue ^ s
		var g = {{.ElementName}}{
			{{- range $i := .SqrtG}}
			{{$i}},{{end}}
		}
		r := uint64({{.SqrtE}})

		// compute legendre symbol
		// t = x^((q-1)/2) = r-1 squaring of x^s
		t.Set(&b)
		for i:=uint64(0); i < r-1; i++ {
			t.Square(&t)
		}
		if t.IsZero() {
			return z.SetZero()
		}
		if !({{- range $i :=  reverse .NbWordsIndexesNoZero}}(t[{{$i}}] == {{index $.One $i}}) &&{{end}}(t[0] == {{index $.One 0}})) {
			// t != 1, we don't have a square root
			return nil
		}
		for {
			var m uint64
			t.Set(&b)

			// for t != 1
			for !({{- range $i :=  reverse .NbWordsIndexesNoZero}}(t[{{$i}}] == {{index $.One $i}}) &&{{end}}(t[0] == {{index $.One 0}})) {
				t.Square(&t)
				m++
			}

			if m == 0 {
				return z.Set(&y)
			}
			// t = g^(2^(r-m-1)) mod q
			ge := int(r - m - 1)
			t = g
			for ge > 0 {
				t.Square(&t)
				ge--
			}

			g.Square(&t)
			y.MulAssign(&t)
			b.MulAssign(&g)
			r = m
		}

	{{- else if .SqrtAtkin}}
		// q ≡ 5 (mod 8)
		// see modSqrt5Mod8Prime in math/big/int.go
		var one, alpha, beta, tx, square {{.ElementName}}
		one.SetOne()
		tx.Double(x)
		alpha.Exp(tx, {{range $i := .SqrtAtkinExponent}}
			{{$i}},{{end}}
		)
		beta.Square(&alpha).
			MulAssign(&tx).
			SubAssign(&one).
			MulAssign(x).
			MulAssign(&alpha)
		
		// TODO is this needed? seems cheaper than computing the Legendre symbol 
		square.Square(&beta)
		if square.Equal(x) {
			return z.Set(&beta)
		} else {
			return nil
		}
	{{- else}}
		panic("not implemented")	
	{{- end}}
}
`
