package component

import "github.com/yohamta/donburi"

type AudioQueueData struct {
	queue [][]byte
}

func (q *AudioQueueData) Enqueue(clip []byte) {
	q.queue = append(q.queue, clip)
}

func (q *AudioQueueData) Dequeue() []byte {
	if q.Empty() {
		return nil
	}

	c := q.queue[0]

	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		q.Reset()
	}

	return c
}

func (q *AudioQueueData) Empty() bool {
	return len(q.queue) == 0
}

func (q *AudioQueueData) Reset() {
	q.queue = q.queue[:0]
}

var AudioQueue = donburi.NewComponentType[AudioQueueData]()
