package audio

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var (
	ErrInvalidVolume = errors.New("volume must be between 0 to 1")
)

var (
	globalIsMuted          bool
	registeredAudioPlayers map[*audio.Player]float64
	mux                    *sync.RWMutex
)

func init() {
	registeredAudioPlayers = make(map[*audio.Player]float64)
	mux = &sync.RWMutex{}

	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			<-t.C

			mux.Lock()
			refreshedAudioPlayers := make(map[*audio.Player]float64)
			for player, lastVolume := range refreshedAudioPlayers {
				if player != nil {
					refreshedAudioPlayers[player] = lastVolume
				}
			}

			registeredAudioPlayers = refreshedAudioPlayers
			mux.Unlock()
		}
	}()
}

func GlobalMute(muted bool) {
	mux.Lock()
	defer mux.Unlock()
	if globalIsMuted == muted {
		return
	}

	globalIsMuted = muted

	for player, lastVolume := range registeredAudioPlayers {
		if player == nil {
			continue
		}

		if globalIsMuted {
			registeredAudioPlayers[player] = player.Volume()
			player.SetVolume(0.0)
		} else {
			player.SetVolume(lastVolume)
		}
	}
}

func IsGloballyMuted() bool {
	mux.RLock()
	defer mux.RUnlock()
	return globalIsMuted
}

type config struct {
	volume *float64
}

type AudioOption func(c *config)

func WithVolume(volume float64) AudioOption {
	return func(c *config) {
		c.volume = &volume
	}
}

// PlayOneShotWav will play the given WAV sound effect stream one time. This is meant for short sound effects only.
func PlayOneShotWav(sfx []byte, opts ...AudioOption) error {
	var c config
	for _, opt := range opts {
		opt(&c)
	}

	m, err := wav.DecodeF32(bytes.NewReader(sfx))
	if err != nil {
		return err
	}

	player, err := audio.CurrentContext().NewPlayerF32(m)
	if err != nil {
		return err
	}

	err = applyConfigAndRegister(c, player)
	if err != nil {
		return err
	}

	player.Play()

	return nil
}

func applyConfigAndRegister(c config, player *audio.Player) error {
	if c.volume != nil {
		if *c.volume < 0 || *c.volume > 1 {
			return ErrInvalidVolume
		}

		player.SetVolume(*c.volume)
	}

	mux.Lock()
	defer mux.Unlock()
	registeredAudioPlayers[player] = player.Volume()
	if globalIsMuted {
		player.SetVolume(0.0)
	}

	return nil
}
