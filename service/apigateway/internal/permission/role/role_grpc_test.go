package rolepermission

import (
	"context"
	"errors"
	"testing"

	rolepb "github.com/MamangRust/microservice-point-of-sale-pb/role"
	userrolepb "github.com/MamangRust/microservice-point-of-sale-pb/user_role"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"google.golang.org/grpc"
)

type fakeRoleQueryClient struct {
	calls    int
	lastUser int32
	res      *rolepb.ApiResponsesRole
	err      error
}

func (f *fakeRoleQueryClient) FindByUserId(
	_ context.Context,
	in *userrolepb.FindByIdUserRoleRequest,
	_ ...grpc.CallOption,
) (*rolepb.ApiResponsesRole, error) {
	f.calls++
	f.lastUser = in.GetUserId()
	return f.res, f.err
}

type fakeRoleCache struct {
	roles   []string
	found   bool
	written []string
}

func (f *fakeRoleCache) GetRoleCache(context.Context, string) ([]string, bool) {
	return f.roles, f.found
}

func (f *fakeRoleCache) SetRoleCache(_ context.Context, _ string, roles []string) {
	f.written = roles
}

func roleResponse(names ...string) *rolepb.ApiResponsesRole {
	data := make([]*rolepb.RoleResponse, 0, len(names))
	for _, name := range names {
		data = append(data, &rolepb.RoleResponse{Name: name})
	}
	return &rolepb.ApiResponsesRole{Status: "success", Data: data}
}

func TestValidateRole_UsesCacheWithoutCallingTheService(t *testing.T) {
	client := &fakeRoleQueryClient{}
	cache := &fakeRoleCache{roles: []string{"Admin"}, found: true}

	perm := NewRolePermissionGRPC(client, logger.NoopLogger{}, cache)

	roles, err := perm.ValidateRole(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected cached roles, got error %v", err)
	}
	if len(roles) != 1 || roles[0] != "Admin" {
		t.Fatalf("expected the cached roles to be returned, got %v", roles)
	}
	if client.calls != 0 {
		t.Fatalf("expected the role service not to be called, got %d calls", client.calls)
	}
}

func TestValidateRole_FetchesFromServiceAndCaches(t *testing.T) {
	client := &fakeRoleQueryClient{res: roleResponse("ROLE_ADMIN", "Cashier")}
	cache := &fakeRoleCache{}

	perm := NewRolePermissionGRPC(client, logger.NoopLogger{}, cache)

	roles, err := perm.ValidateRole(context.Background(), 42)
	if err != nil {
		t.Fatalf("expected roles from the role service, got error %v", err)
	}
	if len(roles) != 2 || roles[0] != "ROLE_ADMIN" || roles[1] != "Cashier" {
		t.Fatalf("unexpected roles %v", roles)
	}
	if client.lastUser != 42 {
		t.Fatalf("expected user 42 to be looked up, got %d", client.lastUser)
	}
	if len(cache.written) != 2 {
		t.Fatalf("expected the roles to be cached, got %v", cache.written)
	}
}

func TestValidateRole_RejectsUserWithoutRoles(t *testing.T) {
	client := &fakeRoleQueryClient{res: roleResponse()}
	cache := &fakeRoleCache{}

	perm := NewRolePermissionGRPC(client, logger.NoopLogger{}, cache)

	if _, err := perm.ValidateRole(context.Background(), 1); err == nil {
		t.Fatal("expected an error for a user without roles")
	}
	if cache.written != nil {
		t.Fatalf("expected nothing to be cached, got %v", cache.written)
	}
}

func TestValidateRole_RejectsWhenServiceFails(t *testing.T) {
	client := &fakeRoleQueryClient{err: errors.New("unavailable")}

	perm := NewRolePermissionGRPC(client, logger.NoopLogger{}, &fakeRoleCache{})

	if _, err := perm.ValidateRole(context.Background(), 1); err == nil {
		t.Fatal("expected an error when the role service fails")
	}
}

func TestCheckRole(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		required []string
		wantErr  bool
	}{
		{name: "matching role", userRole: "Admin", required: []string{"Admin", "ROLE_ADMIN"}},
		{name: "matching alternate role", userRole: "ROLE_ADMIN", required: []string{"Admin", "ROLE_ADMIN"}},
		{name: "no matching role", userRole: "Cashier", required: []string{"Admin", "ROLE_ADMIN"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fakeRoleQueryClient{res: roleResponse(tt.userRole)}
			perm := NewRolePermissionGRPC(client, logger.NoopLogger{}, &fakeRoleCache{})

			err := perm.CheckRole(context.Background(), 1, tt.required...)
			if tt.wantErr && err == nil {
				t.Fatal("expected the role check to fail")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected the role check to pass, got %v", err)
			}
		})
	}
}

func TestValidateRole_WithoutClientFailsInsteadOfPanicking(t *testing.T) {
	perm := NewRolePermissionGRPC(nil, logger.NoopLogger{}, &fakeRoleCache{})

	if _, err := perm.ValidateRole(context.Background(), 1); err == nil {
		t.Fatal("expected an error when no role client is configured")
	}
}
