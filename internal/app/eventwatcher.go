package app

import (
	"context"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// EventWatcher emits app events whenever events are triggered
type EventWatcher struct {
	ses            models.StreamEventSource
	configProvider *config.Provider
	authHandler    *handlers.Auth

	ctx    context.Context
	logger *zap.Logger

	stop     chan struct{}
	stopDone chan struct{}
}

// NewEventWatcher returns a new EventWatcher
func NewEventWatcher(
	streamEventSource models.StreamEventSource,
	configProvider *config.Provider,
	authHandler *handlers.Auth,
	ctx context.Context,
	logger *zap.Logger,
) *EventWatcher {
	return &EventWatcher{
		ses:            streamEventSource,
		configProvider: configProvider,
		authHandler:    authHandler,
		ctx:            ctx,
		logger:         logger.Named("app-event-watcher"),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

func EventsEmit(ctx context.Context, eventName string, optionalData ...interface{}) {
	if !IsHeadless(ctx) {
		runtime.EventsEmit(ctx, eventName, optionalData...)
	}
}

// Serve runs the service for the app event watcher
func (s *EventWatcher) Serve() {
	defer close(s.stopDone)
	streamCh, streamChID := s.ses.Subscribe()
	cfgUpdatesCh, cfgUpdatesChID := s.configProvider.UpdateEvents.Subscribe()
	cfgErrorsCh, cfgErrorsChID := s.configProvider.ErrorEvents.Subscribe()
	s.logger.Info("Running")

	for {
		select {
		case event := <-streamCh:
			_, isAddStream := event.Type.(models.AddStream)
			_, isRemoveStream := event.Type.(models.RemoveStream)
			_, isUpdateIDs := event.Type.(models.UpdateIDs)
			if isAddStream || isRemoveStream || isUpdateIDs {
				EventsEmit(s.ctx, "StreamChange")
			}
		case <-cfgUpdatesCh:
			s.authHandler.RefreshConfig()
			EventsEmit(s.ctx, "ConfigChange")
		case msg := <-cfgErrorsCh:
			EventsEmit(s.ctx, "ErrorEvent", msg)
		case <-s.stop:
			s.logger.Info("Stopping...")
			s.ses.Unsubscribe(streamChID)
			s.configProvider.UpdateEvents.Unsubscribe(cfgUpdatesChID)
			s.configProvider.ErrorEvents.Unsubscribe(cfgErrorsChID)
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning.
func (s *EventWatcher) Stop() {
	close(s.stop)
	<-s.stopDone
}
