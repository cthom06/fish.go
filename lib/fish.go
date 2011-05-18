package fish

import (
	"io"
	"rand"
	"fmt"
	"os"
)

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

type Error os.Error

var error = os.NewError("something smells fishy...")
var NoError = os.NewError("")

type codebox struct {
	D map[uint64]int
	MaxY int
	MaxX map[int]int
}

func pair(x, y int) uint64 {
	return uint64(x) | (uint64(y) << 32)
}

func newcodebox() *codebox {
	return &codebox{ make(map[uint64]int), 0, make(map[int]int) }
}

func (c *codebox) Set(x, y, v int) {
	c.D[pair(x, y)] = v
	if y > c.MaxY { c.MaxY = y }
	if t, ok := c.MaxX[y]; ok {
		if x > t { c.MaxX[y] = x }
	} else if x > 0 {
		c.MaxX[y] = x
	}
}

func (c *codebox) Get(x, y int) (int, Error) {
	if v, ok := c.D[pair(x, y)]; ok {
		return v, nil
	}
	return int(' '), error
}

type register struct {
	Pop bool
	Val int
}

type runtime struct {
	Stacks [][]int
	Box *codebox
	Dir byte
	Pos [2]int
	Register []register
}

func NewRuntime(code []byte, stacks ...[]int) (*runtime) {
	if len(stacks) == 0 {
		stacks = [][]int{[]int{}}
	}
	tmp := &runtime{
		stacks,
		newcodebox(),
		RIGHT,
		[2]int{0,0},
		[]register{register{false,0}},
	}
	x := 0
	y := 0
	for _, v := range code {
		if v == '\n' {
			y++
			x = 0
		} else {
			tmp.Box.Set(x, y, int(v))
			x++
		}
	}
	return tmp
}

func (r *runtime) Move() {
	switch r.Dir {
	case UP:
		if r.Pos[0] == 0 {
			r.Pos[0] = r.Box.MaxY
		} else r.Pos[0]--
	case DOWN:
		if r.Pos[0] == r.Box.MaxY {
			r.Pos[0] = 0
		} else r.Pos[0]++
	case LEFT:
		if r.Pos[1] == 0 {
			r.Pos[1],_ = r.Box.MaxX[r.Pos[0]]
		} else r.Pos[1]--
	case RIGHT:
		mx,_ := r.Box.MaxX[r.Pos[0]]
		if r.Pos[1] == mx {
			r.Pos[1] = 0
		} else r.Pos[1]++
	}
}

func (r *runtime) Get(x, y int) (int, Error) {
	return r.Box.Get(x, y)
}

func (r *runtime) Set(x, y, v int) {
	r.Box.Set(x, y, v)
}

func (r *runtime) Read() int {
	t,_ := r.Box.Get(r.Pos[1],r.Pos[0])
	return t
}

func (r *runtime) Push(v int) {
	r.Stacks[len(r.Stacks) - 1] = append(r.Stacks[len(r.Stacks) - 1], v)
}

func (r *runtime) Pop() (int, Error) {
	if len(r.Stacks[len(r.Stacks) - 1]) > 0 {
		v := r.Stacks[len(r.Stacks) - 1][len(r.Stacks[len(r.Stacks) - 1]) - 1]
		r.Stacks[len(r.Stacks) - 1] = r.Stacks[len(r.Stacks) - 1][0:len(r.Stacks[len(r.Stacks) - 1]) - 1]
		return v, nil
	}
	return 0, error
}

func (r *runtime) Len() int {
	return len(r.Stacks[len(r.Stacks) - 1])
}

func (r *runtime) Split(size int) Error {
	tmp := make([]int, size)
	for size--; size >= 0; size-- {
		e := Error(nil)
		if tmp[size], e = r.Pop(); e != nil {
			return error
		}
	}
	r.Stacks = append(r.Stacks, tmp)
	r.Register = append(r.Register, register{false,0})
	return nil
}

func (r *runtime) Merge() Error {
	if len(r.Stacks) <= 1 {
		return error
	}
	tmp := r.Stacks[len(r.Stacks) - 1]
	r.Stacks = r.Stacks[0:len(r.Stacks) - 1]
	for _, v := range tmp {
		r.Push(v)
	}
	r.Register = r.Register[:len(r.Register) - 1]
	return nil
}

func (r *runtime) RShift() {
	t := r.Stacks[len(r.Stacks) - 1]
	if len(t) > 0 {
		v := t[len(t) - 1]
		for i := len(t) - 1; i > 0; i-- {
			t[i] = t[i - 1]
		}
		t[0] = v
	}
}

func (r *runtime) LShift() {
	t := r.Stacks[len(r.Stacks) - 1]
	if len(t) > 0 {
		v := t[0]
		for i := 0; i < len(t) - 2; i++ {
			t[i] = t[i+1]
		}
		t[len(t) - 1] = v
	}
}

func (r *runtime) Reverse() {
	t := r.Stacks[len(r.Stacks) - 1]
	tmp := make([]int, len(t))
	j := 0
	for i := len(t) - 1; i >= 0; i-- {
		tmp[i] = t[j]
		j++
	}
	r.Stacks[len(r.Stacks) - 1] = tmp
}

