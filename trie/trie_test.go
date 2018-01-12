package trie

import (
	"testing"

	"github.com/elgs/gostrgen"
)

func TestTraverse_ABCD(t *testing.T) {
	trie := &Trie{
		maxIndex: 3,
		transitions: map[Transition]int32{
			Transition{id: 0, label: byte('a')}: 1,
			Transition{id: 1, label: byte('b')}: 2,
			Transition{id: 2, label: byte('c')}: 3,
		},
		values: map[int32]uint64{
			3: uint64(42),
		},
	}

	destination, rest := trie.traverseWith([]byte("abcd"))

	if rest == nil {
		t.Errorf("Result from traversing with 'abcd' when 'abc' is in trie is null, but has to be 'd'\n")
	}

	if len(rest) != 1 {
		t.Errorf("Result from traversing with 'abcd' when 'abc' is in trie is %s, but has to be 'd'\n", string(rest))
	}

	if destination != int32(3) {
		t.Errorf("Destination is %d but it has to be 3\n", destination)
	}
}

func TestTraverse_ABC(t *testing.T) {
	trie := &Trie{
		maxIndex: 3,
		transitions: map[Transition]int32{
			Transition{id: 0, label: byte('a')}: 1,
			Transition{id: 1, label: byte('b')}: 2,
			Transition{id: 2, label: byte('c')}: 3,
		},
		values: map[int32]uint64{
			3: uint64(42),
		},
	}

	destination, rest := trie.traverseWith([]byte("abc"))

	if rest != nil {
		t.Errorf("Result from traversing with 'abc' when 'abc' is in trie is not null: %s\n", string(rest))
	}

	if destination != int32(3) {
		t.Errorf("Destination is %d but it has to be 3\n", destination)
	}
}

func TestGet_ABC(t *testing.T) {
	trie := &Trie{
		maxIndex: 3,
		transitions: map[Transition]int32{
			Transition{id: 0, label: byte('a')}: 1,
			Transition{id: 1, label: byte('b')}: 2,
			Transition{id: 2, label: byte('c')}: 3,
		},
		values: map[int32]uint64{
			3: uint64(42),
		},
	}

	value := trie.Get([]byte("abc"))

	if *value != uint64(42) {
		t.Errorf("Should get value=42 when traversing with word in the trie, but got %d instead\n", *value)
	}
}

func TestTrie_Put_and_Get(t *testing.T) {
	trie := New()

	trie.Put([]byte("foo"), 42)
	trie.Put([]byte("fob"), 43)
	trie.Put([]byte("bar"), 44)

	if *trie.Get([]byte("foo")) != 42 {
		t.Errorf("cannot get foo")
	}

	if *trie.Get([]byte("fob")) != 43 {
		t.Errorf("cannot get fob")
	}

	if *trie.Get([]byte("bar")) != 44 {
		t.Errorf("cannot get bar")
	}

	if trie.GetOrPut([]byte("qux"), 5) != 5 {
		t.Errorf("cannot put or get qux when it's not present")
	}

	if *trie.Get([]byte("qux")) != 5 {
		t.Errorf("cannot get qux")
	}

	if trie.GetOrPut([]byte("qux"), 1) != 5 {
		t.Errorf("cannot put or get qux when it is present")
	}
}

func makeBenchmarkWords(size int) [][]byte {
	randomStrings := make([][]byte, size)

	for i := 0; i < size; i++ {
		str, err := gostrgen.RandGen(100, gostrgen.Lower, "", "")
		if err != nil {
			str = "abc"
		}
		randomStrings[i] = []byte(str)
	}

	return randomStrings
}

func benchmarkPut(numWords int, b *testing.B) {
	words := makeBenchmarkWords(numWords)

	trie := New()

	b.ResetTimer()
	b.ReportAllocs()

	for j := 0; j < b.N; j++ {
		for i := uint64(0); i < uint64(numWords); i++ {
			trie.Put(words[i], i)
		}
	}
}

func benchmarkGet(numWords int, b *testing.B) {
	words := makeBenchmarkWords(numWords)

	trie := New()

	for i := uint64(0); i < uint64(numWords); i++ {
		trie.Put(words[i], i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for j := 0; j < b.N; j++ {
		for i := uint64(0); i < uint64(numWords); i++ {
			trie.Get(words[i])
		}
	}
}

func BenchmarkPut10(b *testing.B)    { benchmarkPut(10, b) }
func BenchmarkPut100(b *testing.B)   { benchmarkPut(100, b) }
func BenchmarkPut1000(b *testing.B)  { benchmarkPut(1000, b) }
func BenchmarkPut10000(b *testing.B) { benchmarkPut(10000, b) }

func BenchmarkGet10(b *testing.B)    { benchmarkGet(10, b) }
func BenchmarkGet100(b *testing.B)   { benchmarkGet(100, b) }
func BenchmarkGet1000(b *testing.B)  { benchmarkGet(1000, b) }
func BenchmarkGet10000(b *testing.B) { benchmarkGet(10000, b) }
