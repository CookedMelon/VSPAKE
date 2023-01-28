package elliptic
import (
	"math/big"
)

func (cd*CurveDetail)Init(){
	cd.Name="secp224"
	cd.P,_= new(big.Int).SetString("26959946667150639794667015087019630673557916260026308143510066298881",10)
	cd.N,_= new(big.Int).SetString("26959946667150639794667015087019625940457807714424391721682722368061",10)
	cd.B,_ = new(big.Int).SetString("b4050a850c04b3abf54132565044b0b7d7bfd8ba270b39432355ffb4",16)
	cd.BasePoint.X,_ = new(big.Int).SetString("b70e0cbd6bb4bf7f321390b94a03c1d356c21122343280d6115c1d21",16)
	cd.BasePoint.Y,_ = new(big.Int).SetString("bd376388b5f723fb4c22dfe6cd4375a05a07476444d5819985007e34",16)
	cd.BitSize=224
}