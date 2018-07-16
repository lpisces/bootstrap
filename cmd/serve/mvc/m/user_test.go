package m

import (
	"fmt"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreate(t *testing.T) {

	// Load default config
	serve.Conf = serve.DefaultConfig()

	// migrate db
	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	u := &User{
		Email:           "iamalazyrat@gmail.com",
		Password:        "hellobootstrap",
		PasswordConfirm: "hellobootstrap",
	}

	exist, err := u.Exist()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exist, false)

	defer func() {
		db, err := GetDB()
		if err != nil {
			t.Fatal(err)
		}
		db.LogMode(true)
		if db.Where("email = ?", "iamalazyrat@gmail.com").First(u).RecordNotFound() {
			t.Log(u.ID)
			t.Fatal(fmt.Errorf("write db failed"))
		}
		db.Delete(u)
	}()

	if err := u.Create(); err != nil {
		t.Fatal(err)
	}

	exist, err = u.Exist()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exist, true)

}
