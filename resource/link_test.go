package resource

import (
	"encoding/json"
	"testing"

	"github.com/cccteam/ccc"
)

func TestLink_EncodeSpanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		link    Link
		want    any
		wantErr bool
	}{
		{
			name: "valid link",
			link: Link{
				ID:       func() ccc.UUID { id, _ := ccc.NewUUID(); return id }(),
				Resource: "test-resource",
				Text:     "test text",
			},
			wantErr: false,
		},
		{
			name: "empty link",
			link: Link{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.link.EncodeSpanner()
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.EncodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Link.EncodeSpanner() returned nil")
			}
		})
	}
}

func TestLink_DecodeSpanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		link    *Link
		val     any
		wantErr bool
	}{
		{
			name: "valid string value",
			link: &Link{},
			val:  `{"id":"123e4567-e89b-12d3-a456-426614174000","resource":"test","text":"test text"}`,
			wantErr: false,
		},
		{
			name: "invalid type",
			link: &Link{},
			val:  123,
			wantErr: true,
		},
		{
			name: "nil value",
			link: &Link{},
			val:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.link.DecodeSpanner(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.DecodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLink_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		link    Link
		wantErr bool
	}{
		{
			name: "valid link",
			link: Link{
				ID:       func() ccc.UUID { id, _ := ccc.NewUUID(); return id }(),
				Resource: "test-resource",
				Text:     "test text",
			},
			wantErr: false,
		},
		{
			name: "empty link",
			link: Link{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.link.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Link.MarshalJSON() returned nil")
			}

			var unmarshaled Link
			if err := json.Unmarshal(got, &unmarshaled); err != nil {
				t.Errorf("Link.MarshalJSON() produced invalid JSON: %v", err)
			}
		})
	}
}

func TestLink_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		link    *Link
		data    []byte
		wantErr bool
	}{
		{
			name: "valid JSON",
			link: &Link{},
			data: []byte(`{"id":"123e4567-e89b-12d3-a456-426614174000","resource":"test","text":"test text"}`),
			wantErr: false,
		},
		{
			name: "null JSON",
			link: &Link{},
			data: []byte("null"),
			wantErr: false,
		},
		{
			name: "nil data",
			link: &Link{},
			data: nil,
			wantErr: false,
		},
		{
			name: "invalid JSON",
			link: &Link{},
			data: []byte("invalid json"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.link.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLink_IsNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		link Link
		want bool
	}{
		{
			name: "null link with nil ID",
			link: Link{},
			want: true,
		},
		{
			name: "valid link with ID",
			link: Link{
				ID:       func() ccc.UUID { id, _ := ccc.NewUUID(); return id }(),
				Resource: "test",
				Text:     "test",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.link.IsNull()
			if got != tt.want {
				t.Errorf("Link.IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullLink_EncodeSpanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nullLink NullLink
		want    any
		wantErr bool
	}{
		{
			name: "valid null link",
			nullLink: NullLink{
				Link: Link{
					ID:       func() ccc.UUID { id, _ := ccc.NewUUID(); return id }(),
					Resource: "test",
					Text:     "test",
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "invalid null link",
			nullLink: NullLink{
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.nullLink.EncodeSpanner()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.EncodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.nullLink.Valid && got == nil {
				t.Errorf("NullLink.EncodeSpanner() returned nil for valid link")
			}
			if !tt.nullLink.Valid && got != nil {
				t.Errorf("NullLink.EncodeSpanner() returned non-nil for invalid link")
			}
		})
	}
}

func TestNullLink_DecodeSpanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nullLink *NullLink
		val     any
		wantErr bool
	}{
		{
			name: "valid string value",
			nullLink: &NullLink{},
			val:  `{"id":"123e4567-e89b-12d3-a456-426614174000","resource":"test","text":"test text"}`,
			wantErr: false,
		},
		{
			name: "nil value",
			nullLink: &NullLink{},
			val:  nil,
			wantErr: false,
		},
		{
			name: "pointer to nil string",
			nullLink: &NullLink{},
			val:  (*string)(nil),
			wantErr: false,
		},
		{
			name: "invalid type",
			nullLink: &NullLink{},
			val:  123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.nullLink.DecodeSpanner(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.DecodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNullLink_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nullLink NullLink
		want    []byte
		wantErr bool
	}{
		{
			name: "valid null link",
			nullLink: NullLink{
				Link: Link{
					ID:       func() ccc.UUID { id, _ := ccc.NewUUID(); return id }(),
					Resource: "test",
					Text:     "test",
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "invalid null link",
			nullLink: NullLink{
				Valid: false,
			},
			want:    []byte("null"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.nullLink.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.nullLink.Valid && string(got) != "null" {
				t.Errorf("NullLink.MarshalJSON() = %v, want %v", string(got), "null")
			}
		})
	}
}

func TestNullLink_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nullLink *NullLink
		data    []byte
		wantErr bool
	}{
		{
			name: "valid JSON",
			nullLink: &NullLink{},
			data: []byte(`{"id":"123e4567-e89b-12d3-a456-426614174000","resource":"test","text":"test text"}`),
			wantErr: false,
		},
		{
			name: "null JSON",
			nullLink: &NullLink{},
			data: []byte("null"),
			wantErr: false,
		},
		{
			name: "nil data",
			nullLink: &NullLink{},
			data: nil,
			wantErr: false,
		},
		{
			name: "invalid JSON",
			nullLink: &NullLink{},
			data: []byte("invalid json"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.nullLink.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}