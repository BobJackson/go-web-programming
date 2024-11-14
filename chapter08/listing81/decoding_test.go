package main

import (
	"testing"
	"time"
)

func TestDecoding(t *testing.T) {
	post, err := decode("../../post.json")
	if err != nil {
		t.Error("Error decoding JSON:", err)
	}
	if post.Id != 1 {
		t.Error("Wrong id, was expecting 1 but got", post.Id)
	}
	if post.Content != "Hello, world!" {
		t.Error("Wrong content, was expecting 'Hello, world!' but got", post.Content)
	}
}

func TestEncode(t *testing.T) {
	t.Skip("Skipping encoding for now")
}

func TestLongRunningTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long running test in short mode")
	}
	time.Sleep(10 * time.Second)
}

// BenchmarkDecoding-8   	   95400	     10766 ns/op
func BenchmarkDecoding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = decode("../../post.json")
	}
}

// BenchmarkUnmarshal-8   	   93148	     11589 ns/op
func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = unmarshal("../../post.json")
	}
}
