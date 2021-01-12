package board

import (
	"errors"
	"fmt"
	"log"
	"math"
)

const (
	// All const base on mm && KN/m^2
	// a1  1.0
	α1 = 1.0
	// FC  11.9 N/mm^2  1.19 x 10 ^ 4 KN/m^2
	FC = 1.19 * 10000
	// b  1000 mm  1 m
	b = 1
	// h 80 mm 0.08 m
	h = 0.07
	// h0  h - 0.02 m
	h0 = h - 0.02
	// FY  270 N/mm^2  2.7 x 10 ^5 KN/m^2
	FY = 2.7 * 100000

	// A .
	A = "A"
	// B .
	B = "B"
	// C .
	C = "C"

	// FIRST .
	FIRST = "1"
	// SECOND .
	SECOND = "2"
	// THIRD .
	THIRD = "3"

	// SURFACE  板顶
	SURFACE = "SURFACE"
	// BUTTOM  板底
	BUTTOM = "BUTTOM"
)

// CalBoard calculate board
func CalBoard() {
	// l01 = 1970 mm, l02 = 2000 mm
	l01 := 1.97
	l02 := 2

	// 板的永久荷载标准值
	// smsG 水磨石 0.65 KN/m^2
	// hntbG 70mm 钢筋混凝土板 0.07 * 25 KN/m^2
	// sjG 20 mm 混合砂浆 0.02 * 17 KN/m^2
	smsConsistantLoad := 0.65
	hntbConsistantLoad := 0.07 * 25
	sjConsistantLoad := 0.02 * 17
	sumCG := smsConsistantLoad + hntbConsistantLoad + sjConsistantLoad

	// 板可变荷载标准值 4.5 KN/m^2
	dynamicLoad := 4.5

	// 永久荷载分项系数 1.3
	// 可变荷载分项系数 1.5
	// g 永久荷载设计值
	g := sumCG * 1.3
	// q 可变荷载设计值
	q := dynamicLoad * 1.5
	// sumLoad 荷载总设计值
	sumLoad := g + q
	// sumLoad == 10.312
	// 取sumLoad = 10.3
	t := sumLoad
	sumLoad = 10.3
	fmt.Println("sumLoad 荷载总设计值", t, "实际取", sumLoad)

	var gq float64
	gq = sumLoad
	//
	MA := -gq * float64(l01*l01) / 16
	M1 := gq * float64(l01*l01) / 14
	MB := -gq * float64(l01*l01) / 11
	MC := -gq * float64(l02*l02) / 14
	// M2 = M3
	M2 := gq * float64(l02*l02) / 16

	Printer(A, MA)
	Printer(FIRST, M1)
	Printer(B, MB)
	Printer(C, MC)
	Printer(SECOND, M2)
}

// GetLocation return the location of the point
func GetLocation(point string) string {
	switch point {
	case A:
		return SURFACE
	case B:
		return SURFACE
	case C:
		return SURFACE
	case FIRST:
		return BUTTOM
	case SECOND:
		return BUTTOM
	}
	panic("[GetLocation] bad point")
}

// Printer print
func Printer(point string, M float64) {
	cal := Calculater(M)
	CalReinforcement(h, cal[2], SURFACE)
	fmt.Printf("M%s: %f, αs: %f, pesi: %f, As: %f \n", point, M, cal[0], cal[1], cal[2])
}

// Calculater calculate
func Calculater(M float64) [3]float64 {
	αs := Calαs(M)
	pesi := CalPesi(αs)
	As := CalAs(pesi)
	return [3]float64{αs, pesi, As}
}

// Calαs cal αs
func Calαs(M float64) float64 {
	return math.Abs(M / (α1 * FC * b * h0 * h0))
}

// CalPesi cal pesi
func CalPesi(αs float64) float64 {
	return math.Abs(1 - math.Sqrt(1-2*αs))
}

// CalAs cal As
// m^2 / 10^6 -> mm^2
func CalAs(pesi float64) float64 {
	return math.Abs(pesi * b * α1 * FC * FY / math.Pow(10, 6))
}

