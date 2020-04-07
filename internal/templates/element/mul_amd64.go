package element

const MontgomeryMultiplicationAMD64 = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

func mulAsm{{.ElementName}}(res,y *{{.ElementName}})

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	res := *x
	mulAsm{{.ElementName}}(&res, y)
	z.Set(&res)
	return z
}

// MulAssign z = z * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	mulAsm{{.ElementName}}(z, x)
	return z 
}
`