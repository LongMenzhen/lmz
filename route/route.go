package route

import (
	"github.com/cyrnicolase/lmz/app"
	"github.com/cyrnicolase/lmz/engine"
)

func init() {
	engine.Registe("somebody", app.SomebodyAction)
	engine.Registe("welcome", app.WelcomeAction)
}
