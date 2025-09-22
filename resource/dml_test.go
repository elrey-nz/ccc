package resource

import (
	"testing"
)

func TestConfig_SetDBType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config Config
		dbType DBType
		want   Config
	}{
		{
			name:   "set spanner db type",
			config: Config{},
			dbType: SpannerDBType,
			want:   Config{DBType: SpannerDBType},
		},
		{
			name:   "set postgres db type",
			config: Config{},
			dbType: PostgresDBType,
			want:   Config{DBType: PostgresDBType},
		},
		{
			name:   "override existing db type",
			config: Config{DBType: SpannerDBType},
			dbType: PostgresDBType,
			want:   Config{DBType: PostgresDBType},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.config.SetDBType(tt.dbType)
			if got != tt.want {
				t.Errorf("Config.SetDBType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetChangeTrackingTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		config              Config
		changeTrackingTable string
		want                Config
	}{
		{
			name:                "set change tracking table",
			config:              Config{},
			changeTrackingTable: "change_tracking",
			want:                Config{ChangeTrackingTable: "change_tracking"},
		},
		{
			name:                "set empty change tracking table",
			config:              Config{},
			changeTrackingTable: "",
			want:                Config{ChangeTrackingTable: ""},
		},
		{
			name:                "override existing change tracking table",
			config:              Config{ChangeTrackingTable: "old_table"},
			changeTrackingTable: "new_table",
			want:                Config{ChangeTrackingTable: "new_table"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.config.SetChangeTrackingTable(tt.changeTrackingTable)
			if got != tt.want {
				t.Errorf("Config.SetChangeTrackingTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetTrackChanges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		config       Config
		trackChanges bool
		want         Config
	}{
		{
			name:         "enable track changes",
			config:       Config{},
			trackChanges: true,
			want:         Config{TrackChanges: true},
		},
		{
			name:         "disable track changes",
			config:       Config{},
			trackChanges: false,
			want:         Config{TrackChanges: false},
		},
		{
			name:         "override existing track changes setting",
			config:       Config{TrackChanges: true},
			trackChanges: false,
			want:         Config{TrackChanges: false},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.config.SetTrackChanges(tt.trackChanges)
			if got != tt.want {
				t.Errorf("Config.SetTrackChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Chaining(t *testing.T) {
	t.Parallel()

	config := Config{}
	result := config.
		SetDBType(SpannerDBType).
		SetChangeTrackingTable("changes").
		SetTrackChanges(true)

	want := Config{
		DBType:              SpannerDBType,
		ChangeTrackingTable: "changes",
		TrackChanges:        true,
	}

	if result != want {
		t.Errorf("Config chaining = %v, want %v", result, want)
	}
}