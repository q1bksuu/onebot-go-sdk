package server

import (
	"context"

	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// NewHTTPEventServer creates an HTTP server for event receiving only.
// It registers generated event handlers and sets EventPath to "/event" when empty.
// Action requests are rejected with 404 to avoid nil handler panics.
func NewHTTPEventServer(cfg HTTPConfig, svc OneBotEventService) *HTTPServer {
	if cfg.EventPath == "" {
		cfg.EventPath = "/event"
	}

	ed := NewEventDispatcher()
	RegisterGeneratedEvents(ed, svc)

	return NewHTTPServer(
		WithHTTPConfig(cfg),
		WithEventHandler(ed),
		WithActionHandler(dispatcher.ActionRequestHandlerFunc(
			func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
				return nil, dispatcher.ErrActionNotFound
			},
		)),
	)
}
