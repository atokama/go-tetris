package main

/* Simple console tetris game. 
	Uses go-termbox for drawing in terminal
	
	Controls:
		left		- h
		right 	- l
		down		- j
		up			- k (for test purposes)
		full down- space
		rotate 	- r 
		exit		- ctrl-c
*/


import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const (
	W        = 10		//width of gamefield
	H        = 20		//height
	nFigures = 7		//number of different figures
)

type axis struct{ x, y int }
type fig_shape uint8
type fig_angle uint8
type cell_texture termbox.Cell

type field [W][H]struct {
	texture cell_texture
	isFree bool
}

type figure struct {
	shape   fig_shape
	angle   fig_angle
	texture cell_texture
	loc     axis
	offset  [4]axis
}

const (			//list of kinds of figure shapes
	O_shaped fig_shape = iota
	Z_shaped
	S_shaped
	I_shaped
	T_shaped
	L_shaped
	J_shaped
)

const (		//list of possible figure rotations (analog clock)
	clock_12 fig_angle = iota
	clock_3
	clock_6
	clock_9
)

// map kind of figure to possible rotation angles
var shapeToAngles = map[fig_shape][]fig_angle{
	O_shaped: {clock_12},		
	Z_shaped: {clock_12, clock_3},
	S_shaped: {clock_12, clock_3},
	I_shaped: {clock_12, clock_3},
	T_shaped: {clock_12, clock_3, clock_6, clock_9},
	L_shaped: {clock_12, clock_3, clock_6, clock_9},
	J_shaped: {clock_12, clock_3, clock_6, clock_9},
}

// texture of figures (utf-8 characters)
var shapeToTexture = map[fig_shape]cell_texture{
	O_shaped: {Ch: '▯'},  //█▯
	Z_shaped: {Ch: '▯'},
	S_shaped: {Ch: '▯'},
	I_shaped: {Ch: '▯'},
	T_shaped: {Ch: '▯'},
	L_shaped: {Ch: '▯'},
	J_shaped: {Ch: '▯'},
}

