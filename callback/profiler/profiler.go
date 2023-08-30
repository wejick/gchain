package profiler

import (
	"context"

	"github.com/wejick/gchain/callback"
)

type Profiler struct {
	Events map[string][]*callback.CallbackData // map of callbackData by sessionID
}

func NewProfiler() *Profiler {
	return &Profiler{
		Events: make(map[string][]*callback.CallbackData),
	}
}

func (P *Profiler) Callback(context context.Context, data callback.CallbackData) {

}
