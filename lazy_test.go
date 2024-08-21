package main

import (
	"fmt"
	"log"
	"maps"
	"math/rand/v2"
	"testing"

	pb "lazy-proto/protos"

	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/proto"
)

// Generate message with random repeated field only
func genRandomProto(name string, size int, withMap bool) (*pb.OuterMessage, []byte) {
	message := &pb.OuterMessage{
		Name:  name,
		Inner: make([]*pb.InnerMessage, size),
		Map:   make(map[string]*pb.InnerMessage, size),
	}

	for i := range message.Inner {
		val := rand.Int32()
		// Fill repeated field
		message.Inner[i] = &pb.InnerMessage{Val: val}

		if withMap {
			// Fill map
			key := fmt.Sprintf("key-%d", val)
			message.Map[key] = &pb.InnerMessage{Val: val}
		}
	}

	serializedMsg, err := proto.Marshal(message)
	if err != nil {
		log.Fatalf("error marshalling proto: %v", err)
	}
	return message, serializedMsg
}

// Merge regular messages
func mergeMessages(msg1, msg2 *pb.OuterMessage) *pb.OuterMessage {
	merged := &pb.OuterMessage{
		Map: make(map[string]*pb.InnerMessage),
	}
	merged.Name = msg1.Name
	merged.Inner = append(msg1.Inner, msg2.Inner...)

	maps.Copy(merged.Map, msg1.Map)
	maps.Copy(merged.Map, msg2.Map)

	return merged
}

// Merge lazy messages
func mergeLazyMessages(msg1, msg2 *pb.LazyOuterMessage) *pb.LazyOuterMessage {
	merged := &pb.LazyOuterMessage{}
	merged.Name = msg1.Name
	merged.Inner = append(msg1.Inner, msg2.Inner...)
	merged.Map = append(msg1.Map, msg2.Map...)
	return merged
}

func TestLazyRepeated(t *testing.T) {
	msg1, serializedMsg1 := genRandomProto("Hello, world.", 1000, false)
	msg2, serializedMsg2 := genRandomProto("Don't panic", 2000, false)

	// Expected merge (manual)
	mergedMsg := mergeMessages(msg1, msg2)
	expectedBytes, err := proto.Marshal(mergedMsg)
	assert.Nilf(t, err, "error marshalling proto")

	// Unmarshal as Lazy
	lazyMsg1 := &pb.LazyOuterMessage{}
	err = proto.Unmarshal(serializedMsg1, lazyMsg1)
	assert.Nilf(t, err, "error lazy-marshalling proto")

	// Unmarshal as Lazy
	lazyMsg2 := &pb.LazyOuterMessage{}
	err = proto.Unmarshal(serializedMsg2, lazyMsg2)
	assert.Nilf(t, err, "error lazy-marshalling proto")

	// Merge Lazy messages
	mergedLazyMsg := mergeLazyMessages(lazyMsg1, lazyMsg2)
	actualBytes, err := proto.Marshal(mergedLazyMsg)
	assert.Nilf(t, err, "error marshalling proto")

	// Exact serialized match
	assert.EqualValues(t, expectedBytes, actualBytes)
}

func TestLazyWithMap(t *testing.T) {
	msg1, serializedMsg1 := genRandomProto("Hello, world.", 1000, true)
	msg2, serializedMsg2 := genRandomProto("Don't panic", 1500, true)

	// Expected merge (manual)
	mergedMsg := mergeMessages(msg1, msg2)
	expectedBytes, err := proto.Marshal(mergedMsg)
	assert.Nilf(t, err, "error marshalling proto")

	// Unmarshal as Lazy
	lazyMsg1 := &pb.LazyOuterMessage{}
	err = proto.Unmarshal(serializedMsg1, lazyMsg1)
	assert.Nilf(t, err, "error lazy-marshalling proto")

	// Unmarshal as Lazy
	lazyMsg2 := &pb.LazyOuterMessage{}
	err = proto.Unmarshal(serializedMsg2, lazyMsg2)
	assert.Nilf(t, err, "error lazy-marshalling proto")

	// Merge Lazy messages
	mergedLazyMsg := mergeLazyMessages(lazyMsg1, lazyMsg2)
	actualBytes, err := proto.Marshal(mergedLazyMsg)
	assert.Nilf(t, err, "error marshalling proto")

	// With Maps, ordered is not guaranteed so we cannot compare serialized protos
	// assert.EqualValues(t, expectedBytes, actualBytes)

	// Instead we compare deserialized version
	actual := &pb.OuterMessage{}
	err = proto.Unmarshal(actualBytes, actual)
	assert.Nilf(t, err, "error unmarshalling proto")

	expected := &pb.OuterMessage{}
	err = proto.Unmarshal(expectedBytes, expected)
	assert.Nilf(t, err, "error unmarshalling proto")

	assert.EqualValues(t, expected, actual)
}

var merged []byte

// Benchmark assumes we want to merge 2 protos which are already serialized
func BenchmarkMerge(b *testing.B) {

	b.Run("Naive", func(b *testing.B) {
		var res []byte
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			_, serializedMsg1 := genRandomProto("Hello, world.", 10000, true)
			_, serializedMsg2 := genRandomProto("Don't panic", 10000, true)
			b.StartTimer()

			msg1 := &pb.OuterMessage{}
			_ = proto.Unmarshal(serializedMsg1, msg1)
			msg2 := &pb.OuterMessage{}
			_ = proto.Unmarshal(serializedMsg2, msg2)

			out := mergeMessages(msg1, msg2)
			res, _ = proto.Marshal(out)
		}
		merged = res
	})

	b.Run("Lazy", func(b *testing.B) {
		var res []byte
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			_, serializedMsg1 := genRandomProto("Hello, world.", 10000, true)
			_, serializedMsg2 := genRandomProto("Don't panic", 10000, true)
			b.StartTimer()

			msg1 := &pb.LazyOuterMessage{}
			_ = proto.Unmarshal(serializedMsg1, msg1)
			msg2 := &pb.LazyOuterMessage{}
			_ = proto.Unmarshal(serializedMsg2, msg2)

			out := mergeLazyMessages(msg1, msg2)
			res, _ = proto.Marshal(out)
		}
		merged = res
	})
}
