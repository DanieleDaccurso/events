package events

import (
	"sort"
	"reflect"
)

type Event interface {}

type onReturnFN func([]reflect.Value)

// DispatchEvents will dispatch a map of events, by their int values, with an empty callback.
// If you need to process the result of your events, consider using DispatchEventsCallback
func DispatchEvents(evList map[int]Event, withContext interface{} )  {
	DispatchEventsCallback(evList, withContext, func(e []reflect.Value) {})
}

// DispatchEventsCallback will dispatch a map of events, by their int values. Unlike DispatchEvent, it will
// execute the additionally passed function if the event does not have an empty return.
// If the return of your event is empty, the callback function will not be executed.
func DispatchEventsCallback(evList map[int]Event, withContext interface{}, callback onReturnFN) {
	// maps are unsorted, we need a sorted slice of map keys
	var evKeys []int
	for k := range evList {
		evKeys = append(evKeys, k)
	}
	sort.Ints(evKeys)

	// generate a reflection value of the event context
	reflectCtx := reflect.ValueOf(withContext)

	// dispatch events in correct order
	for _, key := range evKeys {
		inst := reflect.ValueOf(evList[key])
		// For performance reasons, this method won't check if there is an Exec method on the struct.
		// Using the correct type of event is done by the container of the events itself.
		// Please refer to the documentation for further instruction
		ret := inst.MethodByName("Exec").Call([]reflect.Value{reflectCtx})
		if len(ret) > 0 {
			callback(ret)
		}
	}
}