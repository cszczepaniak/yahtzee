package hand

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
)

type Hand []int

func New() Hand {
	return make([]int, 5)
}

func (h Hand) String() string {
	sb := strings.Builder{}
	for i, n := range h {
		sb.WriteString(strconv.Itoa(n))
		if i != len(h)-1 {
			sb.WriteString(`,`)
		}
	}
	return sb.String()
}

func (h Hand) RollAll() {
	for i := range h {
		// ignore the error because we know we're in the range of h
		_ = h.rollOne(i)
	}
}

func (h Hand) Roll(idxs []int) error {
	for _, i := range idxs {
		if err := h.rollOne(i); err != nil {
			return err
		}
	}
	return nil
}

func (h Hand) rollOne(i int) error {
	if i < 0 || i > len(h) {
		return errors.New(`index outside bounds of hand`)
	}
	h[i] = rand.Intn(6) + 1
	return nil
}
