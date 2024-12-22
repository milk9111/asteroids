package component

import "github.com/yohamta/donburi"

type collision struct {
	collidedWith *donburi.Entry
	layer        CollisionLayer
}

type CollisionQueueData struct {
	queue []collision
}

func (q *CollisionQueueData) Enqueue(collidedWith *donburi.Entry, layer CollisionLayer) {
	q.queue = append(q.queue, collision{
		collidedWith: collidedWith,
		layer:        layer,
	})
}

func (q *CollisionQueueData) Dequeue() (*donburi.Entry, CollisionLayer) {
	if q.Empty() {
		return nil, CollisionLayerUnknown
	}

	c := q.queue[0]

	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		q.Reset()
	}

	return c.collidedWith, c.layer
}

func (q *CollisionQueueData) Empty() bool {
	return len(q.queue) == 0
}

func (q *CollisionQueueData) Reset() {
	q.queue = q.queue[:0]
}

var CollisionQueue = donburi.NewComponentType[CollisionQueueData]()
