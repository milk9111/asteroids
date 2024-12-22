package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var errNoScoreQueueFound = newComponentError("no score queue found")

type ScoreQueueData struct {
	queue []int
}

func (q *ScoreQueueData) Enqueue(amount int) {
	q.queue = append(q.queue, amount)
}

func (q *ScoreQueueData) Dequeue() int {
	if q.Empty() {
		return 0
	}

	a := q.queue[0]

	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		q.Reset()
	}

	return a
}

func (q *ScoreQueueData) Empty() bool {
	return len(q.queue) == 0
}

func (q *ScoreQueueData) Reset() {
	q.queue = q.queue[:0]
}

var ScoreQueue = donburi.NewComponentType[ScoreQueueData]()

func MustFindScoreQueue(w donburi.World) *ScoreQueueData {
	scoreQueue, ok := donburi.NewQuery(filter.Contains(ScoreQueue)).First(w)
	if !ok {
		panic(errNoScoreQueueFound)
	}
	return ScoreQueue.Get(scoreQueue)
}
