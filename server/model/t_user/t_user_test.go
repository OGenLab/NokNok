package t_user

import (
	"testing"

	"lightspeed.2dao3.com/nokserver/model"
)

func TestCreateUser(t *testing.T) {
	var (
		m = model.NewModelDefault()
	)
	defer m.Close()

	user := T_user{
		F_uid: 1,
	}

	if err := Create(m.DB, &user); err != nil {
		t.Fatal(err)
	}
}
