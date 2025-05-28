package utils_test

import (
	"testing"

	"github.com/HarikrishnanBalagopal/fluidsim/utils"
)

func TestStep(t *testing.T) {
	t.Run("try multiple steps", func(t *testing.T) {
		var time float32 = 0.1
		const dt float32 = 0.1
		for time < 10.0 {
			utils.Step(time, dt)
			time += dt
		}
	})
}

func BenchmarkStep(t *testing.B) {
	t.Run("run multiple steps", func(t *testing.B) {
		var time float32 = 0.1
		const dt float32 = 0.1
		for time < 10.0 {
			utils.Step(time, dt)
			time += .1
		}
	})
}
