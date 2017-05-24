package events

import (
	"testing"
	"reflect"
)

//////////////////////
// MOCKS
/////////////////////
type tMockCtxA struct{ a int }
type tMockEvA struct{}
type tMockEvB struct{}
type tMockInvalidEv struct {}

func (t *tMockEvA) Exec(ctx *tMockCtxA) { ctx.a = 54 }
func (t *tMockEvB) Exec(ctx *tMockCtxA) { ctx.a = 123 }

//////////////////////
// MOCKS END
/////////////////////

// BEGIN HAPPYPATH

// Begin of the "happy-path" tests
// The following tests until "END HAPPYPATH" are the tests in which the default
// expected functionality in case of correct usage is tested.

func TestEventCollection_AddEvent(t *testing.T) {
	e := new(EventCollection)

	evA := new(tMockEvA)
	evB := new(tMockEvB)
	e.AddEvent(evA, 50)
	// test counting of highest priority "hp"
	if e.hp != 50 {
		t.Error("TestEventCollection_AddEven: Did not set HP")
	}

	e.AddEvent(evB, 30)
	// test counting of highest priority "hp"
	if e.hp != 50 {
		t.Error("TestEventCollection_AddEven: Falsely updated hp")
	}

	// New events should be appended at the end of the event queue, regardless of their priority
	if e.events[0].ev != evA || e.events[1].ev != evB {
		t.Error("TestEventCollection_AddEvent: Failed event slice test")
	}

	// Sort the event collection
	e.sort()

	// After sorting the event collection, the slice should now be ordered differently
	if e.events[0].ev != evB || e.events[1].ev != evA {
		t.Error("TestEventCollection_AddEvent: Failed event slice order test")
	}
}

func TestEventCollection_AppendEvent(t *testing.T) {
	e := new(EventCollection)

	evA := new(tMockEvA)
	evB := new(tMockEvB)

	// Append first event and let hp count to 1
	e.AppendEvent(evA)
	if e.hp != 1 {
		t.Error("TestEventCollection_AppendEvent: Failed to count hp on append")
	}

	// Append second event and let hp count to 2
	e.AppendEvent(evB)
	if e.hp != 2 {
		t.Error("TestEventCollection_AppendEvent: Failed to increase hp on append")
	}

	// The events should be in the same order as appended
	if e.events[0].ev != evA || e.events[1].ev != evB {
		t.Error("TestEventCollection_AppendEvent: Failed slice order")
	}

	// After a sort, the order should not be changed, since events have only
	// been appended
	e.sort()
	if e.events[0].ev != evA || e.events[1].ev != evB {
		t.Error("TestEventCollection_AppendEvent: Failed to maintain slice order")
	}
}

func TestEventCollection_Dispatch(t *testing.T) {
	e := new(EventCollection)
	evA := new(tMockEvA)
	evB := new(tMockEvB)
	ctx := new(tMockCtxA)
	e.AppendEvent(evA)
	e.AppendEvent(evB)

	e.Dispatch(ctx)

	if ctx.a != 123 {
		t.Error("TestEventCollection_Dispatch: Failed dispatch")
	}
}

func TestEventCollection_DispatchCallback(t *testing.T) {
	counter := 0
	e := new(EventCollection)
	evB := new(tMockEvB)
	evA := new(tMockEvA)
	ctx := new(tMockCtxA)
	e.AppendEvent(evA)
	e.AppendEvent(evB)

	e.DispatchCallback(ctx, func(v []reflect.Value) {
		counter ++
	})

	if counter != 2 {
		t.Error("TestEventCollection_DispatchCallback: Failed to execute callback")
	}
}

func TestDispatchEvent(t *testing.T) {
	evA := new(tMockEvA)
	ctx := new(tMockCtxA)

	DispatchEvent(evA, ctx)

	if ctx.a != 54 {
		t.Error("TestDispatchEvent: Failed single event dispatch")
	}
}

func TestDispatchEventCallback(t *testing.T) {
	counter := 0
	evB := new(tMockEvB)
	ctx := new(tMockCtxA)

	DispatchEventCallback(evB, ctx, func(v []reflect.Value) {
		counter ++
	})

	if counter != 1 {
		t.Error("TestDispatchEventCallback: Failed to execute callback")
	}
}

// END HAPPYPATH

// The following test, tests expected failure
func TestFailEventFormat(t *testing.T) {
	iEv := new(tMockInvalidEv)
	DispatchEvent(iEv, nil)
	if LastError == "" {
		t.Error("TestFailEventFormat: Failed to fail execution")
	}
}