package elliptic

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
)

type JacobianPoint struct {
	X, Y, Z *big.Int
}
type CurvePoint struct {
	X, Y *big.Int
}
type CurveDetail struct {
	p         *big.Int    // the order of the underlying field
	N         *big.Int    // the order of the base point
	B         *big.Int    // the constant of the curve equation
	BasePoint CurvePoint  // (x,y) of the base point
	BitSize   int         // the size of the underlying field
	Name      string      // the canonical name of the curve
	P, Q      *CurvePoint //P,Q为椭圆曲线上的点
}

// type cd interface {

//		IfOnCurve(cp*CurvePoint) bool
//		Mult(cp*CurvePoint, k*big.Int) (ans*CurvePoint)
//		BaseMult(k*big.Int)(ans*CurvePoint)
//		Add(cp1,cp2 *CurvePoint)(ans*CurvePoint)
//		Double(cp *CurvePoint)(ans*CurvePoint)
//		polynomial(x *big.Int) *big.Int
//	}
//
// get x^3 - 3x + b
func (curve *CurveDetail) polynomial(x *big.Int) *big.Int {
	xxx := new(big.Int).Mul(x, x)
	xxx.Mul(xxx, x)
	x3 := new(big.Int).Lsh(x, 1)
	x3.Add(x3, x)

	ans := new(big.Int).Sub(xxx, x3)
	ans.Add(ans, curve.B)
	ans.Mod(ans, curve.p)

	return ans
}
func (curve *CurveDetail) IfOnCurve(cp *CurvePoint) bool {
	if cp.X.Sign() < 0 || cp.X.Cmp(curve.p) >= 0 ||
		cp.Y.Sign() < 0 || cp.Y.Cmp(curve.p) >= 0 {
		return false
	}

	// check y² = x³ - 3x + b
	yy := new(big.Int).Mul(cp.Y, cp.Y)
	yy.Mod(yy, curve.p)

	return yy.Cmp(curve.polynomial(cp.X)) == 0
}

func getZ(cp *CurvePoint) (z *big.Int) {
	z = new(big.Int)
	if cp.X.Sign() != 0 || cp.Y.Sign() != 0 {
		z.SetInt64(1)
		return z
	}
	return z
}
func (curve *CurveDetail) Jacobian2Curve(jp *JacobianPoint) (cp *CurvePoint) {
	cp = new(CurvePoint)
	// 无穷远点返回为（0,0）
	if jp.Z.Sign() == 0 {
		// cp.X , cp.Y= new(big.Int) , new(big.Int)
		return
	}
	zin := new(big.Int).ModInverse(jp.Z, curve.p)
	zin2 := new(big.Int).Mul(zin, zin)
	cp.X = new(big.Int).Mul(jp.X, zin2)
	cp.X.Mod(cp.X, curve.p)
	zin3 := new(big.Int).Mul(zin2, zin)
	cp.Y = new(big.Int).Mul(jp.Y, zin3)
	cp.Y.Mod(cp.Y, curve.p)
	return
}

