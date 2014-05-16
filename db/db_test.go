package db

import (
  "testing"
)

// Test creating and deleting user
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

// Test creating and deleting resource
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

// Test creating and deleting permission
func Test_CreatePermission(t *testing.T) {
  // create user
  user := User{Name: "Sau Sheong", Email: "sausheong@me.com", Password: "123"}
  if err := DB.Save(&user).Error; err != nil {
    t.Errorf("Cannot create user, error is %v", err)
  }

  // create resource
  resource := Resource{Name: "Resource2"}
  if err := DB.Save(&resource).Error; err != nil {
    t.Errorf("Cannot create resource, error is %v", err)
  }

  // create permission
  permission := Permission{User_id: user.Uuid, Resource_id: resource.Uuid}
  if err := DB.Save(&permission).Error; err != nil {
    t.Errorf("Cannot create permission, error is %v", err)
  }

  // clean up

  if err := DB.Delete(&permission).Error; err != nil {
    t.Errorf("Cannot delete permission after creating, error is %v", err)
  }
  if err := DB.Delete(&resource).Error; err != nil {
    t.Errorf("Cannot delete resource after creating, error is %v", err)
  }
  if err := DB.Delete(&user).Error; err != nil {
    t.Errorf("Cannot delete user after creating, error is %v", err)
  }
}

func Test_ActivateUser(t *testing.T) {
  user := create_user("Sau Sheong", "sausheong@me.com", "123")
  defer delete_user(user)

  if user.Activated == true {
    t.Errorf("User shouldn't be activated")
  }
  if user.ActivationToken == "" {
    t.Errorf("User shouldn't have an empty activation token")
  }

  err := user.Activate(user.ActivationToken)
  if err != nil {
    t.Errorf("User cannot be activated because %v", err)
  }

  if user.Activated == false {
    t.Errorf("User not activated yet")
  }

}

func Test_DeactivateUser(t *testing.T) {
  user := create_user("Sau Sheong", "sausheong@me.com", "123")
  defer delete_user(user)

  err := user.Activate(user.ActivationToken)
  if err != nil {
    t.Errorf("User cannot be activated because %v", err)
  }

  if user.Activated == false {
    t.Errorf("User not activated yet")
  }

  err = user.Deactivate()
  if err != nil {
    t.Errorf("User cannot be activated because %v", err)
  }

  if user.Activated == true {
    t.Errorf("User is still activated")
  }


}


func create_user(name string, email string, password string) User {
  user := User{Name: name, Email: email, Password: password}
  err := DB.DB().Ping()
  if err != nil {
    return user
  } else {
    DB.Save(&user)
    return user
  }
}

func delete_user(user User) (err error) {
  db_err := DB.DB().Ping()
  if db_err != nil {
    return
  } else {
    err = DB.Delete(&user).Error
    return
  }
}
