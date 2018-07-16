package m

import (
	//"fmt"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken(t *testing.T) {
	// Load default config
	serve.Conf = serve.DefaultConfig()

	// migrate db
	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	u := &User{
		ID: 1,
	}

	token, err := NewToken(TypeActivate, u)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, token.UserID, u.ID)
	assert.Equal(t, token.Type, TypeActivate)
	assert.Equal(t, token.Status, StatusValid)

	if err := token.UsedBy(u); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, token.Status, StatusInvalid)

	err = token.UsedBy(u)
	assert.NotEqual(t, err, nil)
}
