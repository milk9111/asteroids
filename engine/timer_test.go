package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimer_OverridePercentDone(t *testing.T) {
	testCases := []struct {
		name                  string
		timer                 *Timer
		inputPercentDone      float64
		expectedCurrentFrames int
	}{
		{
			name: "with_percent_done_0_should_return_current_frame_0",
			timer: &Timer{
				currentFrames: 500,
				targetFrames:  1000,
			},
			inputPercentDone:      0.0,
			expectedCurrentFrames: 0,
		},
		{
			name: "with_percent_done_0.5_should_return_current_frame_500",
			timer: &Timer{
				currentFrames: 101,
				targetFrames:  1000,
			},
			inputPercentDone:      0.5,
			expectedCurrentFrames: 500,
		},
		{
			name: "with_percent_done_1.0_should_return_current_frame_1000",
			timer: &Timer{
				currentFrames: 101,
				targetFrames:  1000,
			},
			inputPercentDone:      1.0,
			expectedCurrentFrames: 1000,
		},
		{
			name: "with_percent_done_0.7487_should_return_current_frame_748",
			timer: &Timer{
				currentFrames: 101,
				targetFrames:  1000,
			},
			inputPercentDone:      0.7487,
			expectedCurrentFrames: 748,
		},
		{
			name: "with_percent_done_0.7487_total_frames_509_should_return_current_frame_381",
			timer: &Timer{
				currentFrames: 101,
				targetFrames:  509,
			},
			inputPercentDone:      0.7487,
			expectedCurrentFrames: 381,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.timer.OverridePercentDone(testCase.inputPercentDone)
			require.Equal(t, testCase.expectedCurrentFrames, testCase.timer.currentFrames)
		})
	}
}
