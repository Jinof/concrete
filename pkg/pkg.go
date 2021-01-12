package pkg

import "fmt"

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
