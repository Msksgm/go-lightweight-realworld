package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewUser(t *testing.T) {
	got, err := NewUser(1, "example@email.com", "userName", "password", time.Now(), time.Now())
	if err != nil {
		t.Error(err)
	}
	want := &User{ID: 1, Email: "example@email.com", UserName: "userName"}
	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(*got, "PasswordHash", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}
