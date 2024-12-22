package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	lifeGainScoreInterval = 5000
	maxPlayerLives        = 5
)

var errNoScoreQueueFound = newSystemError("no score queue found")

type Score struct {
	query *donburi.Query
	game  *component.GameData

	nextLifeGainScore int
}

func NewScore() *Score {
	return &Score{
		query: donburi.NewQuery(
			filter.Contains(component.ScoreQueue),
		),
		nextLifeGainScore: lifeGainScoreInterval,
	}
}

func (s *Score) Update(w donburi.World) {
	if s.game == nil {
		s.game = component.MustFindGame(w)
	}

	e, ok := s.query.First(w)
	if !ok {
		panic(errNoScoreQueueFound)
	}

	scoreQueue := component.ScoreQueue.Get(e)
	for !scoreQueue.Empty() {
		amount := scoreQueue.Dequeue()
		if amount == 0 {
			continue
		}

		s.game.Score += amount

		if s.game.Score >= s.nextLifeGainScore && s.game.PlayerLives < maxPlayerLives {
			s.game.PlayerLives++
			s.nextLifeGainScore += lifeGainScoreInterval
		}
	}

	scoreQueue.Reset()
}
