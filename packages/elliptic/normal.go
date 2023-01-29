package elliptic

import "math/big"

func GetEmptyCurvePoint() *CurvePoint {
	cp := new(CurvePoint)
	cp.X = new(big.Int)
	cp.Y = new(big.Int)
	return cp
}
