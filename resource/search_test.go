package resource

import (
	"testing"
)

func TestSearchKey_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sk   SearchKey
		want string
	}{
		{
			name: "basic search key",
			sk:   SearchKey("test_key"),
			want: "test_key",
		},
		{
			name: "empty search key",
			sk:   SearchKey(""),
			want: "",
		},
		{
			name: "search key with special characters",
			sk:   SearchKey("key-with-dashes_and_underscores"),
			want: "key-with-dashes_and_underscores",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.sk.String()
			if got != tt.want {
				t.Errorf("SearchKey.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		typ    SearchType
		values map[SearchKey]string
		want   *Search
	}{
		{
			name: "substring search",
			typ:  SubString,
			values: map[SearchKey]string{
				"title": "test search",
			},
			want: &Search{
				typ:    SubString,
				values: map[SearchKey]string{"title": "test search"},
			},
		},
		{
			name: "fulltext search",
			typ:  FullText,
			values: map[SearchKey]string{
				"content": "full text search",
			},
			want: &Search{
				typ:    FullText,
				values: map[SearchKey]string{"content": "full text search"},
			},
		},
		{
			name: "ngram search",
			typ:  Ngram,
			values: map[SearchKey]string{
				"description": "ngram search",
			},
			want: &Search{
				typ:    Ngram,
				values: map[SearchKey]string{"description": "ngram search"},
			},
		},
		{
			name:   "empty search",
			typ:    SubString,
			values: map[SearchKey]string{},
			want: &Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSearch(tt.typ, tt.values)
			if got.typ != tt.want.typ {
				t.Errorf("NewSearch() typ = %v, want %v", got.typ, tt.want.typ)
			}
			if len(got.values) != len(tt.want.values) {
				t.Errorf("NewSearch() values length = %v, want %v", len(got.values), len(tt.want.values))
			}
			for k, v := range got.values {
				if tt.want.values[k] != v {
					t.Errorf("NewSearch() values[%v] = %v, want %v", k, v, tt.want.values[k])
				}
			}
		})
	}
}

func TestSearch_spannerStmt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		search  Search
		wantErr bool
	}{
		{
			name: "substring search with single key",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title": "test search",
				},
			},
			wantErr: false,
		},
		{
			name: "substring search with multiple terms",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title": "test search multiple terms",
				},
			},
			wantErr: false,
		},
		{
			name: "substring search with empty values",
			search: Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
			wantErr: true,
		},
		{
			name: "substring search with multiple keys",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title":   "test",
					"content": "search",
				},
			},
			wantErr: true,
		},
		{
			name: "fulltext search not implemented",
			search: Search{
				typ: FullText,
				values: map[SearchKey]string{
					"title": "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ngram search not implemented",
			search: Search{
				typ: Ngram,
				values: map[SearchKey]string{
					"title": "test",
				},
			},
			wantErr: true,
		},
		{
			name: "unsupported search type",
			search: Search{
				typ: SearchType("unsupported"),
				values: map[SearchKey]string{
					"title": "test",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.search.spannerStmt()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.spannerStmt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearch_parseToSearchSubstring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		search  Search
		wantErr bool
	}{
		{
			name: "single key single term",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "single key multiple terms",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title": "test search multiple",
				},
			},
			wantErr: false,
		},
		{
			name: "empty values",
			search: Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
			wantErr: true,
		},
		{
			name: "multiple keys",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title":   "test",
					"content": "search",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.search.parseToSearchSubstring()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.parseToSearchSubstring() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearch_parseToNgramScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		search  Search
		wantErr bool
	}{
		{
			name: "single key single term",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "single key multiple terms",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title": "test search multiple",
				},
			},
			wantErr: false,
		},
		{
			name: "empty values",
			search: Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
			wantErr: true,
		},
		{
			name: "multiple keys",
			search: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"title":   "test",
					"content": "search",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.search.parseToNgramScore()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.parseToNgramScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}