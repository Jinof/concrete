package bridge

import (
	"fmt"
	"math"
	"github.com/Jinof/concrete/board"
)

var (
	g       float64 // g 永久荷载设计值
	q       float64 // q 可变荷载设计值
	sumLoad float64 // sumLoad 荷载总设计值
)

// CalBridge calculate the bridge
func CalBridge() {
	// 主梁截面 250x650 mm^2
	// l01 边跨 l02 中间跨
	l01 := (6600 - 130 - 250 / 2) / 1000.0
	l02 := (6600 - 250) / 1000.0

	var gq float64
	gq = sumLoad

	MA := -gq * math.Pow(l01, 2) / 24
	M1 := gq * math.Pow(l01, 2) / 14
	MB := -gq * math.Pow(l01, 2) / 11
	MC := -gq * math.Pow(l02, 2) / 16
	// M2 = M3
	M2 := gq * math.Pow(l02, 2) / 14

	VA := 0.5 * gq * l01
	VBl := 0.55 * gq * l01
	// VBr == Vc
	VC := 0.55 * (gq) * l02

	fmt.Println("MA", MA)
	fmt.Println("M1", M1)
	fmt.Println("MB", MB)
	fmt.Println("MC", MC)
	fmt.Println("M2", M2)
	fmt.Println("VA", VA)
	fmt.Println("VB1", VBl)
	fmt.Println("VC", VC)
}

func init() {
	CalLoad()
}

// CalLoad calculate the load of bridge
func CalLoad() {
	// 板传来永久荷载
	gBoard := board.GetG()
	// 次梁自重 0.2x(0.5-0.07)x25x1.3
	selfG := 0.2 * (0.5 - 0.07) * 25 * 1.3
	// 次梁粉刷 0.02x(0.5-0.07)x2x17x1.3
	paintLayer := 0.02 * (0.5 - 0.07) * 2 * 17 * 1.3
	// 小计
	g = gBoard + selfG + paintLayer

	// 板传来的可变荷载
	qBoard := board.GetQ()

	// 梁的可变荷载分项系数为 2.2
	// 可变荷载设计值
	q = qBoard * 2.2

	// 荷载设计总值
	sumLoad = g + q
	fmt.Printf("g 永久荷载设计值 %f, q 可变荷载设计值 %f, sumLoad 荷载总设计值 %f \n", g, q, sumLoad)
}
