package abstractions

import "context"

type EventDispatcher interface {
	DispatchEvents(ctx context.Context) error
}
