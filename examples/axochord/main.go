//-----------------------------------------------------------------------------
/*

OmniKeys 3 x 12 chord keys

*/
//-----------------------------------------------------------------------------

package main

import (
	"math"
	"strings"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

var b_r0 = 13.0 * 0.5    // major radius
var b_r1 = 7.0 * 0.5     // minor radius
var b_h0 = 6.0           // cavity height for button body
var b_h1 = 1.5           // thru panel thickness
var b_dv = 22.0          // vertical inter-button distance
var b_dh = 20.0          // horizontal inter-button distance
var b_theta = DtoR(20.0) // button angle

const BUTTONS_V = 3 // number of vertical buttons
const BUTTONS_H = 3 //12 // number of horizontal buttons

func button_cavity() SDF3 {
	p := NewPolygon()
	p.Add(0, -(b_h0 + b_h1))
	p.Add(b_r0, 0).Rel()
	p.Add(0, b_h0).Rel()
	p.Add(b_r1-b_r0, 0).Rel()
	p.Add(0, b_h1).Rel()
	p.Add(b_r0-b_r1, 0).Rel()
	p.Add(0, b_h0).Rel()
	p.Add(b_r1-b_r0, 0).Rel()
	p.Add(0, b_h1).Rel()
	p.Add(-b_r1, 0).Rel()
	return Revolve3D(Polygon2D(p.Vertices()))
}

// return the button matrix
func buttons() SDF3 {
	// single key column
	d := BUTTONS_V * b_dv
	p := V3{-math.Sin(b_theta) * d, math.Cos(b_theta) * d, 0}
	col := LineOf3D(button_cavity(), V3{}, p, strings.Repeat("x", BUTTONS_V))
	// multiple key columns
	d = BUTTONS_H * b_dh
	p = V3{d, 0, 0}
	matrix := LineOf3D(col, V3{}, p, strings.Repeat("x", BUTTONS_H))
	// centered on the origin
	d = (BUTTONS_V - 1) * b_dv
	dx := 0.5 * (((BUTTONS_H - 1) * b_dh) - (d * math.Sin(b_theta)))
	dy := 0.5 * d * math.Cos(b_theta)
	return Transform3D(matrix, Translate3d(V3{-dx, -dy, 0}))
}

//-----------------------------------------------------------------------------

func panel() SDF3 {
	v := (BUTTONS_V - 1) * b_dv
	vx := v * math.Sin(b_theta)
	vy := v * math.Cos(b_theta)

	sx := ((BUTTONS_H-1)*b_dh + vx) * 1.5
	sy := vy * 1.9

	pp := &PanelParms{
		Size:         V2{sx, sy},
		CornerRadius: 5.0,
		HoleDiameter: 3.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	// extrude to 3d
	return Extrude3D(Panel2D(pp), 2.0*(b_h0+b_h1))
}

//-----------------------------------------------------------------------------

func main() {
	s := Difference3D(panel(), buttons())
	upper := Cut3D(s, V3{}, V3{0, 0, 1})
	lower := Cut3D(s, V3{}, V3{0, 0, -1})

	RenderSTL(upper, 400, "upper.stl")
	RenderSTL(lower, 400, "lower.stl")
}

//-----------------------------------------------------------------------------