// 求Jacobian加重射影坐标
func (curve *CurveDetail) JacobianAdd(jp1, jp2 *JacobianPoint) (ans *JacobianPoint) {
	ans = new(JacobianPoint)
	if jp1.Z.Sign() == 0 {
		data, _ := json.Marshal(jp2)
		json.Unmarshal(data, ans)
		return
	}
	if jp2.Z.Sign() == 0 {
		data, _ := json.Marshal(jp1)
		json.Unmarshal(data, ans)
		return
	}
	z1z1 := new(big.Int).Mul(jp1.Z, jp1.Z)
	z1z1.Mod(z1z1, curve.p)
	z2z2 := new(big.Int).Mul(jp2.Z, jp2.Z)
	z2z2.Mod(z2z2, curve.p)
	t1 := new(big.Int).Mul(jp1.X, z2z2)
	t1.Mod(t1, curve.p)
	t2 := new(big.Int).Mul(jp2.X, z1z1)
	t2.Mod(t2, curve.p)
	t3 := new(big.Int).Sub(t2, t1)
	//求mod消耗太大了
	s1 := t3.Sign()
	if t3.Sign() < 0 {
		t3.Add(t3, curve.p)
	}
	i := new(big.Int).Lsh(t3, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(t3, i)

	r1 := new(big.Int).Mul(jp1.Y, jp2.Z)
	r1.Mul(r1, z2z2)
	r1.Mod(r1, curve.p)
	r2 := new(big.Int).Mul(jp2.Y, jp1.Z)
	r2.Mul(r2, z1z1)
	r2.Mod(r2, curve.p)
	r3 := new(big.Int).Sub(r2, r1)
	s2 := r3.Sign()
	if r3.Sign() < 0 {
		r3.Add(r3, curve.p)
	}

	if s1 == 0 && s2 == 0 {
		ans = curve.JacobianDouble(jp1)
		return
	}
	r3.Lsh(r3, 1)
	v := new(big.Int).Mul(t1, i)
	ans.X = new(big.Int)
	ans.X.Set(r3)
	ans.X.Mul(ans.X, ans.X)
	ans.X.Sub(ans.X, j)
	ans.X.Sub(ans.X, v)
	ans.X.Sub(ans.X, v)
	ans.X.Mod(ans.X, curve.p)

	ans.Y = new(big.Int)
	ans.Y.Set(r3)
	v.Sub(v, ans.X)
	ans.Y.Mul(ans.Y, v)
	r1.Mul(r1, j)
	r1.Lsh(r1, 1)
	ans.Y.Sub(ans.Y, r1)
	ans.Y.Mod(ans.Y, curve.p)

	ans.Z = new(big.Int)
	ans.Z.Add(jp1.Z, jp2.Z)
	ans.Z.Mul(ans.Z, ans.Z)
	ans.Z.Sub(ans.Z, z1z1)
	ans.Z.Sub(ans.Z, z2z2)
	ans.Z.Mul(ans.Z, t3)
	ans.Z.Mod(ans.Z, curve.p)
	return
}
func (curve *CurveDetail) JacobianDouble(jp *JacobianPoint) (ans *JacobianPoint) {
	ans = new(JacobianPoint)
	zz := new(big.Int).Mul(jp.Z, jp.Z)
	zz.Mod(zz, curve.p)
	yy := new(big.Int).Mul(jp.Y, jp.Y)
	yy.Mod(yy, curve.p)
	u1 := new(big.Int).Sub(jp.X, zz)
	if u1.Sign() == -1 {
		u1.Add(u1, curve.p)
	}
	u2 := new(big.Int).Add(jp.X, zz)
	u1.Mul(u1, u2)
	u2.Set(u1)
	u1.Lsh(u1, 1)
	u1.Add(u1, u2)

	v1 := u2.Mul(jp.X, yy)
	v8 := new(big.Int).Lsh(v1, 3)
	v8.Mod(v8, curve.p)

	ans.X = new(big.Int).Mul(u1, u1)
	ans.X.Sub(ans.X, v8)
	ans.X.Mod(ans.X, curve.p)

	ans.Z = new(big.Int).Add(jp.Y, jp.Z)
	ans.Z.Mul(ans.Z, ans.Z)
	ans.Z.Sub(ans.Z, yy)
	ans.Z.Sub(ans.Z, zz)
	ans.Z.Mod(ans.Z, curve.p)

	v1.Lsh(v1, 2)
	v1.Sub(v1, ans.X)
	if v1.Sign() < 0 {
		v1.Add(v1, curve.p)
	}
	ans.Y = new(big.Int).Mul(u1, v1)
	yy.Mul(yy, yy)
	yy.Lsh(yy, 3)
	yy.Mod(yy, curve.p)
	ans.Y.Sub(ans.Y, yy)
	ans.Y.Mod(ans.Y, curve.p)
	return
}

func (curve *CurveDetail) Add(cp1, cp2 *CurvePoint) (ans *CurvePoint) {
	z1 := getZ(cp1)
	z2 := getZ(cp2)
	jp1 := new(JacobianPoint)
	jp1.X = cp1.X
	jp1.Y = cp1.Y
	jp1.Z = z1
	jp2 := new(JacobianPoint)
	jp2.X = cp2.X
	jp2.Y = cp2.Y
	jp2.Z = z2
	ans = curve.Jacobian2Curve(curve.JacobianAdd(jp1, jp2))
	return
}

func (curve *CurveDetail) Double(cp *CurvePoint) (ans *CurvePoint) {
	z := getZ(cp)
	jp := new(JacobianPoint)
	jp.X = cp.X
	jp.Y = cp.Y
	jp.Z = z
	ans = curve.Jacobian2Curve(curve.JacobianDouble(jp))
	return
}

func (curve *CurveDetail) Mult(cp *CurvePoint, k []byte) (ans *CurvePoint) {
	B := new(JacobianPoint)
	B.X = cp.X
	B.Y = cp.Y
	B.Z = new(big.Int).SetInt64(1)
	nB := new(JacobianPoint)
	nB.X = new(big.Int)
	nB.Y = new(big.Int)
	nB.Z = new(big.Int)
	for _, byte := range k {
		for i := 0; i < 8; i++ {
			nB = curve.JacobianDouble(nB)
			if byte&0x80 == 0x80 {
				nB = curve.JacobianAdd(B, nB)
			}
			byte = byte << 1
		}
	}
	return curve.Jacobian2Curve(nB)
}
func (curve *CurveDetail) BaseMult(k []byte) (ans *CurvePoint) {
	GP := new(CurvePoint)
	GP.X = curve.BasePoint.X
	GP.Y = curve.BasePoint.Y
	return curve.Mult(GP, k)
}

func (curve *CurveDetail) GetRandPoint() (ans *CurvePoint) {
	randInt := make([]byte, 32)
	rand.Read(randInt)
	ans = curve.BaseMult(randInt)
	return
}
func (curve *CurveDetail) GetNeg(cp *CurvePoint) (ans *CurvePoint) {
	cp.Y = new(big.Int).Mod(new(big.Int).Neg(cp.Y), curve.p)
	return
}
