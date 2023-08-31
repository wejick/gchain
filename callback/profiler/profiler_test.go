package profiler

import (
	"context"
	"testing"

	"github.com/wejick/gchain/callback"
)

func TestProfiler(t *testing.T) {
	// create a new profiler
	p := NewProfiler()

	// create a context
	ctx := context.Background()

	// define the test cases
	testCases := []struct {
		name     string
		data     []callback.CallbackData
		session  string
		expected []callback.CallbackData
	}{
		{
			name: "single session",
			data: []callback.CallbackData{
				{SessionID: "session1", EventName: "hello"},
				{SessionID: "session1", EventName: "how are you?"},
			},
			session: "session1",
			expected: []callback.CallbackData{
				{SessionID: "session1", EventName: "hello"},
				{SessionID: "session1", EventName: "how are you?"},
			},
		},
		{
			name: "multiple sessions",
			data: []callback.CallbackData{
				{SessionID: "session1", EventName: "hello"},
				{SessionID: "session2", EventName: "world"},
				{SessionID: "session1", EventName: "how are you?"},
				{SessionID: "session2", EventName: "I'm fine, thanks!"},
			},
			session: "session2",
			expected: []callback.CallbackData{
				{SessionID: "session2", EventName: "world"},
				{SessionID: "session2", EventName: "I'm fine, thanks!"},
			},
		},
		{
			name: "empty session",
			data: []callback.CallbackData{
				{SessionID: "session1", EventName: "hello"},
				{SessionID: "session2", EventName: "world"},
			},
			session:  "session3",
			expected: []callback.CallbackData{},
		},
	}

	// run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// call the Callback method with the callback data
			for _, data := range tc.data {
				p.Callback(ctx, data)
			}

			// get the events for the session
			events := p.GetEvents(tc.session)

			// check that the events are correct
			if len(events) != len(tc.expected) {
				t.Errorf("unexpected number of events for session %s: got %d, want %d", tc.session, len(events), len(tc.expected))
			} else {
				for i, e := range events {
					if e.SessionID != tc.expected[i].SessionID || e.EventName != tc.expected[i].EventName {
						t.Errorf("unexpected event for session %s at index %d: got %v, want %v", tc.session, i, e, tc.expected[i])
					}
				}
			}
		})
	}
}
