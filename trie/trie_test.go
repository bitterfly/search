package trie

import (
	"testing"

	"github.com/elgs/gostrgen"
	"github.com/stretchr/testify/assert"
)

func TestTraverse_ABCD(t *testing.T) {
	assert := assert.New(t)

	trie := &Trie{
		MaxIndex: 3,
		Transitions: map[Transition]int32{
			Transition{Id: 0, Label: byte('a')}: 1,
			Transition{Id: 1, Label: byte('b')}: 2,
			Transition{Id: 2, Label: byte('c')}: 3,
		},
		Values: map[int32]int32{
			3: int32(42),
		},
	}

	destination, rest := trie.traverseWith([]byte("abcd"))

	if rest == nil {
		t.Errorf("Result from traversing with 'abcd' when 'abc' is in trie is null, but has to be 'd'\n")
	}

	if len(rest) != 1 {
		t.Errorf("Result from traversing with 'abcd' when 'abc' is in trie is %s, but has to be 'd'\n", string(rest))
	}

	assert.Equal(destination, int32(3))
}

func TestTraverse_ABC(t *testing.T) {
	assert := assert.New(t)

	trie := &Trie{
		MaxIndex: 3,
		Transitions: map[Transition]int32{
			Transition{Id: 0, Label: byte('a')}: 1,
			Transition{Id: 1, Label: byte('b')}: 2,
			Transition{Id: 2, Label: byte('c')}: 3,
		},
		Values: map[int32]int32{
			3: int32(42),
		},
	}

	destination, rest := trie.traverseWith([]byte("abc"))

	if rest != nil {
		t.Errorf("Result from traversing with 'abc' when 'abc' is in trie is not null: %s\n", string(rest))
	}

	assert.Equal(destination, int32(3))
}

func TestGet_ABC(t *testing.T) {
	assert := assert.New(t)

	trie := &Trie{
		MaxIndex: 3,
		Transitions: map[Transition]int32{
			Transition{Id: 0, Label: byte('a')}: 1,
			Transition{Id: 1, Label: byte('b')}: 2,
			Transition{Id: 2, Label: byte('c')}: 3,
		},
		Values: map[int32]int32{
			3: int32(42),
		},
	}

	assert.Equal(*trie.Get([]byte("abc")), int32(42))
}

func TestTrie_Put_and_Get(t *testing.T) {
	assert := assert.New(t)
	trie := New()

	trie.Put([]byte("foo"), 42)
	trie.Put([]byte("fob"), 43)
	trie.Put([]byte("bar"), 44)

	assert.Equal(*trie.Get([]byte("foo")), int32(42))
	assert.Equal(*trie.Get([]byte("fob")), int32(43))
	assert.Equal(*trie.Get([]byte("bar")), int32(44))
	assert.Equal(trie.GetOrPut([]byte("qux"), 5), int32(5))
	assert.Equal(*trie.Get([]byte("qux")), int32(5))
	assert.Equal(trie.GetOrPut([]byte("qux"), 1), int32(5))
}

func TestTrie_PutLambda(t *testing.T) {
	assert := assert.New(t)
	trie := New()

	trie.Put([]byte("foo"), 42)
	trie.Put([]byte("fob"), 43)
	trie.Put([]byte("bar"), 44)

	trie.PutLambda([]byte("fob"), func(x int32) int32 { return x + 10 }, 0)
	assert.Equal(*trie.Get([]byte("fob")), int32(53))
}

func TestTrie_Walk(t *testing.T) {
	assert := assert.New(t)

	trie := New()

	trie.Put([]byte("fo"), 41)
	trie.Put([]byte("foo"), 42)
	trie.Put([]byte("fob"), 43)
	trie.Put([]byte("bar"), 44)

	var words []string
	var values []int32

	trie.Walk(func(word []byte, value int32) {
		words = append(words, string(word))
		values = append(values, value)
	})

	assert.ElementsMatch(words, []string{"bar", "fo", "fob", "foo"})
	assert.ElementsMatch(values, []int32{44, 41, 43, 42})
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
		for i := int32(0); i < int32(numWords); i++ {
			trie.Put(words[i], i)
		}
	}
}

func benchmarkGet(numWords int, b *testing.B) {
	words := makeBenchmarkWords(numWords)

	trie := New()

	for i := int32(0); i < int32(numWords); i++ {
		trie.Put(words[i], i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for j := 0; j < b.N; j++ {
		for i := int32(0); i < int32(numWords); i++ {
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
