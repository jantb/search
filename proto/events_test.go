package proto

import (
	"testing"
)

var data = `support
		hello=3`

func TestEvent_Search_match_not_match(t *testing.T) {
	e := Event{
		BloomDirty: false,
		Data:       data}
	e.GenerateBloom()
	if !e.ShouldAddAndGetIndexes([]string{"support"}) {
		t.Fail()
	}
	//if e.ShouldAddAndGetIndexes([]string{"supports"}) {
	//	t.Fail()
	//}
}

func TestEvent_Search_field(t *testing.T) {
	e := Event{
		BloomDirty: false,
		Data:       data}

	e.GenerateBloom()
	if !e.ShouldAddAndGetIndexes([]string{"hello>2"}) {
		t.Fail()
	}
}
