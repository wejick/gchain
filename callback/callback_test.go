package callback

import (
	"context"
	"fmt"
	"testing"
)

func TestManager_RegisterCallback(t *testing.T) {
	// Create a new Manager
	manager := NewManager()

	// Register some callbacks
	manager.RegisterCallback("event1", callback1)
	manager.RegisterCallback("event2", callback2)
	manager.RegisterCallback("event1", callback3)

	// Register duplicate callback for event1
	manager.RegisterCallback("event1", callback1)

	// Inspect the registered callbacks
	registeredCallbacks := manager.inspect()

	// Assert the number of registered callbacks
	if len(registeredCallbacks) != 2 {
		t.Errorf("Expected 2 registered events, got %d", len(registeredCallbacks))
	}

	// Assert the registered callbacks for a specific event
	// this will check if there's duplicate or not
	event1Callbacks := registeredCallbacks["event1"]
	if len(event1Callbacks) != 2 {
		t.Errorf("Expected 2 registered callbacks for event1, got %d", len(event1Callbacks))
	}

	event2Callbacks := registeredCallbacks["event2"]
	if len(event2Callbacks) != 1 {
		t.Errorf("Expected 1 registered callback for event2, got %d", len(event2Callbacks))
	}
}

func TestManager_TriggerEvent(t *testing.T) {
	// Create a new Manager
	manager := NewManager()

	// Variables to track callback invocations
	var callback1Invoked, callback2Invoked bool

	// Register some callbacks
	manager.RegisterCallback("event1", func(ctx context.Context, data CallbackData) {
		callback1Invoked = true
	})

	manager.RegisterCallback("event2", func(ctx context.Context, data CallbackData) {
		callback2Invoked = true
	})

	// Prepare test data
	ctx := context.Background()
	data := CallbackData{
		RunID:        "123",
		EventName:    "event1",
		FunctionName: "test",
		Input:        map[string]string{"input1": "value1"},
		Output:       map[string]string{"output1": "value1"},
		Data:         "some data",
	}

	// Reset invocation tracking variables
	callback1Invoked = false
	callback2Invoked = false

	// Trigger the event
	manager.TriggerEvent(ctx, "event1", data)
	manager.TriggerEvent(ctx, "event2", data)

	// Assert callback invocations
	if !callback1Invoked {
		t.Errorf("Callback 1 was not invoked")
	}

	if !callback2Invoked {
		t.Errorf("Callback 2 was not invoked")
	}
}

// Example callback functions
func callback1(ctx context.Context, data CallbackData) {
	fmt.Println("Callback 1 called")
	// Add assertions or verification specific to callback1
}

func callback2(ctx context.Context, data CallbackData) {
	fmt.Println("Callback 2 called")
	// Add assertions or verification specific to callback2
}

func callback3(ctx context.Context, data CallbackData) {
	fmt.Println("Callback 3 called")
	// Add assertions or verification specific to callback3
}
