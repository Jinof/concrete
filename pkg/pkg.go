package pkg

import (
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


// Calculater calculate
func Calculater(M float64) [3]float64 {
	αs := Calαs(M)
	pesi := CalPesi(αs)
	As := CalAs(pesi)
	return [3]float64{αs, pesi, As}
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