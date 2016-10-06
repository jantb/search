package proto

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var data = `support
		hello=3`

func TestEvent_Search_match_not_match(t *testing.T) {
	e := Event{
		BloomDirty:true,
		Data:data, }
	assert.Equal(t, true, e.ShouldAddAndGetIndexes([]string{"support"}))
	assert.Equal(t, false, e.ShouldAddAndGetIndexes([]string{"supports"}))
}

func TestEvent_Search_field(t *testing.T) {
	e := Event{
		BloomDirty:true,
		Data:data, }

	e.GenerateBloom()
	assert.Equal(t, true, e.ShouldAddAndGetIndexes([]string{"hello>2"}))
}