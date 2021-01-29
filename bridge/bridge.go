package bridge

import (
	"fmt"
	"math"

	"github.com/Jinof/concrete/board"
	"github.com/Jinof/concrete/pkg"
)

var (
	g          float64 // g 永久荷载设计值
	q          float64 // q 可变荷载设计值
	sumLoad    float64 // sumLoad 荷载总设计值
	gBoard     float64
	selfG      float64
	paintLayer float64
	qBoard     float64
	h0         float64
	h01        float64
	h02        float64
)

// CalBridge calculate the bridge
func CalBridge() {
	fmt.Println("一排纵向钢筋 h0", h01)
	fmt.Println("两排纵向钢筋 h0", h02)
	fmt.Printf("gBoard 板传来的永久荷载 %f, selfG 次梁自重 %f, 次梁粉刷 %f \n", gBoard, selfG, paintLayer)
	fmt.Printf("g 永久荷载设计值 %f, q 可变荷载设计值 %f, sumLoad 荷载总设计值 %f \n", g, q, sumLoad)
	// 主梁截面 250x650 mm^2
	// l01 边跨 l02 中间跨
	l01 := (pkg.L1*1000 - 130 - 250/2) / 1000.0
	l02 := (pkg.L1*1000 - 250) / 1000.0
	fmt.Printf("l01 %f, l02 %f \n", l01*math.Pow(10, 3), l02*math.Pow(10, 3))

	var gq float64
	gq = sumLoad

	MA := -gq * math.Pow(l01, 2) / 24
	M1 := gq * math.Pow(l01, 2) / 14
	MB := -gq * math.Pow(l01, 2) / 11
	// M2 = M3
	M2 := gq * math.Pow(l02, 2) / 16
	MC := -gq * math.Pow(l02, 2) / 14

	VA := 0.5 * gq * l01
	VBl := 0.55 * gq * l01
	// VBr == Vc
	VC := 0.55 * (gq) * l02

	fmt.Println("MA", MA)
	fmt.Println("M1", M1)
	fmt.Println("MB", MB)
	fmt.Println("M2", M2)
	fmt.Println("MC", MC)
	fmt.Println("VA", VA)
	fmt.Println("VB1", VBl)
	fmt.Println("VC", VC)

	var rbas []pkg.RealBridgeAs
	counter := [][]pkg.RealBridgeAs{}

	rbas = Calculator(pkg.A, MA)
	counter = append(counter, rbas)

	rbas = Calculator(pkg.FIRST, M1)
	counter = append(counter, rbas)

	rbas = Calculator(pkg.B, MB)
	counter = append(counter, rbas)

	rbas = Calculator(pkg.SECOND, M2)
	counter = append(counter, rbas)

	rbas = Calculator(pkg.C, MC)
	counter = append(counter, rbas)

	pkg.BestBridgeReinforcementChoice(counter)
}

// Calculator cal
func Calculator(point string, M float64) (rba []pkg.RealBridgeAs) {
	// 翼缘宽度
	bf1 := CalFlangeWidth() / math.Pow(10, 3)

	isT := pkg.CheckBridgeT(point)
	var tType int
	if isT {
		tType = pkg.CheckBridgeTtype(M, bf1, h0)
	}
	αs := pkg.CalBridgeαs(tType, M, bf1, h0)
	pesi := pkg.CalPesi(αs)
	As := pkg.CalBridgeAs(tType, pesi, bf1, h0)
	if isT {
		fmt.Printf("M%s h0 %f, tType 第%d种, αs %f, pesi %f, As %f \n", point, h0, tType, αs, pesi, As)
	} else {
		fmt.Printf("M%s h0 %f, 支座截面, αs %f, pesi %f, As %f \n", point, h0, αs, pesi, As)
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

// CalFlangeWidth cal flange width
// bf1 (mm)
func CalFlangeWidth() float64 {
	// 翼缘宽度取 l/3, b + sn, b + 12hf 中的最小值
	return min(min(pkg.L1*math.Pow(10, 3)/3.0, 200+2000), 200.0+12*pkg.PKGhf*math.Pow(10, 3))
}

func min(a1, a2 float64) float64 {
	if a1 > a2 {
		return a2
	}
	return a1
}

func init() {
	CalLoad()
}

// CalLoad calculate the load of bridge
func CalLoad() {
	// 板传来永久荷载
	gBoard = board.GetG() * 2.2
	// 次梁自重 0.2x(0.5-0.07)x25x1.3
	selfG = 0.2 * (0.5 - 0.07) * 25 * 1.3
	// 次梁粉刷 0.02x(0.5-0.07)x2x17x1.3
	paintLayer = 0.02 * (0.5 - 0.07) * 2 * 17 * 1.3
	// 小计
	g = gBoard + selfG + paintLayer

	// 板传来的可变荷载
	qBoard = board.GetQ()

	// 梁的可变荷载分项系数为 2.2
	// 可变荷载设计值
	q = qBoard * 2.2

	// 荷载设计总值
	sumLoad = g + q

	// 一排纵向钢筋直径
	h01 = pkg.CalSingleLayerReinforcementH0()
	h0 = h01
	// 两排纵向钢筋直径
	h02 = pkg.CalDoubleLayerReinforcementH0()
	h0 = h02
	// 优先用一排纵向钢筋
	h0 = h01
}

// GetG return g
func GetG() float64 {
	return g
}

// GetQ return q
func GetQ() float64 {
	return q
}
