package pkg

import (
	"errors"
	"fmt"
	"math"
)

const (
	// All const base on mm && KN/m^2

	// PKGα1 α1 1.0
	PKGα1 = 1.0
	// FC  11.9 N/mm^2  1.19 x 10 ^ 4 KN/m^2
	FC = 1.19 * 10000
	// PKGb b  1000 mm  1 m
	PKGb = 1
	// PKGh h 80 mm 0.08 m
	PKGh = 0.07
	// PKGh0  h - 0.02 m
	PKGh0 = PKGh - 0.02
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
	panic(fmt.Sprintf("[GetLocation] bad point %s", point))
}

// BasicCalculaterData def
type BasicCalculaterData []float64

// BasicCalculater calculate
func BasicCalculater(M float64) (bcal BasicCalculaterData) {
	αs := Calαs(M)
	pesi := CalPesi(αs)
	As := CalAs(pesi)
	return BasicCalculaterData{αs, pesi, As}
}

// Calαs cal αs
func Calαs(M float64) float64 {
	return math.Abs(M / (PKGα1 * FC * PKGb * PKGh0 * PKGh0))
}

// CalPesi cal pesi
func CalPesi(αs float64) float64 {
	return math.Abs(1 - math.Sqrt(1-2*αs))
}

// CalAs cal As
// m^2 / 10^6 -> mm^2
func CalAs(pesi float64) float64 {
	return math.Abs(pesi * PKGb * PKGα1 * FC * FY / math.Pow(10, 6))
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

// Reinforcement def
type Reinforcement struct {
	Counter []ReinforcementData
}

// ReinforcementData def
type ReinforcementData struct {
	rfType string // single or double
	space  float64
	as     float64
	d      float64
	d1     float64
	d2     float64
}

// Reinforcement def
var reinforcement Reinforcement

// NewReinforcement return reinforcement
func NewReinforcement() Reinforcement {
	return reinforcement
}

// CalRealSingleReinforcement cal with one Reinforcement
// s space (mm)
// diameter (mm)
// r real (mm^2) is the real reinforcement
// d (mm) is the diameter
func (rf *Reinforcement) CalRealSingleReinforcement(spaces []float64, diameters []float64, As float64) (err error) {
	const rfType = "single"
	var r, s float64
	exist := false
	for _, s = range spaces {
		for _, d := range diameters {
			r = realAsOfBoardForSingleReinforcement(s, d)
			if checkAs(r, As) {
				exist = true
				data := ReinforcementData{
					rfType: rfType,
					space:  s,
					as:     r,
					d:      d,
				}
				rf.Counter = append(rf.Counter, data)
			}
		}
	}
	if exist {
		return
	}
	return errors.New("bad cal realSingleReinforcement")
}

// CalRealDoubleReinforcement cal with one Reinforcement
// space (mm)
// diameter (mm)
// real (mm^2) is the real reinforcement
// d1 (mm) is the shorter diameter
// d2 (mm) is the longer diameter
func (rf *Reinforcement) CalRealDoubleReinforcement(spaces []float64, diameters []float64, As float64) (err error) {
	const rfType = "double"
	exist := false
	for _, s := range spaces {
		for i := 0; i < len(diameters)-1; i++ {
			cal := realAsOfBoardForDoubleReinforcement(s, diameters[i], diameters[i+1])
			if checkAs(cal, As) {
				exist = true
				data := ReinforcementData{
					rfType: rfType,
					space:  s,
					as:     cal,
					d1:     diameters[i],
					d2:     diameters[i+1],
				}
				rf.Counter = append(rf.Counter, data)
			}
			cal = realAsOfBoardForDoubleReinforcement(s, diameters[i+1], diameters[i])
			if checkAs(cal, As) {
				exist = true
				data := ReinforcementData{
					rfType: rfType,
					space:  s,
					as:     cal,
					d1:     diameters[i],
					d2:     diameters[i+1],
				}
				rf.Counter = append(rf.Counter, data)
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

// BestChoice select the best choice
func BestChoice(counters [][]ReinforcementData) {
	type recordData struct {
		point string
		ReinforcementData
	}
	// points := []string{A, FIRST, B, C, SECOND}

	// 以MA配筋为基准, 统计间距相同配筋
	tag := make([]bool, len(counters[1]))
	c := make(map[int]int, len(counters[1]))
	for i := 0; i < len(counters); i++ {
		for j := 0; j < len(counters[i]); j++ {
			for l, t := range tag {
				if !t {
					if counters[1][l] == counters[i][j] {
						c[l]++
					}
				}
			}
		}
	}

	
	bestExist := false
	bestCount := 0

	for k, v := range c {
		if v >= 5 {
			tag[k] = true
			bestCount = 5
			bestExist = true
		}
	}

	if !bestExist {
		for k, v := range c {
			if v >= 4 {
				tag[k] = true
				bestCount = 4
				bestExist = true
			}
		}
	}

	if bestExist {

		fmt.Printf("\n\n 最多重合配筋情况%d次 \n", bestCount)
		for i, v := range tag {
			if v {
				fmt.Printf("%v \n", counters[1][i])
			}
		}
		fmt.Println()
	}
}