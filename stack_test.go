package fstack

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestStackSimple(t *testing.T) {
	stack, err := CreateStack("temp.stack")
	if err != nil {
		t.Fatal(err)
	}
	if stack.Depth() != 0 {
		t.Fatal("Corrupted empty stack")
	}
	_, err = stack.Push(nil, []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	stack.Close()

	stack, err = OpenStack("temp.stack")
	if err != nil {
		t.Fatal(err)
	}
	if stack.Depth() != 1 {
		t.Fatal("Corrupted no-empty stack")
	}
	_, data, err := stack.Peak()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("Non-consisten data when peak: expected [hello], got %v. Block: %#v", string(data), stack.currentBlock)
	}
	if stack.Depth() != 1 {
		t.Fatal("Peak must not corrupt stack")
	}
	_, data, err = stack.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatal("Non-consisten data when pop")
	}
	stack.Close()
}

func TestStackMultiPush(t *testing.T) {
	N := 100
	data := "AAABBBCCC"
	header := "112233"
	stack, err := CreateStack("temp.stack")
	defer stack.Close()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < N; i++ {
		depth, err := stack.Push([]byte(header), []byte(data))
		if err != nil {
			t.Fatal(err)
		}
		if i+1 != depth {
			t.Fatal("Expected push returns depth")
		}
	}
	if stack.Depth() != N {
		t.Fatal("Not all elements updated depth: current is", stack.Depth(), "but expected", N)
	}
}

func TestStackIterate(t *testing.T) {
	N := 100
	data := "AAABBBCCC"
	header := "112233"
	stack, err := CreateStack("temp.stack")
	defer stack.Close()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < N; i++ {
		head := fmt.Sprintf("%v-%v", header, i)
		body := fmt.Sprintf("%v-%v", data, i)
		depth, err := stack.Push([]byte(head), []byte(body))
		if err != nil {
			t.Fatal(err)
		}
		if i+1 != depth {
			t.Fatal("Expected push returns depth")
		}
	}
	if stack.Depth() != N {
		t.Fatal("Not all elements updated depth: current is", stack.Depth(), "but expected", N)
	}
	// Now iterate
	i := 0
	stack.IterateForward(func(depth int, headStream io.Reader, body io.Reader) bool {
		head, err := ioutil.ReadAll(headStream)
		if err != nil {
			t.Fatal(err)
		}
		content, err := ioutil.ReadAll(body)
		if err != nil {
			t.Fatal(err)
		}
		// Expected
		ehead := fmt.Sprintf("%v-%v", header, depth)
		ebody := fmt.Sprintf("%v-%v", data, depth)
		if ehead != string(head) {
			t.Fatal("Unexpected header", ehead, "!=", string(head))
		}
		if ebody != string(content) {
			t.Fatal("Unexpected body", ebody, "!=", string(content))
		}
		i++
		return true
	})
}

func TestStackMultiGet(t *testing.T) {
	N := 100
	data := "AAABBBCCC"
	header := "112233"
	TestStackMultiPush(t)
	stack, err := OpenStack("temp.stack")
	if err != nil {
		t.Fatal(err)
	}
	defer stack.Close()
	if stack.Depth() != N {
		t.Fatal("Not all elements updated depth: current is", stack.Depth(), "but expected", N)
	}
	for i := N - 1; i >= 0; i-- {
		h, d, err := stack.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if i != stack.Depth() {
			t.Fatal("Expected pop returns depth")
		}
		if string(h) != header {
			t.Fatal("header not matched")
		}
		if string(d) != data {
			t.Fatal("body not matched")
		}
	}
	if stack.Depth() != 0 {
		t.Fatal("Not all elements poped: current is", stack.Depth(), "but expected", 0)
	}
}

func BenchmarkStackMultiPush(b *testing.B) {
	stack, err := CreateStack("temp.stack")
	if err != nil {
		b.Fatal(err)
	}
	data := "AAABBBCCC"
	header := "112233"
	defer stack.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depth, err := stack.Push([]byte(header), []byte(data))
		if err != nil {
			b.Fatal(err)
		}
		if i+1 != depth {
			b.Fatal("Expected push returns depth")
		}
	}
}

func BenchmarkStackMultiPop(b *testing.B) {
	BenchmarkStackMultiPush(b)
	stack, err := OpenStack("temp.stack")
	if err != nil {
		b.Fatal(err)
	}
	data := "AAABBBCCC"
	header := "112233"
	defer stack.Close()
	b.ResetTimer()
	for i := b.N - 1; i >= 0; i-- {
		h, d, err := stack.Pop()
		if err != nil {
			b.Fatal(err)
		}
		if i != stack.Depth() {
			b.Fatal("Expected pop returns depth")
		}
		if string(h) != header {
			b.Fatal("header not matched")
		}
		if string(d) != data {
			b.Fatal("body not matched")
		}

	}
}

func BenchmarkStackMultiPeak(b *testing.B) {
	stack, err := CreateStack("temp.stack")
	if err != nil {
		b.Fatal(err)
	}
	data := "AAABBBCCC"
	header := "112233"
	defer stack.Close()
	depth, err := stack.Push([]byte(header), []byte(data))
	if err != nil {
		b.Fatal(err)
	}
	if 1 != depth {
		b.Fatal("Expected push returns depth")
	}
	b.ResetTimer()
	for i := b.N - 1; i >= 0; i-- {
		h, d, err := stack.Peak()
		if err != nil {
			b.Fatal(err)
		}
		if 1 != stack.Depth() {
			b.Fatal("Expected peak  no change depth")
		}
		if string(h) != header {
			b.Fatal("header not matched")
		}
		if string(d) != data {
			b.Fatal("body not matched")
		}

	}
}

func BenchmarkStackMultiPushPop(b *testing.B) {
	stack, err := CreateStack("temp.stack")
	if err != nil {
		b.Fatal(err)
	}
	data := "AAABBBCCC"
	header := "112233"
	defer stack.Close()
	b.ResetTimer()
	for i := b.N - 1; i >= 0; i-- {
		depth, err := stack.Push([]byte(header), []byte(data))
		if err != nil {
			b.Fatal(err)
		}
		if depth != 1 {
			b.Fatal("Expected push returns depth")
		}
		h, d, err := stack.Pop()
		if err != nil {
			b.Fatal(err)
		}
		if stack.Depth() != 0 {
			b.Fatal("Expected Pop returns depth")
		}
		if string(h) != header {
			b.Fatal("header not matched")
		}
		if string(d) != data {
			b.Fatal("body not matched")
		}

	}
}
