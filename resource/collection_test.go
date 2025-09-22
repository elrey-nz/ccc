package resource

import (
	"testing"

	"github.com/cccteam/ccc/accesstypes"
)

func TestNewCollection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *Collection
	}{
		{
			name: "new collection",
			want: &Collection{}, // When collectResourcePermissions is false, it returns empty Collection
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCollection()
			if got == nil {
				t.Errorf("NewCollection() returned nil")
				return
			}
			// When collectResourcePermissions is false, all fields are nil
			// This is the expected behavior
		})
	}
}

func TestCollection_AddResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		c          *Collection
		scope      accesstypes.PermissionScope
		permission accesstypes.Permission
		res        accesstypes.Resource
		wantErr    bool
	}{
		{
			name: "add valid resource",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore, 2),
			},
			scope:      "test-scope",
			permission: accesstypes.Create,
			res:        "test-resource",
			wantErr:    false,
		},
		{
			name: "add null permission",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore, 2),
			},
			scope:      "test-scope",
			permission: accesstypes.NullPermission,
			res:        "test-resource",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.c.AddResource(tt.scope, tt.permission, tt.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.AddResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_AddMethodResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		c          *Collection
		scope      accesstypes.PermissionScope
		permission accesstypes.Permission
		res        accesstypes.Resource
		wantErr    bool
	}{
		{
			name: "add valid method resource",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore, 2),
			},
			scope:      "test-scope",
			permission: accesstypes.Execute,
			res:        "test-method",
			wantErr:    false,
		},
		{
			name: "add null permission",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore, 2),
			},
			scope:      "test-scope",
			permission: accesstypes.NullPermission,
			res:        "test-method",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.c.AddMethodResource(tt.scope, tt.permission, tt.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.AddMethodResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_IsResourceImmutable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		scope accesstypes.PermissionScope
		res   accesstypes.Resource
		want bool
	}{
		{
			name: "immutable resource",
			c: &Collection{
				immutableFields: map[accesstypes.PermissionScope]immutableFieldMap{
					"test-scope": {
						"test-resource": {
							"test-tag": struct{}{},
						},
					},
				},
			},
			scope: "test-scope",
			res:   "test-resource.test-tag",
			want:  true,
		},
		{
			name: "mutable resource",
			c: &Collection{
				immutableFields: map[accesstypes.PermissionScope]immutableFieldMap{
					"test-scope": {
						"test-resource": {
							"other-tag": struct{}{},
						},
					},
				},
			},
			scope: "test-scope",
			res:   "test-resource.test-tag",
			want:  false,
		},
		{
			name: "non-existent scope",
			c: &Collection{
				immutableFields: make(map[accesstypes.PermissionScope]immutableFieldMap),
			},
			scope: "non-existent-scope",
			res:   "test-resource.test-tag",
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.IsResourceImmutable(tt.scope, tt.res)
			if got != tt.want {
				t.Errorf("Collection.IsResourceImmutable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_permissions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		want []accesstypes.Permission
	}{
		{
			name: "empty collection",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore),
				tagStore:      make(map[accesstypes.PermissionScope]tagStore),
			},
			want: []accesstypes.Permission{},
		},
		{
			name: "collection with permissions",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Create, accesstypes.Read},
					},
				},
				tagStore: map[accesstypes.PermissionScope]tagStore{
					"scope2": {
						"resource2": {
							"tag1": {accesstypes.Update},
						},
					},
				},
			},
			want: []accesstypes.Permission{accesstypes.Create, accesstypes.Read, accesstypes.Update},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.permissions()
			if len(got) != len(tt.want) {
				t.Errorf("Collection.permissions() length = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestCollection_resourcePermissions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		want []accesstypes.Permission
	}{
		{
			name: "empty collection",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore),
				tagStore:      make(map[accesstypes.PermissionScope]tagStore),
			},
			want: []accesstypes.Permission{},
		},
		{
			name: "collection with execute permission filtered out",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Create, accesstypes.Read, accesstypes.Execute},
					},
				},
			},
			want: []accesstypes.Permission{accesstypes.Create, accesstypes.Read},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.resourcePermissions()
			if len(got) != len(tt.want) {
				t.Errorf("Collection.resourcePermissions() length = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestCollection_Resources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		want []accesstypes.Resource
	}{
		{
			name: "empty collection",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore),
			},
			want: []accesstypes.Resource{},
		},
		{
			name: "collection with resources",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Create, accesstypes.Read},
						"resource2": {accesstypes.Execute}, // Should be filtered out
					},
				},
			},
			want: []accesstypes.Resource{"resource1"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.Resources()
			if len(got) != len(tt.want) {
				t.Errorf("Collection.Resources() length = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestCollection_ResourceExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		r    accesstypes.Resource
		want bool
	}{
		{
			name: "existing resource",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Create, accesstypes.Read},
					},
				},
			},
			r:    "resource1",
			want: true,
		},
		{
			name: "non-existing resource",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Create, accesstypes.Read},
					},
				},
			},
			r:    "resource2",
			want: false,
		},
		{
			name: "resource with only execute permission",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Execute},
					},
				},
			},
			r:    "resource1",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.ResourceExists(tt.r)
			if got != tt.want {
				t.Errorf("Collection.ResourceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_tags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		want map[accesstypes.Resource][]accesstypes.Tag
	}{
		{
			name: "empty collection",
			c: &Collection{
				tagStore: make(map[accesstypes.PermissionScope]tagStore),
			},
			want: map[accesstypes.Resource][]accesstypes.Tag{},
		},
		{
			name: "collection with tags",
			c: &Collection{
				tagStore: map[accesstypes.PermissionScope]tagStore{
					"scope1": {
						"resource1": {
							"tag1": {accesstypes.Create},
							"tag2": {accesstypes.Read},
						},
					},
				},
			},
			want: map[accesstypes.Resource][]accesstypes.Tag{
				"resource1": {"tag1", "tag2"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.tags()
			if len(got) != len(tt.want) {
				t.Errorf("Collection.tags() length = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestCollection_domains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		want []accesstypes.PermissionScope
	}{
		{
			name: "empty collection",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore),
			},
			want: []accesstypes.PermissionScope{},
		},
		{
			name: "collection with domains",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {},
					"scope2": {},
				},
			},
			want: []accesstypes.PermissionScope{"scope1", "scope2"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.domains()
			if len(got) != len(tt.want) {
				t.Errorf("Collection.domains() length = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestCollection_Scope(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		resource accesstypes.Resource
		want accesstypes.PermissionScope
	}{
		{
			name: "resource in resource store",
			c: &Collection{
				resourceStore: map[accesstypes.PermissionScope]resourceStore{
					"scope1": {
						"resource1": {accesstypes.Create},
					},
				},
			},
			resource: "resource1",
			want:     "scope1",
		},
		{
			name: "resource in tag store",
			c: &Collection{
				tagStore: map[accesstypes.PermissionScope]tagStore{
					"scope1": {
						"resource1": {
							"tag1": {accesstypes.Create},
						},
					},
				},
			},
			resource: "resource1.tag1",
			want:     "scope1",
		},
		{
			name: "non-existent resource",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore),
				tagStore:      make(map[accesstypes.PermissionScope]tagStore),
			},
			resource: "non-existent",
			want:     "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.Scope(tt.resource)
			if got != tt.want {
				t.Errorf("Collection.Scope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_TypescriptData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    *Collection
		want TypescriptData
	}{
		{
			name: "empty collection",
			c: &Collection{
				resourceStore: make(map[accesstypes.PermissionScope]resourceStore),
				tagStore:      make(map[accesstypes.PermissionScope]tagStore),
			},
			want: TypescriptData{
				Permissions:           []accesstypes.Permission{},
				ResourcePermissions:   []accesstypes.Permission{},
				Resources:             []accesstypes.Resource{},
				ResourceTags:          map[accesstypes.Resource][]accesstypes.Tag{},
				ResourcePermissionMap: map[accesstypes.Resource]map[accesstypes.Permission]bool{},
				Domains:               []accesstypes.PermissionScope{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.TypescriptData()
			if got.Permissions == nil {
				t.Errorf("Collection.TypescriptData() Permissions is nil")
			}
			if got.ResourcePermissions == nil {
				t.Errorf("Collection.TypescriptData() ResourcePermissions is nil")
			}
			if got.Resources == nil {
				t.Errorf("Collection.TypescriptData() Resources is nil")
			}
			if got.ResourceTags == nil {
				t.Errorf("Collection.TypescriptData() ResourceTags is nil")
			}
			if got.ResourcePermissionMap == nil {
				t.Errorf("Collection.TypescriptData() ResourcePermissionMap is nil")
			}
			if got.Domains == nil {
				t.Errorf("Collection.TypescriptData() Domains is nil")
			}
		})
	}
}