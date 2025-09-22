package resource

import (
	"context"
	"errors"
	"testing"
)

type mockTxnFuncRunner struct {
	executeFunc func(context.Context, func(context.Context, TxnBuffer) error) error
}

func (m *mockTxnFuncRunner) ExecuteFunc(ctx context.Context, fn func(context.Context, TxnBuffer) error) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, fn)
	}
	return nil
}

type mockSpannerBuffer struct {
	spannerBufferFunc func(context.Context, TxnBuffer, ...string) error
}

func (m *mockSpannerBuffer) SpannerBuffer(ctx context.Context, txn TxnBuffer, eventSource ...string) error {
	if m.spannerBufferFunc != nil {
		return m.spannerBufferFunc(ctx, txn, eventSource...)
	}
	return nil
}

func TestNewCommitBuffer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		db             TxnFuncRunner
		eventSource    string
		autoCommitSize int
		want           *CommitBuffer
	}{
		{
			name:           "basic commit buffer",
			db:             &mockTxnFuncRunner{},
			eventSource:    "test-source",
			autoCommitSize: 10,
			want: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 10,
				buffer:         make([]SpannerBuffer, 0, 10),
			},
		},
		{
			name:           "zero auto commit size",
			db:             &mockTxnFuncRunner{},
			eventSource:    "test-source",
			autoCommitSize: 0,
			want: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 0,
				buffer:         make([]SpannerBuffer, 0, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCommitBuffer(tt.db, tt.eventSource, tt.autoCommitSize)
			// Can't compare function pointers directly, just check they're not nil
			if got.db == nil {
				t.Errorf("NewCommitBuffer() db is nil")
			}
			if got.eventSource != tt.want.eventSource {
				t.Errorf("NewCommitBuffer() eventSource = %v, want %v", got.eventSource, tt.want.eventSource)
			}
			if got.autoCommitSize != tt.want.autoCommitSize {
				t.Errorf("NewCommitBuffer() autoCommitSize = %v, want %v", got.autoCommitSize, tt.want.autoCommitSize)
			}
			if cap(got.buffer) != cap(tt.want.buffer) {
				t.Errorf("NewCommitBuffer() buffer capacity = %v, want %v", cap(got.buffer), cap(tt.want.buffer))
			}
		})
	}
}

func TestCommitBuffer_Buffer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		cb             *CommitBuffer
		ps             []SpannerBuffer
		wantErr        bool
		wantBufferLen  int
		expectCommit   bool
	}{
		{
			name: "buffer without auto commit",
			cb: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 0,
				buffer:         make([]SpannerBuffer, 0, 5),
			},
			ps: []SpannerBuffer{
				&mockSpannerBuffer{},
			},
			wantErr:       false,
			wantBufferLen: 1,
			expectCommit:  false,
		},
		{
			name: "buffer with auto commit triggered",
			cb: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 2,
				buffer:         make([]SpannerBuffer, 0, 2),
			},
			ps: []SpannerBuffer{
				&mockSpannerBuffer{},
				&mockSpannerBuffer{},
			},
			wantErr:       false,
			wantBufferLen: 0, // Should be committed and cleared
			expectCommit:  true,
		},
		{
			name: "buffer with auto commit not triggered",
			cb: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 5,
				buffer:         make([]SpannerBuffer, 0, 5),
			},
			ps: []SpannerBuffer{
				&mockSpannerBuffer{},
			},
			wantErr:       false,
			wantBufferLen: 1,
			expectCommit:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commitCalled := false
			tt.cb.db = &mockTxnFuncRunner{
				executeFunc: func(ctx context.Context, fn func(context.Context, TxnBuffer) error) error {
					commitCalled = true
					return fn(ctx, nil)
				},
			}

			err := tt.cb.Buffer(context.Background(), tt.ps...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommitBuffer.Buffer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.cb.buffer) != tt.wantBufferLen {
				t.Errorf("CommitBuffer.Buffer() buffer length = %v, want %v", len(tt.cb.buffer), tt.wantBufferLen)
			}
			if commitCalled != tt.expectCommit {
				t.Errorf("CommitBuffer.Buffer() commit called = %v, want %v", commitCalled, tt.expectCommit)
			}
		})
	}
}

func TestCommitBuffer_Buffer_Error(t *testing.T) {
	t.Parallel()

	cb := &CommitBuffer{
		db: &mockTxnFuncRunner{
			executeFunc: func(ctx context.Context, fn func(context.Context, TxnBuffer) error) error {
				return errors.New("commit error")
			},
		},
		eventSource:    "test-source",
		autoCommitSize: 1,
		buffer:         make([]SpannerBuffer, 0, 1),
	}

	err := cb.Buffer(context.Background(), &mockSpannerBuffer{})
	if err == nil {
		t.Errorf("CommitBuffer.Buffer() expected error, got nil")
	}
}

func TestCommitBuffer_Commit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cb      *CommitBuffer
		wantErr bool
	}{
		{
			name: "commit with empty buffer",
			cb: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 0,
				buffer:         make([]SpannerBuffer, 0, 5),
			},
			wantErr: false,
		},
		{
			name: "commit with buffered items",
			cb: &CommitBuffer{
				db:             &mockTxnFuncRunner{},
				eventSource:    "test-source",
				autoCommitSize: 0,
				buffer: []SpannerBuffer{
					&mockSpannerBuffer{},
					&mockSpannerBuffer{},
				},
			},
			wantErr: false,
		},
		{
			name: "commit with error",
			cb: &CommitBuffer{
				db: &mockTxnFuncRunner{
					executeFunc: func(ctx context.Context, fn func(context.Context, TxnBuffer) error) error {
						return errors.New("commit error")
					},
				},
				eventSource:    "test-source",
				autoCommitSize: 0,
				buffer: []SpannerBuffer{
					&mockSpannerBuffer{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			initialBufferLen := len(tt.cb.buffer)
			err := tt.cb.Commit(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("CommitBuffer.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && initialBufferLen > 0 {
				if len(tt.cb.buffer) != 0 {
					t.Errorf("CommitBuffer.Commit() buffer should be cleared after successful commit")
				}
			}
		})
	}
}

func TestCommitBuffer_Commit_SpannerBufferError(t *testing.T) {
	t.Parallel()

	cb := &CommitBuffer{
		db: &mockTxnFuncRunner{
			executeFunc: func(ctx context.Context, fn func(context.Context, TxnBuffer) error) error {
				return fn(ctx, nil)
			},
		},
		eventSource:    "test-source",
		autoCommitSize: 0,
		buffer: []SpannerBuffer{
			&mockSpannerBuffer{
				spannerBufferFunc: func(ctx context.Context, txn TxnBuffer, eventSource ...string) error {
					return errors.New("spanner buffer error")
				},
			},
		},
	}

	err := cb.Commit(context.Background())
	if err == nil {
		t.Errorf("CommitBuffer.Commit() expected error, got nil")
	}
}