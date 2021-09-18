package scorer_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/yahtzee/hand"
	"github.com/cszczepaniak/yahtzee/scorer"
)

func TestBadHands(t *testing.T) {
	tests := []struct {
		desc   string
		hand   hand.Hand
		expErr error
	}{{
		desc:   `too many in hand`,
		hand:   []int{1, 2, 3, 4, 5, 6},
		expErr: scorer.ErrInvalidHand,
	}, {
		desc:   `too few in hand`,
		hand:   []int{1, 2, 3},
		expErr: scorer.ErrInvalidHand,
	}, {
		desc:   `nil hand`,
		hand:   nil,
		expErr: scorer.ErrInvalidHand,
	}, {
		desc:   `invalid die`,
		hand:   []int{1, 2, 3, 4, 10},
		expErr: scorer.ErrInvalidDie,
	}, {
		desc:   `valid`,
		hand:   []int{1, 2, 3, 4, 5},
		expErr: nil,
	}}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			_, err := scorer.Score(tc.hand, scorer.NewYahtzeeScorer())
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestSingleDieScorer(t *testing.T) {
	for dieVal := 1; dieVal <= 6; dieVal++ {
		for nOfDie := 0; nOfDie <= 5; nOfDie++ {
			shuffleAndAssertScore(
				t, generateNOfK(nOfDie, dieVal),
				scorer.NewSingleDieScorer(dieVal),
				dieVal*nOfDie,
			)
		}
	}
}
func TestNOfAKindScorer(t *testing.T) {
	tests := []struct {
		n        int
		hand     hand.Hand
		expScore int
	}{{
		n:        3,
		hand:     []int{1, 2, 3, 4, 5},
		expScore: 0,
	}, {
		n:        3,
		hand:     []int{1, 2, 3, 4, 1},
		expScore: 0,
	}, {
		n:        3,
		hand:     []int{1, 1, 1, 3, 4},
		expScore: 10,
	}, {
		n:        3,
		hand:     []int{5, 5, 5, 5, 4},
		expScore: 24,
	}, {
		n:        3,
		hand:     []int{6, 6, 6, 6, 6},
		expScore: 30,
	}, {
		n:        4,
		hand:     []int{1, 1, 1, 3, 4},
		expScore: 0,
	}, {
		n:        4,
		hand:     []int{5, 5, 5, 5, 4},
		expScore: 24,
	}, {
		n:        4,
		hand:     []int{6, 6, 6, 6, 6},
		expScore: 30,
	}}
	for _, tc := range tests {
		shuffleAndAssertScore(t, tc.hand, scorer.NewNOfAKindScorer(tc.n), tc.expScore)
	}
}

func TestStraightScorer(t *testing.T) {
	tests := []struct {
		straightKind scorer.StraightKind
		hand         hand.Hand
		expScore     int
	}{{
		straightKind: scorer.SmallStraight,
		hand:         []int{1, 2, 3, 4, 5},
		expScore:     30,
	}, {
		straightKind: scorer.SmallStraight,
		hand:         []int{2, 3, 4, 5, 6},
		expScore:     30,
	}, {
		straightKind: scorer.SmallStraight,
		hand:         []int{1, 2, 3, 4, 1},
		expScore:     30,
	}, {
		straightKind: scorer.SmallStraight,
		hand:         []int{2, 3, 4, 5, 5},
		expScore:     30,
	}, {
		straightKind: scorer.SmallStraight,
		hand:         []int{1, 3, 4, 5, 6},
		expScore:     30,
	}, {
		straightKind: scorer.SmallStraight,
		hand:         []int{1, 1, 1, 3, 4},
		expScore:     0,
	}, {
		straightKind: scorer.LargeStraight,
		hand:         []int{1, 2, 3, 4, 5},
		expScore:     40,
	}, {
		straightKind: scorer.LargeStraight,
		hand:         []int{2, 3, 4, 5, 6},
		expScore:     40,
	}, {
		straightKind: scorer.LargeStraight,
		hand:         []int{1, 2, 3, 4, 4},
		expScore:     0,
	}, {
		straightKind: scorer.LargeStraight,
		hand:         []int{1, 1, 1, 3, 4},
		expScore:     0,
	}}
	for _, tc := range tests {
		shuffleAndAssertScore(t, tc.hand, scorer.NewStraightScorer(tc.straightKind), tc.expScore)
	}
}

func TestFullHouseScorer(t *testing.T) {
	tests := []struct {
		hand     []int
		expScore int
	}{{
		hand:     []int{1, 2, 3, 4, 5},
		expScore: 0,
	}, {
		hand:     []int{1, 1, 1, 1, 2},
		expScore: 0,
	}, {
		hand:     []int{1, 1, 2, 2, 3},
		expScore: 0,
	}, {
		hand:     []int{1, 1, 1, 2, 2},
		expScore: 25,
	}, {
		hand:     []int{1, 1, 2, 2, 2},
		expScore: 25,
	}, {
		hand:     []int{5, 5, 6, 6, 6},
		expScore: 25,
	}}
	for _, tc := range tests {
		shuffleAndAssertScore(t, tc.hand, scorer.NewFullHouseScorer(), tc.expScore)
	}
}

func TestChanceScorer(t *testing.T) {
	tests := []struct {
		hand     []int
		expScore int
	}{{
		hand:     []int{1, 1, 1, 1, 1},
		expScore: 5,
	}, {
		hand:     []int{1, 2, 3, 4, 5},
		expScore: 15,
	}, {
		hand:     []int{1, 4, 5, 5, 6},
		expScore: 21,
	}, {
		hand:     []int{2, 2, 4, 4, 5},
		expScore: 17,
	}, {
		hand:     []int{6, 6, 6, 6, 6},
		expScore: 30,
	}}
	for _, tc := range tests {
		shuffleAndAssertScore(t, tc.hand, scorer.NewChanceScorer(), tc.expScore)
	}
}

func TestYahtzeeScorer(t *testing.T) {
	for dieVal := 1; dieVal <= 6; dieVal++ {
		for nOfDie := 1; nOfDie <= 5; nOfDie++ {
			hand := generateNOfK(nOfDie, dieVal)
			exp := 0
			if nOfDie == 5 {
				exp = 50
			}
			shuffleAndAssertScore(t, hand, scorer.NewYahtzeeScorer(), exp)
		}
	}
}

func generateNOfK(n, k int) hand.Hand {
	res := make([]int, 5)
	fillWith := k%6 + 1
	for i := 0; i < n; i++ {
		res[i] = k
	}
	for i := n; i < len(res); i++ {
		res[i] = fillWith
	}
	return res
}

func shuffleAndAssertScore(t *testing.T, hand hand.Hand, strat scorer.ScoringStrategy, expScore int) {
	// shuffle and test 5 times to make sure order doesn't matter
	for i := 0; i < 5; i++ {
		rand.Shuffle(len(hand), func(i, j int) {
			hand[i], hand[j] = hand[j], hand[i]
		})
		score, err := scorer.Score(hand, strat)
		require.NoError(t, err)
		require.Equal(t, expScore, score, hand)
	}
}
