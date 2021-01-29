package grider

import (
	"fmt"

	"github.com/Jinof/concrete/bridge"
	"github.com/Jinof/concrete/pkg"
)

var (
	// G 永久荷载设计值
	G float64
	// Q 可变荷载设计值
	Q       float64
	gGirder float64
	selfG   float64
	qBridge float64
	h0      float64
	h01     float64
	h02     float64
	l01     float64
	l02     float64
	bf1     float64
)

// CalGirder calculate the girder
func CalGirder() {
	fmt.Printf("\n主梁设计\n")
	fmt.Printf("gBridge 次梁传来的永久荷载 %f, qBridge 次梁传来的可变荷载 %f, selfG 次梁自重 %f \n", gGirder, qBridge, selfG)
	fmt.Printf("G 永久荷载设计值 %f, Q 可变荷载设计值 %f \n", G, Q)

	M1 := 322.0
	MB := -351.0
	MB = MB + (G+Q)/2
	M21 := 167.4
	M22 := -47.9
	Calculator(pkg.FIRST, M1)
	Calculator(pkg.B, MB)
	Calculator(pkg.SECOND, M21)
	Calculator(pkg.SECOND, M22)
}

func init() {
	CalLoad()
}

// Calculator cal
func Calculator(point string, M float64) (rba []pkg.RealBridgeAs) {

	isT := pkg.CheckBridgeT(point)
	var tType int
	if isT {
		tType = pkg.CheckBridgeTtype(M, bf1, h0)
	}
	αs := pkg.CalBridgeαs(tType, M, bf1, h0)
	gama := pkg.CalGama(αs)
	As := pkg.CalGirderAs(gama, M, h0)
	if isT {
		fmt.Printf("M%s %f, h0 %f, tType 第%d种, αs %f, gama %f, As %f \n", point, M, h0, tType, αs, gama, As)
	} else {
		fmt.Printf("M%s %f, h0 %f, 支座截面, αs %f, gama %f, As %f \n", point, M, h0, αs, gama, As)
	}

	rba = CalReinforcement(As)
	fmt.Printf("以下为M%s的可能配筋 \n", point)
	for i, v := range rba {
		fmt.Println(i+1, v)
	}
	return
}

// CalReinforcement cal the reinforcement for bridge
func CalReinforcement(As float64) (rab []pkg.RealBridgeAs) {
	ns := pkg.CalBridgeReinforcementNum()
	diameters := pkg.NewBridgeDiameter()
	rab = pkg.CalBridgeRealAs(ns, diameters, As)
	return rab
}

// CalLoad calculate the load of bridge
func CalLoad() {
	// 次梁传来永久荷载
	gGirder = bridge.GetG() * 6.6
	// 主梁自重(含粉刷) [(0.65-0.07)x0.3x2.2x25+(0.65-0.07)x2x2.2x0.34]x1.3
	selfG = ((pkg.PKGGirderh-pkg.PKGh)*0.3*2.2*25 + (pkg.PKGGirderh-pkg.PKGh)*2*2.2*0.34) * 1.3
	// 小计
	G = gGirder + selfG

	// 次梁传来的可变荷载
	qBridge = bridge.GetQ()

	// 梁的可变荷载分项系数为 6.6
	// 可变荷载设计值
	Q = qBridge * 6.6

	l01 = pkg.L1 + 120 - 200
	l02 = pkg.L1

	bf1 = 2.1
	h01 = 0.595
	h02 = 0.570
	h0 = h02
}
