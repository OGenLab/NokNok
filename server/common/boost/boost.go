package boost

import (
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_boost"
)

var BoostSet = map[int]*t_boost.T_boost{}

func init() {
	var (
		m = model.NewModel()
	)
	defer m.Close()

	boots, err := t_boost.Select(m.DB, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	for i := range boots {
		BoostSet[boots[i].F_level] = boots[i]
	}
}
