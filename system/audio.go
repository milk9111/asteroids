package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine/audio"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Audio struct {
	query *donburi.Query
}

func NewAudio() *Audio {
	return &Audio{
		query: donburi.NewQuery(
			filter.Contains(
				component.AudioQueue,
			),
		),
	}
}

func (a *Audio) Update(w donburi.World) {
	a.query.Each(w, func(e *donburi.Entry) {
		audioQueue := component.AudioQueue.Get(e)
		for !audioQueue.Empty() {
			clip := audioQueue.Dequeue()
			if len(clip) == 0 {
				continue
			}

			err := audio.PlayOneShotWav(clip, audio.WithVolume(0.5))
			if err != nil {
				panic(err)
			}
		}

		audioQueue.Reset()
	})
}
