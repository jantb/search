package searchfor

import (
	"github.com/jantb/search/proto"
	"testing"
)

var data = `The Go language has built-in facilities, as well as library support,
		for writing concurrent programs. Concurrency refers not only to CPU parallelism,
		but also to asynchrony: letting slow operations like a database or network-read run
		while the program does other work, as is common in event-based servers.[64]The primary
		concurrency construct is the goroutine, a type of light-weight process. A function call
		prefixed with the go keyword starts a function in a new goroutine. The language specification
		does not specify how goroutines should be implemented, but current implementations multiplex a Go
		process's goroutines onto a smaller set of operating system threads, similar to the scheduling
		performed in Erlang.[65]:10While a standard library package featuring most of the classical
		concurrency control structures (mutex locks, etc.) is available,[65]:151–152 idiomatic concurrent
		programs instead prefer channels, which provide send messages between goroutines.[66] Optional buffers
		store messages in FIFO order[50]:43 and allow sending goroutines to proceed before their messages are
		received.Channels are typed, so that a channel of type chan T can only be used to transfer messages of
		type T. Special syntax is used to operate on them; <-ch is an expression that causes the executing
		goroutine to block until a value comes in over the channel ch, while ch <- x sends the value x
		(possibly blocking until another goroutine receives the value). The built-in switch-like select
		statement can be used to implement non-blocking communication on multiple channels; see below for an
		example. Go has a memory model describing how goroutines must use channels or other operations to
		safely share data.The existence of channels sets Go apart from actor model-style concurrent languages
		like Erlang, where messages are addressed directly to actors (corresponding to goroutines); the actor
		style can be simulated in Go by maintaining a one-to-one correspondence between goroutines and
		channels, but the language allows multiple goroutines to share a channel, or a single goroutine
		to send and receive on multiple channels.[65]:147From these tools one can build concurrent constructs
		like worker pools, pipelines (in which, say, a file is decompressed and parsed as it downloads),
		background calls with timeout, "fan-out" parallel calls to a set of services, and others.[67]
		Channels have also found uses further from the usual notion of interprocess communication, like
		serving as a concurrency-safe list of recycled buffers,[68] implementing coroutines (which helped
		inspire the name goroutine),[69] and implementing iterators.[70]Concurrency-related structural
		conventions of Go (channels and alternative channel inputs) are derived from Tony Hoare's
		communicating sequential processes model. Unlike previous concurrent programming languages such as
		Occam or Limbo (a language on which Go co-designer Rob Pike worked),[71] Go does not provide any
		built-in notion of safe or verifiable concurrency.[72] While the communicating-processes model is
		favored in Go, it is not the only one: all goroutines in a program share a single address space.
		This means that mutable objects and pointers can be shared between goroutines; see § Lack of race
		condition safety, below.`

func TestEvent_Search(t *testing.T) {
	e := proto.Event{
		BloomDirty: true,
		Data:       data}
	if !e.ShouldAddAndGetIndexes([]string{"suppor"}) {
		t.Fail()
	}
	if e.ShouldAddAndGetIndexes([]string{"supporss"}) {
		t.Fail()
	}
}

func TestEvent_Search_Not(t *testing.T) {
	e := proto.Event{
		BloomDirty: true,
		Data:       data}
	if e.ShouldAddAndGetIndexes([]string{"!suppor"}) {
		t.Fail()
	}
	if !e.ShouldAddAndGetIndexes([]string{"!supporss"}) {
		t.Fail()
	}
}

func BenchmarkEvent_Search_Worst_Case(b *testing.B) {
	e := proto.Event{
		BloomDirty: true,
		Data:       data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ShouldAddAndGetIndexes([]string{"below"})
	}
}

func BenchmarkEvent_Search_Bloom_Worst_Case(b *testing.B) {
	e := proto.Event{
		BloomDirty: true,
		Data:       data}
	e.GenerateBloom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ShouldAddAndGetIndexes([]string{"below"})
	}
}

func BenchmarkEvent_Search(b *testing.B) {
	e := proto.Event{
		BloomDirty: true,
		Data:       data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ShouldAddAndGetIndexes([]string{"Go"})
	}
}
func BenchmarkEvent_Search_Bloom(b *testing.B) {
	e := proto.Event{
		BloomDirty: true,
		Data:       data}
	e.GenerateBloom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ShouldAddAndGetIndexes([]string{"Go"})
	}
}
