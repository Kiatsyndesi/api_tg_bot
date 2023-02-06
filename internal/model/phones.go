package model

type Phone struct {
	ID    uint64
	Title string
}

type EventType uint8

type EventStatus uint8

const (
	Created = iota
	Updated
	Removed

	Deffered EventStatus = iota
	Processed
)

type PhonesEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Phone
}
