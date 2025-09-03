package property

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSourceList_Contains(t *testing.T) {
	testCases := []struct {
		name       string
		sourceName string
		sources    []Source
		wantResult bool
	}{
		{
			name:       "source does not exist",
			sourceName: "anotherMapSource",
			sources: []Source{
				NewMapSource("anyMapSource", make(map[string]any)),
			},
			wantResult: false,
		},
		{
			name:       "source exists",
			sourceName: "anyMapSource",
			sources: []Source{
				NewMapSource("anyMapSource", make(map[string]any)),
				NewMapSource("anotherMapSource", make(map[string]any)),
			},
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.sources...)

			// when
			result := sourceList.Contains(tc.sourceName)

			// then
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestSourceList_Find(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))

	testCases := []struct {
		name       string
		sourceName string
		sources    []Source
		wantExists bool
		wantSource Source
	}{
		{
			name:       "source does not exist",
			sourceName: "anotherMapSource",
			sources: []Source{
				anyMapSource,
			},
			wantExists: false,
		},
		{
			name:       "source exists",
			sourceName: "anyMapSource",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			wantExists: true,
			wantSource: anyMapSource,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.sources...)

			// when
			source, exists := sourceList.Find(tc.sourceName)

			// then
			assert.Equal(t, tc.wantExists, exists)
			assert.Equal(t, tc.wantSource, source)
		})
	}
}

func TestSourceList_AddFirst(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapSource("anyMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []Source
		source      Source
		wantPanic   error
		wantSources []Source
	}{
		{
			name: "nil source",
			sources: []Source{
				anyMapSource,
			},
			source:    nil,
			wantPanic: errors.New("nil source"),
		},
		{
			name: "any source",
			sources: []Source{
				anyMapSource,
			},
			source:      anotherMapSource,
			wantSources: []Source{anotherMapSource, anyMapSource},
		},
		{
			name:        "with empty sources",
			sources:     []Source{},
			source:      anotherMapSource,
			wantSources: []Source{anotherMapSource},
		},
		{
			name: "existing source",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			wantSources: []Source{otherMapSource, anotherMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.sources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					sourceList.AddFirst(tc.source)
				})
				return
			}

			sourceList.AddFirst(tc.source)

			// then
			assert.Equal(t, tc.wantSources, sourceList.sources)
		})
	}
}

func TestSourceList_AddLast(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapSource("anyMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []Source
		source      Source
		wantPanic   error
		wantSources []Source
	}{
		{
			name: "nil source",
			sources: []Source{
				anyMapSource,
			},
			source:    nil,
			wantPanic: errors.New("nil source"),
		},
		{
			name: "any source",
			sources: []Source{
				anyMapSource,
			},
			source:      anotherMapSource,
			wantSources: []Source{anyMapSource, anotherMapSource},
		},
		{
			name: "existing source",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			wantSources: []Source{anotherMapSource, otherMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.wantSources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					sourceList.AddLast(tc.source)
				})
				return
			}

			sourceList.AddLast(tc.source)

			// then
			assert.Equal(t, tc.wantSources, sourceList.sources)
		})
	}
}

func TestSourceList_AddAtIndex(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapSource("otherMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []Source
		index       int
		source      Source
		wantPanic   error
		wantSources []Source
	}{
		{
			name: "negative index",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			index:     -5,
			wantPanic: errors.New("negative index"),
		},
		{
			name: "nil source",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			index:     1,
			wantPanic: errors.New("nil source"),
		},
		{
			name: "add last",
			sources: []Source{
				anyMapSource,
			},
			source:      anotherMapSource,
			index:       1,
			wantSources: []Source{anyMapSource, anotherMapSource},
		},
		{
			name: "add first",
			sources: []Source{
				anyMapSource,
			},
			source:      anotherMapSource,
			index:       0,
			wantSources: []Source{anotherMapSource, anyMapSource},
		},
		{
			name: "add at index",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			index:       1,
			wantSources: []Source{anyMapSource, otherMapSource, anotherMapSource},
		},
		{
			name: "existing source",
			sources: []Source{
				anotherMapSource,
				anyMapSource,
				otherMapSource,
			},
			source:      anyMapSource,
			index:       2,
			wantSources: []Source{anotherMapSource, otherMapSource, anyMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.sources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					sourceList.AddAtIndex(tc.index, tc.source)
				})
				return
			}

			sourceList.AddAtIndex(tc.index, tc.source)

			// then
			assert.Equal(t, tc.wantSources, sourceList.sources)
		})
	}
}

func TestSourceList_Remove(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))

	testCases := []struct {
		name       string
		sourceName string
		sources    []Source
		wantSource Source
		wantLen    int
	}{
		{
			name:       "source does not exist",
			sourceName: "anotherMapSource",
			sources: []Source{
				anyMapSource,
			},
			wantSource: nil,
			wantLen:    1,
		},
		{
			name:       "source exists",
			sourceName: "anyMapSource",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			wantSource: anyMapSource,
			wantLen:    1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.sources...)

			// when
			source := sourceList.Remove(tc.sourceName)

			// then
			assert.Equal(t, tc.wantSource, source)
			assert.Equal(t, tc.wantLen, len(sourceList.sources))
		})
	}
}

func TestSourceList_Replace(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))
	otherMapSource := NewMapSource("anyMapSource", make(map[string]any))

	testCases := []struct {
		name        string
		sources     []Source
		sourceName  string
		source      Source
		wantPanic   error
		wantSources []Source
	}{
		{
			name:       "nil source",
			sourceName: "anySourceName",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:    nil,
			wantPanic: errors.New("nil source"),
		},
		{
			name:       "replace with existing",
			sourceName: "anyMapSource",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:      otherMapSource,
			wantSources: []Source{otherMapSource, anyMapSource},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.wantSources...)

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					sourceList.Replace(tc.sourceName, tc.source)
				})
				return
			}

			sourceList.Replace(tc.sourceName, tc.source)

			// then
			assert.Equal(t, tc.wantSources, sourceList.sources)
		})
	}
}

func TestSourceList_Count(t *testing.T) {
	// given
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))
	sourceList := SourcesAsList(anyMapSource, anotherMapSource)

	// when
	count := sourceList.Count()

	// then
	assert.Equal(t, 2, count)
}

func TestSourceList_PrecedenceOf(t *testing.T) {
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))

	testCases := []struct {
		name             string
		sources          []Source
		source           Source
		wantPrecedenceOf int
	}{
		{
			name: "nil source",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:           nil,
			wantPrecedenceOf: -1,
		},
		{
			name: "source does not exist",
			sources: []Source{
				anyMapSource,
			},
			source:           anotherMapSource,
			wantPrecedenceOf: -1,
		},
		{
			name: "source exists",
			sources: []Source{
				anyMapSource,
				anotherMapSource,
			},
			source:           anotherMapSource,
			wantPrecedenceOf: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			sourceList := SourcesAsList(tc.sources...)

			// when
			precedenceOf := sourceList.PrecedenceOf(tc.source)

			// then
			assert.Equal(t, tc.wantPrecedenceOf, precedenceOf)
		})
	}
}

func TestSourceList_Slice(t *testing.T) {
	// given
	anyMapSource := NewMapSource("anyMapSource", make(map[string]any))
	anotherMapSource := NewMapSource("anotherMapSource", make(map[string]any))
	sourceList := SourcesAsList(anyMapSource, anotherMapSource)

	// when
	slice := sourceList.Slice()

	// then
	assert.Equal(t, []Source{anyMapSource, anotherMapSource}, slice)
}