// CalReinforcement cal reinforcement
// h (mm) is the height of the board
// location is SURFACE or BUTTOM
// As (mm^2) is the calculated reinforcement
func CalReinforcement(h float64, As float64, location string) (err error) {

	var spaces = NewSpace(h)
	var diameters = NewDiameter(location)

	// 优先使用同一种钢筋
	err = realSingleReinforcement(spaces, diameters, As)
	if err != nil {
		log.Println("[CalReinforcement]", err)
	}
	// 当无法使用单种钢筋时，采用两种极差为1的钢筋
	err = realDoubleReinforcement(spaces, diameters, As)
	if err != nil {
		log.Fatal("[CalReinforcement]", err)
	}
	return
}

// realSingleReinforcement cal with one Reinforcement
// s space (mm)
// diameter (mm)
// r real (mm^2) is the real reinforcement
// d (mm) is the diameter
func realSingleReinforcement(spaces []float64, diameters []float64, As float64) (err error) {
	var r, s float64
	exist := false
	for _, s = range spaces {
		for _, d := range diameters {
			r = realAsOfBoardForSingleReinforcement(s, d)
			if checkAs(r, As) {
				exist = true
				fmt.Println("single", r, s, d)
			}
		}
	}
	if exist {
		return
	}
	return errors.New("bad cal realSingleReinforcement")
}

// realDoubleReinforcement cal with one Reinforcement
// space (mm)
// diameter (mm)
// real (mm^2) is the real reinforcement
// d1 (mm) is the shorter diameter
// d2 (mm) is the longer diameter
func realDoubleReinforcement(spaces []float64, diameters []float64, As float64) (err error) {
	exist := false
	for _, s := range spaces {
		for i := 0; i < len(diameters)-1; i++ {
			cal := realAsOfBoardForDoubleReinforcement(s, diameters[i], diameters[i+1])
			if checkAs(cal, As) {
				exist = true
				fmt.Println("double", cal, s, diameters[i], diameters[i+1])
			}
		}
	}
	if exist {
		return
	}
	return errors.New("bad cal realDoubleReinforcement")
}

// calas calculate the As (mm^2)
// d (mm) is the diameter of the reinforcement
// s (mm) is the space between reinforcement
func realAsOfBoardForDoubleReinforcement(s, d1, d2 float64) float64 {
	n := int(1000 / s)
	n1 := n / 2
	n2 := n - n1
	// let sinner reinforcement more
	if n1 < n2 {
		n1, n2 = n2, n1
	}
	return (math.Pi*math.Pow(d1, 2)/4)*float64(n1) + (math.Pi*math.Pow(d2, 2)/4)*float64(n2)
}

// calas calculate the As (mm^2)
// d (mm) is the diameter of the reinforcement
// s (mm) is the space between reinforcement
func realAsOfBoardForSingleReinforcement(s, d float64) float64 {
	return (math.Pi * math.Pow(d, 2) / 4) * (1000 / s)
}

// checkAs the real As should
// less than (1 - percent) * As
// and greater than (1 + percent) * As
func checkAs(cal float64, As float64) bool {
	const percent = 0.2
	if cal > As*(1-percent/4) && cal < As*(1+percent) {
		return true
	}
	return false
}

// NewDiameter generate a array of diameters for SURFACE OF BUTTOM
func NewDiameter(location string) (diameters []float64) {
	// 板的钢筋直径取 6, 8, 10, 12 mm
	// 板面钢筋 > 8 mm
	// 板底钢筋 > 6 mm
	if location == SURFACE {
		return []float64{8, 10, 12}
	} else if location == BUTTOM {
		return []float64{6, 8, 10, 12}
	}
	panic("wrong location which should be SURFACE or BUTTOM")
}

// NewSpace init the space bewteen reinforcement of board.
// h (mm) is the height of the board.
func NewSpace(h float64) (space []float64) {
	// return []float64{140.0}
	if h <= 150.0 {
		for i := 70.0; i < 200.0; i += 10.0 {
			space = append(space, i)
		}
	} else {
		for i := 70.0; i < 250.0 && i <= 1.5*h; i += 10.0 {
			space = append(space, i)
		}
	}
	return
}