//relative coordinates of figure for every possible 
//rotation angle
var figOffset = map[string][4]axis{
	"shape_O_clock_12": {{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	"shape_Z_clock_12": {{0, 0}, {0, 1}, {1, 0}, {1, -1}},
	"shape_Z_clock_3":  {{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	"shape_S_clock_12": {{0, -1}, {0, 0}, {1, 0}, {1, 1}},
	"shape_S_clock_3":  {{0, 1}, {1, 1}, {1, 0}, {2, 0}},
	"shape_I_clock_12": {{0, -1}, {0, 0}, {0, 1}, {0, 2}},
	"shape_I_clock_3":  {{-1, 1}, {0, 1}, {1, 1}, {2, 1}},
	"shape_T_clock_12": {{0, 0}, {-1, 1}, {0, 1}, {1, 1}},
	"shape_T_clock_3":  {{0, 0}, {0, 1}, {0, 2}, {1, 1}},
	"shape_T_clock_6":  {{-1, 1}, {0, 1}, {1, 1}, {0, 2}},
	"shape_T_clock_9":  {{-1, 1}, {0, 1}, {0, 0}, {0, 2}},
	"shape_L_clock_12": {{0, -1}, {0, 0}, {0, 1}, {1, 1}},
	"shape_L_clock_3":  {{0, 1}, {0, 0}, {1, 0}, {2, 0}},
	"shape_L_clock_6":  {{0, -1}, {1, -1}, {1, 0}, {1, 1}},
	"shape_L_clock_9":  {{1, 0}, {1, 1}, {0, 1}, {-1, 1}},
	"shape_J_clock_12": {{0, 1}, {1, 1}, {1, 0}, {1, -1}},
	"shape_J_clock_3":  {{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	"shape_J_clock_6":  {{0, 1}, {0, 0}, {0, -1}, {1, -1}},
	"shape_J_clock_9":  {{-1, 0}, {0, 0}, {1, 1}, {1, 0}},
}

//set relative coordinates 
func (f *figure) setOffset() {
	s := f.shape
	a := f.angle
	switch {
	case s == O_shaped && a == clock_12:
		f.offset = figOffset["shape_O_clock_12"]
	case s == Z_shaped && a == clock_12:
		f.offset = figOffset["shape_Z_clock_12"]
	case s == Z_shaped && a == clock_3:
		f.offset = figOffset["shape_Z_clock_3"]
	case s == S_shaped && a == clock_12:
		f.offset = figOffset["shape_S_clock_12"]
	case s == S_shaped && a == clock_3:
		f.offset = figOffset["shape_S_clock_3"]
	case s == I_shaped && a == clock_12:
		f.offset = figOffset["shape_I_clock_12"]
	case s == I_shaped && a == clock_3:
		f.offset = figOffset["shape_I_clock_3"]
	case s == T_shaped && a == clock_12:
		f.offset = figOffset["shape_T_clock_12"]
	case s == T_shaped && a == clock_3:
		f.offset = figOffset["shape_T_clock_3"]
	case s == T_shaped && a == clock_6:
		f.offset = figOffset["shape_T_clock_6"]
	case s == T_shaped && a == clock_9:
		f.offset = figOffset["shape_T_clock_9"]
	case s == L_shaped && a == clock_12:
		f.offset = figOffset["shape_L_clock_12"]
	case s == L_shaped && a == clock_3:
		f.offset = figOffset["shape_L_clock_3"]
	case s == L_shaped && a == clock_6:
		f.offset = figOffset["shape_L_clock_6"]
	case s == L_shaped && a == clock_9:
		f.offset = figOffset["shape_L_clock_9"]
	case s == J_shaped && a == clock_12:
		f.offset = figOffset["shape_J_clock_12"]
	case s == J_shaped && a == clock_3:
		f.offset = figOffset["shape_J_clock_3"]
	case s == J_shaped && a == clock_6:
		f.offset = figOffset["shape_J_clock_6"]
	case s == J_shaped && a == clock_9:
		f.offset = figOffset["shape_J_clock_9"]
	}
}

//make new game field
func newField() field {
	var p field
	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			p[i][j].texture = cell_texture{Ch: ' '}
			p[i][j].isFree = true
		}
	}
	return p
}

// put figure onto the gamefield
func (p *field) putFigure(f figure) {
	var w, h int
	for c := 0; c < 4; c++ {
		w = f.loc.x + f.offset[c].x
		h = f.loc.y + f.offset[c].y
		//p[w][h] = field{texture: f.texture, isFree: false}
		p[w][h].texture = f.texture
		p[w][h].isFree = false
	}
}

// clear full-filled horizontal lines
func (p *field) clearLines() {
	var isFullLine bool
	for h := H-1; h >= 0; h-- {
		isFullLine = true
		for w := 0; w < W; w++ {
			if p[w][h].isFree == true {
				isFullLine = false
				break
			}
		}
		if isFullLine {
			for ww := 0; ww < W; ww++ {
				for hh := h; hh > 0; hh-- {
					p[ww][hh] = p[ww][hh-1]
				}
			}
			for ww := 0; ww < W; ww++ {
				p[ww][0].isFree = true
			}
			p.clearLines() //check recursion
		}
	}
}

//make new figure with random shape, apply texture and place
//on top of gamefield
func newFigure() figure {
	var f figure
	f.shape = fig_shape(rand.Int31n(nFigures))
	f.angle = shapeToAngles[f.shape][0]
	f.texture = shapeToTexture[f.shape]
	f.loc = axis{4, 0}
	f.setOffset()
	return f
}

//check is figure not placed on already occupied square
func (f *figure) isFits(p field) bool {
	var x, y int
	for c := 0; c < 4; c++ {
		x = f.loc.x + f.offset[c].x
		y = f.loc.y + f.offset[c].y
		if !( x>=0 && x<W && y<H ) {
			return false
		} else if y>=0 && p[x][y].isFree == false {
			return false
		}
	}
	return true
}

func (f *figure) isFullyVisible() bool {
	for c := 0; c < 4; c++ {
		if f.loc.y + f.offset[c].y < 0 {
			return false
		}
	}
	return true
}

func (f *figure) rotate(p field) {
	//rotate through available angles
	n := len(shapeToAngles[f.shape])
	f.angle = fig_angle((int(f.angle) + 1) % n)
	f.setOffset()
	if !f.isFits(p) {
		f.angle = fig_angle( (int(f.angle) + n - 1) % n )
		f.setOffset()
	}
}

func (f *figure) moveLeft(p field) {
	if f.loc.x--; !f.isFits(p) { 
		f.loc.x++ 
	}
}

func (f *figure) moveRight(p field) { 
	if f.loc.x++; !f.isFits(p) { 
		f.loc.x-- 
	} 
} 

func (f *figure) moveDown(p field) bool { 
	if f.loc.y++; !f.isFits(p) {
		f.loc.y-- 
		return false //not possible to move down
	}
	return true 
}

func (f *figure) fullDown(p field) {
	ok := f.moveDown(p)
	for ok { 
		ok = f.moveDown(p)
	}
}

func (f *figure) moveUp(p field) { //for test
	if f.loc.y--; !f.isFits(p) {
		f.loc.y++
	} 
}

func main() {
	var ok bool
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.HideCursor()

	rand.Seed(time.Now().UnixNano())
	p := newField()
	f := newFigure()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			termbox.Interrupt()
		}
	}()

	draw_screen(p, f)
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				break mainloop
			case termbox.KeySpace:
				f.fullDown(p)
				//p.putFigure(f)
			default:
				switch ev.Ch {
				case 'h':
					f.moveLeft(p)
				case 'l':
					f.moveRight(p)
				case 'j':
					f.moveDown(p)
				case 'k':
					f.moveUp(p)
				case 'r':
					f.rotate(p)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			ok = f.moveDown(p)
			if !ok {
				if f.isFullyVisible() {
					p.putFigure(f)
					f = newFigure()
				} else {
					break mainloop //gameover
				}
			}
		}
		p.clearLines()
		draw_screen(p, f)
	}
}
