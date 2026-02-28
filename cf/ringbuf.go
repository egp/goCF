// ringbuf.go v1
package cf

import "strings"

// RingBuf keeps the last N strings (fingerprints) and can detect repeats.
type RingBuf struct {
	buf    []string
	next   int
	size   int
	filled bool
}

func NewRingBuf(n int) *RingBuf {
	if n <= 0 {
		n = 1
	}
	return &RingBuf{buf: make([]string, n)}
}

func (r *RingBuf) Cap() int { return len(r.buf) }

func (r *RingBuf) Len() int {
	if r.filled {
		return len(r.buf)
	}
	return r.size
}

func (r *RingBuf) Add(s string) {
	r.buf[r.next] = s
	r.next++
	if r.next >= len(r.buf) {
		r.next = 0
		r.filled = true
	}
	if !r.filled && r.size < len(r.buf) {
		r.size++
	}
}

// Count returns how many times s appears in the current window.
func (r *RingBuf) Count(s string) int {
	n := r.Len()
	if n == 0 {
		return 0
	}
	c := 0
	for i := 0; i < n; i++ {
		if r.buf[i] == s {
			c++
		}
	}
	return c
}

// Dump returns the window from oldest->newest, one per line, prefixed with index.
// Useful for embedding in errors.
func (r *RingBuf) Dump() string {
	n := r.Len()
	if n == 0 {
		return ""
	}
	var b strings.Builder

	// Determine oldest index.
	start := 0
	if r.filled {
		start = r.next
	}

	for i := 0; i < n; i++ {
		idx := (start + i) % len(r.buf)
		b.WriteString("#")
		b.WriteString(itoa(i))
		b.WriteString(" ")
		b.WriteString(r.buf[idx])
		if i+1 < n {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// small local int->string without fmt (keeps allocations small; correctness only here).
func itoa(x int) string {
	if x == 0 {
		return "0"
	}
	if x < 0 {
		x = -x
	}
	var tmp [32]byte
	i := len(tmp)
	for x > 0 {
		i--
		tmp[i] = byte('0' + (x % 10))
		x /= 10
	}
	return string(tmp[i:])
}

// ringbuf.go v1
