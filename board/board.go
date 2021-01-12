package board

import (
	"fmt"
	"log"
	"math"

	"github.com/Jinof/concrete/pkg"
)

var (
	g       float64 // g 永久荷载设计值
	q       float64 // q 可变荷载设计值
	sumLoad float64 // sumLoad 荷载总设计值
	sumCG   float64 // 板的永久荷载标准值小计
)

// CalBoard calculate board
func CalBoard() {
	// l01 边跨, l02 中间跨
	// l01 = 1970 mm, l02 = 2000 mm
	l01 := 1.97
	l02 := 2.0

	var gq float64
	gq = sumLoad

	MA := -gq * math.Pow(l01, 2) / 16
	M1 := gq * math.Pow(l01, 2) / 14
	MB := -gq * math.Pow(l01, 2) / 11
	MC := -gq * math.Pow(l02, 2) / 14
	// M2 = M3
	M2 := gq * math.Pow(l02, 2) / 16

	Printer(pkg.A, MA)
	Printer(pkg.FIRST, M1)
	Printer(pkg.B, MB)
	Printer(pkg.C, MC)
	Printer(pkg.SECOND, M2)
}

// Printer print
func Printer(point string, M float64) {
	cal := pkg.Calculater(M)
	CalReinforcement(pkg.PKGh, cal[2], pkg.GetLocation(point))
	fmt.Printf("以上为M%s的可能配筋情况 \n", point)
	fmt.Printf("M%s: %f, αs: %f, pesi: %f, As: %f, Location: %s \n", point, M, cal[0], cal[1], cal[2], pkg.GetLocation(point))
}

// CalReinforcement cal reinforcement
// h (mm) is the height of the board
// location is SURFACE or BUTTOM
// As (mm^2) is the calculated reinforcement
func CalReinforcement(h float64, As float64, location string) (err error) {

	var spaces = pkg.NewSpace(h)
	var diameters = pkg.NewDiameter(location)

	// 优先使用同一种钢筋
	err = pkg.CalRealSingleReinforcement(spaces, diameters, As)
	if err != nil {
		log.Println("[CalReinforcement]", err)
	}
	// 当无法使用单种钢筋时，采用两种极差为1的钢筋
	err = pkg.CalRealDoubleReinforcement(spaces, diameters, As)
	if err != nil {
		log.Fatal("[CalReinforcement]", err)
	}
	return
}

// CalLoad calculate the load of board
func CalLoad() {
	// 板的永久荷载标准值
	// smsG 水磨石 0.65 KN/m^2
	// hntbG 70mm 钢筋混凝土板 0.07 * 25 KN/m^2
	// sjG 20 mm 混合砂浆 0.02 * 17 KN/m^2
	smsConsistantLoad := 0.65
	hntbConsistantLoad := 0.07 * 25
	sjConsistantLoad := 0.02 * 17
	sumCG = smsConsistantLoad + hntbConsistantLoad + sjConsistantLoad

	// 板可变荷载标准值 4.5 KN/m^2
	dynamicLoad := 4.5

	// 永久荷载分项系数 1.3
	// 可变荷载分项系数 1.5
	g = sumCG * 1.3
	q := dynamicLoad * 1.5
	sumLoad = g + q
	// sumLoad == 10.312
	// 取sumLoad = 10.3
	t := sumLoad
	sumLoad = 10.3

	fmt.Println("sumLoad 荷载总设计值", t, "实际取", sumLoad)
}

func init() {
	CalLoad()
}

// GetG return g
func GetG() float64 {
	return g
}