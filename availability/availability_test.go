package availability

import (
	"context"
	"testing"
)

func TestHolder_EventListeners(t *testing.T) {

}

func TestHolder_OnAvailabilityChangeEvent(t *testing.T) {
	holder := NewHolder()
	changeEvent := NewChangeEvent(context.Background(), StateCorrect)

	holder.OnAvailabilityChangeEvent(context.Background(), changeEvent)

}
