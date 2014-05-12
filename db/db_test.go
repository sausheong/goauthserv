package db

import (
	"testing"
)

func Test_CreateUser(t *testing.T) {
	user := User{Name: "Sau Sheong", Email: "sausheong@me.com", Password: "123"}

	if err := DB.Save(&user).Error; err != nil {
		t.Errorf("Cannot create user, error is %v", err)
	}
	if user.Uuid == "" {
		t.Errorf("User has no UUID")
	}
	if user.Password == "" {
		t.Errorf("User has no Password")
	}
	if user.Password == "123" {
		t.Errorf("User password has not been hashed")
	}

	if err := DB.Delete(&user).Error; err != nil {
		t.Errorf("Cannot delete user after creating, error is %v", err)
	}
}

func Test_CreateResource(t *testing.T) {
	resource := Resource{Name: "Resource1"}

	if err := DB.Save(&resource).Error; err != nil {
		t.Errorf("Cannot create resource, error is %v", err)
	}

	if resource.Name != "Resource1" {
		t.Errorf("Resource name is wrong")
	}

	if err := DB.Delete(&resource).Error; err != nil {
		t.Errorf("Cannot delete resource after creating, error is %v", err)
	}

}
