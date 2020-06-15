package cmd

import (
	"fmt"

	"github.com/consensys/bavard"
)

func generateSubASM(b *bavard.Assembly, F *field) error {

	b.FuncHeader("sub"+F.ElementName, 24)
	// registers
	b.Reset()
	var regX bavard.Register

	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
	}

	regX = b.PopRegister()
	regY := b.PopRegister()
	b.MOVQ("x+8(FP)", regX, "dereference x")
	b.MOVQ("y+16(FP)", regY, "dereference y")
	// z = x - y mod q

	for i := 0; i < F.NbWords; i++ {
		b.MOVQ(regX.At(i), regT[i], fmt.Sprintf("t[%d] = x[%d]", i, i))
	}

	b.XORQ(bavard.DX, bavard.DX)

	b.SUBQ(regY.At(0), regT[0])
	for i := 1; i < F.NbWords; i++ {
		b.SBBQ(regY.At(i), regT[i])
	}

	// reduction, if borrow is set (CMOVQCS)
	b.PushRegister(regY)
	regQ := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regQ[i] = b.PopRegister()
		b.MOVQ(F.Q[i], regQ[i])
	}
	for i := 0; i < F.NbWords; i++ {
		b.CMOVQCC(bavard.DX, regQ[i])
	}
	b.ADDQ(regQ[0], regT[0])
	for i := 1; i < F.NbWords; i++ {
		b.ADCQ(regQ[i], regT[i])
	}

	b.MOVQ("res+0(FP)", regX, "dereference res")

	for i := 0; i < F.NbWords; i++ {
		b.MOVQ(regT[i], regX.At(i))
	}

	b.RET()

	return nil
}
