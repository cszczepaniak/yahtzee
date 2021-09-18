package hand_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/yahtzee/hand"
)

func TestRollAll(t *testing.T) {
	h := hand.New()
	h.RollAll()
	for _, n := range h {
		require.Greater(t, n, 0)
	}
}

func TestRollSome(t *testing.T) {
	h := hand.New()
	shouldRoll := []int{1, 2, 4}
	shouldRollMap := map[int]struct{}{1: {}, 2: {}, 4: {}}
	h.Roll(shouldRoll)
	for i := range h {
		if _, ok := shouldRollMap[i]; ok {
			require.Greater(t, h[i], 0)
			continue
		}
		require.Zero(t, h[i])
	}
}
