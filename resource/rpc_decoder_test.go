package resource

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cccteam/ccc/accesstypes"
	"github.com/cccteam/httpio"
)

type mockUserPermissions struct {
	user    accesstypes.User
	domain  accesstypes.Domain
	checkFunc func(context.Context, accesstypes.Permission, ...accesstypes.Resource) (bool, []accesstypes.Resource, error)
}

func (m *mockUserPermissions) User() accesstypes.User {
	return m.user
}

func (m *mockUserPermissions) Domain() accesstypes.Domain {
	return m.domain
}

func (m *mockUserPermissions) Check(ctx context.Context, perm accesstypes.Permission, resources ...accesstypes.Resource) (bool, []accesstypes.Resource, error) {
	if m.checkFunc != nil {
		return m.checkFunc(ctx, perm, resources...)
	}
	return true, nil, nil
}

type testRequest struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type mockValidator struct {
	structFunc       func(interface{}) error
	structPartialFunc func(interface{}, ...string) error
}

func (m *mockValidator) Struct(s interface{}) error {
	if m.structFunc != nil {
		return m.structFunc(s)
	}
	return nil
}

func (m *mockValidator) StructPartial(s interface{}, fields ...string) error {
	if m.structPartialFunc != nil {
		return m.structPartialFunc(s, fields...)
	}
	return nil
}

func TestNewRPCDecoder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		userPermissions   func(*http.Request) UserPermissions
		methodName        accesstypes.Resource
		perm              accesstypes.Permission
		wantErr           bool
	}{
		{
			name: "valid RPC decoder",
			userPermissions: func(*http.Request) UserPermissions {
				return &mockUserPermissions{
					user:   accesstypes.User("testuser"),
					domain: accesstypes.Domain("testdomain"),
				}
			},
			methodName: "test-method",
			perm:       accesstypes.Execute,
			wantErr:    false,
		},
		{
			name: "invalid struct decoder creation",
			userPermissions: func(*http.Request) UserPermissions {
				return &mockUserPermissions{
					user:   accesstypes.User("testuser"),
					domain: accesstypes.Domain("testdomain"),
				}
			},
			methodName: "test-method",
			perm:       accesstypes.Execute,
			wantErr:    false, // NewStructDecoder should work for testRequest
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewRPCDecoder[testRequest](tt.userPermissions, tt.methodName, tt.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRPCDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewRPCDecoder() returned nil")
			}
		})
	}
}

func TestRPCDecoder_WithValidator(t *testing.T) {
	t.Parallel()

	userPermissions := func(*http.Request) UserPermissions {
		return &mockUserPermissions{
			user:   "testuser",
			domain: "testdomain",
		}
	}

	decoder, err := NewRPCDecoder[testRequest](userPermissions, "test-method", accesstypes.Execute)
	if err != nil {
		t.Fatalf("NewRPCDecoder() error = %v", err)
	}

	validator := &mockValidator{
		structFunc: func(s interface{}) error {
			if req, ok := s.(*testRequest); ok && req.Name == "" {
				return httpio.NewBadRequestMessage("name is required")
			}
			return nil
		},
		structPartialFunc: func(s interface{}, fields ...string) error {
			return nil
		},
	}

	got := decoder.WithValidator(validator)
	if got == nil {
		t.Errorf("RPCDecoder.WithValidator() returned nil")
	}
	if got == decoder {
		t.Errorf("RPCDecoder.WithValidator() returned same instance")
	}
}

func TestRPCDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		decoder         *RPCDecoder[testRequest]
		request         *http.Request
		wantErr         bool
		wantForbidden   bool
	}{
		{
			name: "successful decode",
			decoder: func() *RPCDecoder[testRequest] {
				decoder, _ := NewStructDecoder[testRequest]()
				return &RPCDecoder[testRequest]{
					d: decoder,
					res: "test-method",
					requiredPermission: accesstypes.Execute,
					userPermissions: func(*http.Request) UserPermissions {
						return &mockUserPermissions{
							user:   accesstypes.User("testuser"),
							domain: accesstypes.Domain("testdomain"),
							checkFunc: func(ctx context.Context, perm accesstypes.Permission, resources ...accesstypes.Resource) (bool, []accesstypes.Resource, error) {
								return true, nil, nil
							},
						}
					},
				}
			}(),
			request: &http.Request{
				Body:   io.NopCloser(strings.NewReader(`{"name":"test","value":123}`)),
				Header: make(http.Header),
			},
			wantErr:       false,
			wantForbidden: false,
		},
		{
			name: "permission denied",
			decoder: func() *RPCDecoder[testRequest] {
				decoder, _ := NewStructDecoder[testRequest]()
				return &RPCDecoder[testRequest]{
					d: decoder,
					res: "test-method",
					requiredPermission: accesstypes.Execute,
					userPermissions: func(*http.Request) UserPermissions {
						return &mockUserPermissions{
							user:   accesstypes.User("testuser"),
							domain: accesstypes.Domain("testdomain"),
							checkFunc: func(ctx context.Context, perm accesstypes.Permission, resources ...accesstypes.Resource) (bool, []accesstypes.Resource, error) {
								return false, []accesstypes.Resource{"missing-resource"}, nil
							},
						}
					},
				}
			}(),
			request: &http.Request{
				Body:   io.NopCloser(strings.NewReader(`{"name":"test","value":123}`)),
				Header: make(http.Header),
			},
			wantErr:       true,
			wantForbidden: true,
		},
		{
			name: "permission check error",
			decoder: func() *RPCDecoder[testRequest] {
				decoder, _ := NewStructDecoder[testRequest]()
				return &RPCDecoder[testRequest]{
					d: decoder,
					res: "test-method",
					requiredPermission: accesstypes.Execute,
					userPermissions: func(*http.Request) UserPermissions {
						return &mockUserPermissions{
							user:   accesstypes.User("testuser"),
							domain: accesstypes.Domain("testdomain"),
							checkFunc: func(ctx context.Context, perm accesstypes.Permission, resources ...accesstypes.Resource) (bool, []accesstypes.Resource, error) {
								return false, nil, httpio.NewInternalServerErrorMessage("permission check failed")
							},
						}
					},
				}
			}(),
			request: &http.Request{
				Body:   io.NopCloser(strings.NewReader(`{"name":"test","value":123}`)),
				Header: make(http.Header),
			},
			wantErr:       true,
			wantForbidden: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.decoder.Decode(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RPCDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.wantForbidden {
					if !httpio.HasForbidden(err) {
						t.Errorf("RPCDecoder.Decode() expected forbidden error, got %T: %v", err, err)
					}
				}
				return
			}
			if got == nil {
				t.Errorf("RPCDecoder.Decode() returned nil result")
			}
		})
	}
}

func TestRPCDecoder_Decode_StructDecoderError(t *testing.T) {
	t.Parallel()

	decoder := &RPCDecoder[testRequest]{
		d: &StructDecoder[testRequest]{},
		res: "test-method",
		requiredPermission: accesstypes.Execute,
		userPermissions: func(*http.Request) UserPermissions {
			return &mockUserPermissions{
				user:   accesstypes.User("testuser"),
				domain: accesstypes.Domain("testdomain"),
			}
		},
	}

	// Create a request with invalid JSON body
	request := &http.Request{
		Body: http.NoBody, // This will cause a decode error
	}

	_, err := decoder.Decode(request)
	if err == nil {
		t.Errorf("RPCDecoder.Decode() expected error, got nil")
	}
}