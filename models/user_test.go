package models

import (
	"fmt"
	"testing"
	"time"
)

func setupUserTestData() uint {
	InitDB()
	email := fmt.Sprintf("user_%d@test.com", time.Now().UnixNano())
	user := User{Email: email}
	user.Password, _ = HashPassword("password")
	created, _ := AddUser(user)
	return created.ID
}

func TestGetUserByEmail(t *testing.T) {
	InitDB()
	email := fmt.Sprintf("getbyemail_%d@test.com", time.Now().UnixNano())
	user := User{Email: email}
	user.Password, _ = HashPassword("password")
	AddUser(user)

	// Existing email
	found, ok := GetUserByEmail(email)
	if !ok {
		t.Errorf("expected user to be found")
	}
	if found.Email != email {
		t.Errorf("expected email %s, got %s", email, found.Email)
	}

	// Non existing email
	_, ok = GetUserByEmail("nonexistent@test.com")
	if ok {
		t.Errorf("expected user not to be found")
	}
}

func TestUpdateUser(t *testing.T) {
	userID := setupUserTestData()

	updated := User{Email: fmt.Sprintf("updated_%d@test.com", time.Now().UnixNano())}
	updated.Password, _ = HashPassword("newpassword")

	result, ok := UpdateUser(userID, updated)
	if !ok {
		t.Errorf("expected update to succeed")
	}
	if result.ID != userID {
		t.Errorf("expected ID %d, got %d", userID, result.ID)
	}

	// Non existing user
	_, ok = UpdateUser(9999, updated)
	if ok {
		t.Errorf("expected update to fail for non existing user")
	}
}

func TestDeleteUser(t *testing.T) {
	userID := setupUserTestData()

	// Delete existing user
	deleted, ok := DeleteUser(userID)
	if !ok {
		t.Errorf("expected delete to succeed")
	}
	if deleted.ID != userID {
		t.Errorf("expected ID %d, got %d", userID, deleted.ID)
	}

	// Delete already deleted user
	_, ok = DeleteUser(userID)
	if ok {
		t.Errorf("expected delete to fail for already deleted user")
	}

	// Delete non existing user
	_, ok = DeleteUser(9999)
	if ok {
		t.Errorf("expected delete to fail for non existing user")
	}
}