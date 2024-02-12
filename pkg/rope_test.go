package pkg

import (
	"os"
	"testing"

	"golang.org/x/exp/rand"
)

const (
	text       string = "Hello to this beaufitul world!"
	panTadeusz string = "pan_tadeusz.txt"
	chunk      int    = 15
)

func TestRope_String(t *testing.T) {
	texts := []string{"Hello world", "", "7_7.7__7  && 777"}
	ropes := []*Rope{
		{value: []byte("Hello world")},
		{},
		{value: []byte("7_7.7__7  && 777")},
	}
	for i, text := range texts {
		if ropes[i].String() != text {
			t.Errorf("rope.String() = %v; want %v", ropes[i], text)
		}
	}
}

func TestNewRope(t *testing.T) {
	got := FromStringWith(text, chunk)
	if got == nil {
		t.Errorf("incorrect initization")
	}

	if string(got.value) != "" {
		t.Errorf("value got %v; want %v", got.value, "")
	}

	if string(got.left.value) != text[:chunk] {
		t.Errorf("left child got %v; want %v", got.left.value, text[:chunk])
	}

	if got.weight != len(text[:chunk]) {
		t.Errorf("weight got %v; want %v", got.weight, len(text[:chunk]))
	}

	if string(got.right.value) != text[chunk:] {
		t.Errorf("right child got %v; want %v", got.right.value, text[chunk:])
	}
}

func TestRope_Concat(t *testing.T) {
	left := FromStringWith(text[:chunk], chunk/2)
	right := FromStringWith(text[chunk:], chunk/2)
	got := left.Concat(right)

	if left.String() != text[:chunk] {
		t.Errorf("unexpected left modification; got %v; want %v", left, text[:chunk])
	}

	if right.String() != text[chunk:] {
		t.Errorf("unexpected right modification; got %v; want %v", right, text[chunk:])
	}

	if got.String() != text {
		t.Errorf("got %v; want %v", got, text)
	}
}

func TestRope_At(t *testing.T) {
	rope := FromStringWith(text, chunk)
	for i, c := range text {
		got, err := rope.At(i)
		if err != nil {
			t.Errorf("error got %v", err)
		}

		if rune(got) != c {
			t.Errorf("value got %v; want %v", got, c)
		}
	}
}

func TestRope_AtErr(t *testing.T) {
	rope := FromStringWith(text, chunk)
	_, err := rope.At(-1)
	if err == nil {
		t.Errorf("expected error upon reading negative position")
	}

	_, err = rope.At(len(text))
	if err == nil {
		t.Errorf("expected error upon out of bounds position")
	}
}

func TestRope_Len(t *testing.T) {
	got := FromStringWith(text, chunk).Len()
	if got != len(text) {
		t.Errorf("got %v; want %v", got, len(text))
	}

	got = FromStringWith(text, 2).Len()
	if got != len(text) {
		t.Errorf("got %v; want %v", got, len(text))
	}
}

func TestRope_CloneWith(t *testing.T) {
	rope := FromStringWith(text, chunk)
	clone := rope.CloneWith(chunk)
	if rope == clone {
		t.Errorf("rope and clone are the same object")
	}

	if clone.String() != rope.String() {
		t.Errorf("got %v; want %v", clone.String(), rope.String())
	}
}

func TestRope_Split(t *testing.T) {
	rope := FromStringWith(text, chunk)
	for pos := range text {
		left := rope.CloneWith(chunk)
		right := left.Split(pos)

		if left.String() != text[:pos] {
			t.Errorf("left got %v; want %v", left.String(), text[:pos])
		}

		if right.String() != text[pos:] {
			t.Errorf("right got %v; want %v", right.String(), text[pos:])
		}
	}
}

func TestRope_AtAfterSplit(t *testing.T) {
	left := FromStringWith(text, chunk)
	pos := 10

	right := left.Split(pos)

	for i, c := range text[:pos] {
		got, err := left.At(i)
		if err != nil {
			t.Errorf("error got %v", err)
		}

		if rune(got) != c {
			t.Errorf("value got %v; want %v", got, c)
		}
	}

	for i, c := range text[pos:] {
		got, err := right.At(i)
		if err != nil {
			t.Errorf("error got %v", err)
		}

		if rune(got) != c {
			t.Errorf("value got %v; want %v", got, c)
		}
	}
}

func TestRope_Insert(t *testing.T) {
	rope := FromStringWith(text, chunk)
	banana := "Banana!"
	for pos := range text {
		got := rope.CloneWith(chunk)
		want := text[:pos] + banana + text[pos:]
		err := got.Insert(pos, []byte(banana))
		if err != nil {
			t.Errorf("unexpected error in insert operation %v", err)
		}

		if got.String() != want {
			t.Errorf("got %v; want %v", got, want)
		}
	}
}

func TestRope_Delete(t *testing.T) {
	rope := FromStringWith(text, chunk)
	ns := []int{1, 2, 3}
	for _, n := range ns {
		for pos := range text {
			got := rope.CloneWith(chunk)
			want := text[:pos] + text[min(pos+n, len(text)):]
			err := got.Delete(pos, n)
			if err != nil {
				t.Errorf("unexpected error in delete operation%v", err)
			}

			if got.String() != want {
				t.Errorf("got %v; want %v", got, want)
			}
		}
	}
}

func BenchmarkString_At(b *testing.B) {
	file, err := os.ReadFile(panTadeusz)
	if err != nil {
		b.Errorf("unexpected error when reading %v, %v", panTadeusz, err)
	}

	str := string(file)
	length := len(str)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := rand.Intn(length)
		_ = str[pos]
	}
}

func BenchmarkRope_At(b *testing.B) {
	file, err := os.ReadFile(panTadeusz)
	if err != nil {
		b.Errorf("unexpected error when reading %v, %v", panTadeusz, err)
	}

	rope := NewWith(file, chunk)
	length := rope.Len()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := rand.Intn(length)
		_, _ = rope.At(pos)
	}
}

func BenchmarkString_Insert(b *testing.B) {
	file, err := os.ReadFile(panTadeusz)
	if err != nil {
		b.Errorf("unexpected error when reading %v, %v", panTadeusz, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str := string(file)
		pos := rand.Intn(len(str))
		_ = str[:pos] + text + str[pos:]
	}
}

func BenchmarkRope_Insert(b *testing.B) {
	file, err := os.ReadFile(panTadeusz)
	if err != nil {
		b.Errorf("unexpected error when reading %v, %v", panTadeusz, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rope := New(file)
		pos := rand.Intn(rope.Len())
		_ = rope.Insert(pos, []byte(text))
	}
}
