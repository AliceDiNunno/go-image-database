package events

import (
	"github.com/AliceDiNunno/go-image-database/src/core/usecases"
	"time"
)

type event struct {
	usecases usecases.Usecases
}

func (e event) Start() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				//Jean Castex Bien Sur
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func NewEventManager(usecases usecases.Usecases) event {
	return event{
		usecases: usecases,
	}
}
