package burl

type EventType int

//event IDs used by burl internally
const (
	EV_UPDATE_UI    EventType = iota //signifies some UI needs to be updated
	EV_CHANGE_STATE                  //state change required --internal--
	EV_QUIT                          //shut it down! --internal--
	EV_TAB_FIELD                     //for when you tab around a UI. NOTE: should this be internal burl behaviour?
	EV_ANIMATION_DONE
	EV_BUTTON_PRESS
	EV_LIST_CYCLE
	EV_MAX_EVENTS
)

//number of custom events.
var customEventNum int

func RegisterCustomEvent() EventType {
	customEventNum++
	return EventType(int(EV_MAX_EVENTS) + customEventNum)
}

//flags set to true for events that are internal only
var internalEvent map[EventType]bool

//Event is the basic unit of the EventStream.
type Event struct {
	ID      EventType
	Message string
	Caller  UIElem //for UI events, record what emitted the event.
}

//Creates a new generic event, with no UI caller.
func NewEvent(id EventType, m string) *Event {
	return &Event{id, m, nil}
}

//Creates a new UI event. Should be emitted internally by the UI system, but consumed externally.
func NewUIEvent(id EventType, m string, call UIElem) *Event {
	return &Event{id, m, call}
}

//EventStream is the queue of Events produced and consumed by the application. Burl may
//also emit events into this stream for optional consumption. Example: animation completion
//events - applicaton may like to enact some behaviour the frame an animation completes.
var eventStream chan *Event

//eventStreamInternal is the queue of events used for Burl-level events.
//Application can also emit Burl-level events to be handled by the engine.
//Example: application can emit QUIT_EVENT, which burl consumes instead of the application.
//All Burl-Events are consumed at the end of the frame.
var eventStreamInternal chan *Event

//Sometimes we end up putting like 10000 UI_UPDATEs in the queue and break it, so we keep track of
//which UPDATE_UI_EVENTs we push into the queue to ensure we don't duplicate them. They are removed from
//the map once popped. There never needs to be more than one since the events SHOULD be being consumed
//once per frame.
var uiEvents map[string]bool

//Allocate event buffer. 1000 events should be enough, right???
func init() {
	eventStream = make(chan *Event, 1000)
	eventStreamInternal = make(chan *Event, 1000)
	uiEvents = make(map[string]bool, 100)
	internalEvent = make(map[EventType]bool, EV_MAX_EVENTS)

	//set which events types are internal to burl
	internalEvent[EV_QUIT] = true
	internalEvent[EV_CHANGE_STATE] = true
	customEventNum = 0
}

//Emits an event into the relevant EventStream. If the stream is full we flush the whole buffer.
//TODO: Is flushing the buffer a little barbaric? We could maybe just consume half of them or something.
//Or maybe just push the oldest one out? Rotate the buffer? Think on this.
func PushEvent(e *Event) {
	//special processing for certain IDs
	switch e.ID {
	case EV_UPDATE_UI:
		if uiEvents[e.Message] {
			return
		}
		uiEvents[e.Message] = true
	}

	if internalEvent[e.ID] {
		addEvent(e, eventStreamInternal)
	} else {
		addEvent(e, eventStream)
	}
}

func addEvent(e *Event, stream chan *Event) {
	if len(stream) == cap(stream) {
		stream = make(chan *Event, 1000)
		LogError("stream buffer overflow! all events flushed. Oh no.")
	}
	stream <- e
}

//Reallocates the eventstream.
func ClearEvents() {
	eventStream = make(chan *Event, 1000)
}

//Reallocates the internal eventstream.
func clearEventsInternal() {
	eventStreamInternal = make(chan *Event, 1000)
}

//Grabs an event from the stream for consumption purposes.
func PopEvent() *Event {
	if len(eventStream) > 0 {
		e := <-eventStream

		if e.ID == EV_UPDATE_UI {
			uiEvents[e.Message] = false
		}

		return e
	} else {
		return nil
	}
}

func popInternalEvent() *Event {
	if len(eventStreamInternal) > 0 {
		e := <-eventStreamInternal
		return e
	} else {
		return nil
	}
}