func (r *runtime) Do(w byte, in io.Reader, out io.Writer) Error {
	switch w {
	case ' ':
	case '<':
		r.Dir = LEFT
	case '>':
		r.Dir = RIGHT
	case 'v':
		r.Dir = DOWN
	case '^':
		r.Dir = UP
	case '/':
		switch r.Dir {
		case UP:
			r.Dir = RIGHT
		case DOWN:
			r.Dir = LEFT
		case LEFT:
			r.Dir = DOWN
		case RIGHT:
			r.Dir = UP
		}
	case '\\':
		switch r.Dir {
		case UP:
			r.Dir = LEFT
		case DOWN:
			r.Dir = RIGHT
		case LEFT:
			r.Dir = UP
		case RIGHT:
			r.Dir = DOWN
		}
	case '|':
		switch r.Dir {
		case LEFT:
			r.Dir = RIGHT
		case RIGHT:
			r.Dir = LEFT
		}
	case '_':
		switch r.Dir {
		case UP:
			r.Dir = DOWN
		case DOWN:
			r.Dir = UP
		}
	case '#':
		switch r.Dir {
		case UP:
			r.Dir = DOWN
		case DOWN:
			r.Dir = UP
		case LEFT:
			r.Dir = RIGHT
		case RIGHT:
			r.Dir = LEFT
		}
	case 'x':
		switch rand.Int() % 4 {
		case 0:
			r.Dir = UP
		case 1:
			r.Dir = DOWN
		case 2:
			r.Dir = LEFT
		case 3:
			r.Dir = RIGHT
		}
	case '+':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		r.Push(v1 + v2)
	case '-':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		r.Push(v2 - v1)
	case '*':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		r.Push(v1 * v2)
	case ',':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		r.Push(v2 / v1)
	case '%':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		r.Push(v2 % v1)
	case '0','1','2','3','4','5','6','7','8','9':
		r.Push(int(w - '0'))
	case 'a','b','c','d','e','f':
		r.Push(int(w - 'a' + 10))
	case '=':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		if v1 == v2 {
			r.Push(1)
		} else r.Push(0)
	case ')':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		if v2 > v1 {
			r.Push(1)
		} else r.Push(0)
	case '(':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		if v2 < v1 {
			r.Push(1)
		} else r.Push(0)
	case '\'', '"':
		c := true
		for c {
			r.Move()
			m := r.Read()
			if int(w) == m {
				c = false
			} else r.Push(m)
		}
	case '!':
		r.Move()
	case '?':
		t,e := r.Pop()
		if e != nil {
			return error
		}
		if t == 0 {
			r.Move()
		}
	case 'l':
		r.Push(r.Len())
	case ':':
		if t, e := r.Pop(); e == nil {
			r.Push(t)
			r.Push(t)
		} else return error
	case '~':
		if _,e := r.Pop(); e != nil {
			return error
		}
	case '$':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		if e1 != nil || e1 != e2 {
			return error
		}
		r.Push(v1)
		r.Push(v2)
	case '@':
		v1, e1 := r.Pop()
		v2, e2 := r.Pop()
		v3, e3 := r.Pop()
		if e1 != nil || e2 != nil || e3 != nil {
			return error
		}
		r.Push(v1)
		r.Push(v3)
		r.Push(v2)
	case '&':
		if r.Register[len(r.Register) - 1].Pop {
			r.Register[len(r.Register) - 1].Pop = false
			r.Push(r.Register[len(r.Register) - 1].Val)
		} else if v, e := r.Pop(); e == nil {
			r.Register[len(r.Register) - 1].Pop = true
			r.Register[len(r.Register) - 1].Val = v
		} else return error
	case 'r':
		r.Reverse()
	case '{':
		r.LShift()
	case '}':
		r.RShift()
	case 'g':
		y, e1 := r.Pop()
		x, e2 := r.Pop()
		if e1 != nil || e2 != nil {
			return error
		}
		if v, e := r.Box.Get(x, y); e == nil {
			r.Push(v)
		} else return error
	case 'p':
		y, e1 := r.Pop()
		x, e2 := r.Pop()
		v, e3 := r.Pop()
		if e1 != nil || e2 != nil || e3 != nil {
			return error
		}
		r.Box.Set(x, y, v)
	case 'o':
		if v, e := r.Pop(); e == nil {
			if _, e = fmt.Fprintf(out, "%c", byte(v)); e != nil {
				return error
			}
		} else return error
	case 'n':
		if v, e := r.Pop(); e == nil {
			if _, e = fmt.Fprintf(out, "%d", v); e != nil {
				return error
			}
		} else return error
	case 'i':
		b := []byte{0}
		for n, e := in.Read(b); n < 1; n, e = in.Read(b) {
			if e != nil {
				r.Push(-1)
				return nil
			}
		}
		r.Push(int(b[0]))
	case '[':
		if v, e := r.Pop(); e == nil {
			if r.Split(v) != nil {
				return error
			}
		} else return error
	case ']':
		if r.Merge() != nil {
			return error
		}
	case ';':
		return NoError
	default:
		return error
	}
	return nil
}

func (r *runtime) Run(in io.Reader, out io.Writer) Error {
	for {
		m := r.Read()
		if m > 255 {
			return error
		} else {
			if e := r.Do(byte(m), in, out); e != nil {
				return e
			}
		}
		r.Move()
	}
	return nil
}

func (r *runtime) Debug(in io.Reader, out, debug io.Writer) Error {
	for {
		m := r.Read()
		if m > 255 {
			return error
		} else {
			fmt.Fprintf(debug, "%c\n", byte(m))
			if e := r.Do(byte(m), in, out); e != nil {
				return e
			}
		}
		r.Move()
	}
	return nil
}
