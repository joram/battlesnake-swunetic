package snek

type Stack []Point

func (s Stack) Push(v Point) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, Point) {
	l := len(s)
	if l == 0 {
		return s, Point{}
	}
	return s[:l-1], s[l-1]
}

func (s Stack) Len() int {
	return len(s)
}
