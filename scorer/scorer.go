package scorer

import (
	"errors"

	"github.com/cszczepaniak/yahtzee/hand"
)

func Score(h hand.Hand, strat ScoringStrategy) (int, error) {
	defer strat.Clear()
	if len(h) != 5 {
		return 0, errors.New(`hand must have length 5`)
	}
	for _, n := range h {
		strat.AtEach(n)
	}
	return strat.Accumulate(), nil
}

type ScoringStrategy interface {
	AtEach(int)
	Accumulate() int
	Clear()
}

type singleDieScorer struct {
	target int
	sum    int
}

func NewSingleDieScorer(target int) ScoringStrategy {
	return &singleDieScorer{
		target: target,
	}
}

func (s *singleDieScorer) AtEach(n int) {
	if n == s.target {
		s.sum += s.target
	}
}

func (s *singleDieScorer) Accumulate() int { return s.sum }

func (s *singleDieScorer) Clear() { s.sum = 0 }

type nOfAKindScorer struct {
	n      int
	counts map[int]int
	sum    int
}

func NewNOfAKindScorer(n int) ScoringStrategy {
	return &nOfAKindScorer{
		n:      n,
		counts: make(map[int]int, 5),
	}
}

func (s *nOfAKindScorer) AtEach(n int) {
	s.counts[n]++
	s.sum += n
}

func (s *nOfAKindScorer) Accumulate() int {
	for _, ct := range s.counts {
		if ct >= s.n {
			return s.sum
		}
	}
	return 0
}

func (s *nOfAKindScorer) Clear() {
	for k := range s.counts {
		delete(s.counts, k)
	}
	s.sum = 0
}

type StraightKind int

const (
	SmallStraight StraightKind = 1
	LargeStraight StraightKind = 2
)

type straightScorer struct {
	kind   StraightKind
	unique map[int]struct{}
}

func NewStraightScorer(k StraightKind) ScoringStrategy {
	return &straightScorer{
		kind:   k,
		unique: make(map[int]struct{}, 5),
	}
}

func (s *straightScorer) AtEach(n int) { s.unique[n] = struct{}{} }

func (s *straightScorer) Accumulate() int {
	if len(s.unique) >= 4 && s.kind == SmallStraight {
		return 30
	}
	if len(s.unique) == 5 && s.kind == LargeStraight {
		return 40
	}
	return 0
}

func (s *straightScorer) Clear() {
	for k := range s.unique {
		delete(s.unique, k)
	}
}

type fullHouseScorer struct {
	counts map[int]int
}

func NewFullHouseScorer() ScoringStrategy {
	return &fullHouseScorer{
		counts: make(map[int]int, 5),
	}
}

func (s *fullHouseScorer) AtEach(n int) { s.counts[n]++ }

func (s *fullHouseScorer) Accumulate() int {
	if len(s.counts) != 2 {
		return 0
	}
	for _, ct := range s.counts {
		if ct == 1 || ct == 4 {
			break
		}
		return 25
	}
	return 0
}

func (s *fullHouseScorer) Clear() {
	for k := range s.counts {
		delete(s.counts, k)
	}
}

type chanceScorer struct {
	sum int
}

func NewChanceScorer() ScoringStrategy { return &chanceScorer{} }

func (s *chanceScorer) AtEach(n int)    { s.sum += n }
func (s *chanceScorer) Accumulate() int { return s.sum }
func (s *chanceScorer) Clear()          { s.sum = 0 }

type YahtzeeScorer struct {
	last      int
	isYahtzee bool
}

func NewYahtzeeScorer() ScoringStrategy { return &YahtzeeScorer{isYahtzee: true} }

func (s *YahtzeeScorer) AtEach(n int) {
	if s.last != 0 && s.last != n {
		s.isYahtzee = false
	}
	s.last = n
}

func (s *YahtzeeScorer) Accumulate() int {
	if s.isYahtzee {
		return 50
	}
	return 0
}

func (s *YahtzeeScorer) Clear() {
	s.last = 0
	s.isYahtzee = true
}
