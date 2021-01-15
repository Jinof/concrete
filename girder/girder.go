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
	l01     float64
	l02     float64
)

// CalGirder calculate the girder
func CalGirder() {
	fmt.Printf("gBridge 次梁传来的永久荷载 %f, qBridge 次梁传来的可变荷载 %f, selfG 次梁自重 %f \n", gGirder, qBridge, selfG)
	fmt.Printf("G 永久荷载设计值 %f, Q 可变荷载设计值 %f \n", G, Q)
}

func init() {
	CalLoad()
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
}
