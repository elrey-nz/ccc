package resource

import (
	"context"
	"testing"
)

func TestUserEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ctx      context.Context
		wantPanic bool
	}{
		{
			name: "empty context causes panic",
			ctx:  context.Background(),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("UserEvent() panicked unexpectedly: %v", r)
					}
				} else if tt.wantPanic {
					t.Errorf("UserEvent() should have panicked but didn't")
				}
			}()
			UserEvent(tt.ctx)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		processName string
		want        string
	}{
		{
			name:        "valid process name",
			processName: "test-process",
			want:        "Process test-process",
		},
		{
			name:        "empty process name",
			processName: "",
			want:        "Process ",
		},
		{
			name:        "special characters in process name",
			processName: "process-with-dashes_and_underscores",
			want:        "Process process-with-dashes_and_underscores",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ProcessEvent(tt.processName)
			if got != tt.want {
				t.Errorf("ProcessEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserProcessEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ctx         context.Context
		processName string
		wantPanic   bool
	}{
		{
			name:        "empty context causes panic",
			ctx:         context.Background(),
			processName: "",
			wantPanic:   true,
		},
		{
			name:        "empty context with process causes panic",
			ctx:         context.Background(),
			processName: "test-process",
			wantPanic:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("UserProcessEvent() panicked unexpectedly: %v", r)
					}
				} else if tt.wantPanic {
					t.Errorf("UserProcessEvent() should have panicked but didn't")
				}
			}()
			UserProcessEvent(tt.ctx, tt.processName)
		})
	}
}