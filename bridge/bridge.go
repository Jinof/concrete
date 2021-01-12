package bridge

import (
	"fmt"

	"github.com/Jinof/concrete/board"
)

var (
	g       float64 // g 永久荷载设计值
	q       float64 // q 可变荷载设计值
	sumLoad float64 // sumLoad 荷载总设计值
)

// CalBridge calculate the bridge
func CalBridge() {
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
