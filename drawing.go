package main

import (
	"github.com/nsf/termbox-go"
)

func fill(x, y, w, h int, c termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, c.Ch, c.Fg, c.Bg)
		}
	}
}

func draw_box(x0, y0 int) {
	const coldef = termbox.ColorDefault
	w := W * 2
	h := H
	termbox.SetCell(x0, y0+h, '\u255a', coldef, coldef)     //└
	termbox.SetCell(x0+w+1, y0+h, '\u255d', coldef, coldef) //┘
	fill(x0+1, y0+h, w, 1, termbox.Cell{Ch: '\u2550'})      //─
	fill(x0, y0, 1, h, termbox.Cell{Ch: '\u2551'})          //║
	fill(x0+w+1, y0, 1, h, termbox.Cell{Ch: '\u2551'})      //║
}

func (p field) draw (x0, y0 int) {
	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			fill(x0 + 2*i, y0 + j, 2, 1, termbox.Cell( p[i][j].texture ))
		}
	}
}

func (f figure) draw(x0, y0 int) {
	var x, y int
	for c := 0; c < 4; c++ {
		x = f.loc.x + f.offset[c].x
		y = f.loc.y + f.offset[c].y
		if x >= 0 && x <= W && y >= 0 && y <= H {
			fill(x0 + 2*x, y0 + y, 2, 1, termbox.Cell( f.texture ))
		}
	}
}

func draw_screen(p field, f figure) {
	const coldef = termbox.ColorDefault
	
	termbox.Clear(coldef, coldef)
	defer termbox.Flush()

	sz_x, sz_y := termbox.Size()
	x0 := (sz_x - 2*W - 2) / 2
	y0 := (sz_y - H - 1) / 2
	draw_box(x0, y0)
	p.draw(x0+1, y0)
	f.draw(x0+1, y0)
}

