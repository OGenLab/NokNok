package protocol

type AccountInfo struct {
	IsPremium           bool   `json:"isPremium"`
	Boost               int    `json:"boost"`
	NextBoost           int    `json:"nextBoost"`
	Coins               uint64 `json:"coins"`
	InvitedCount        int    `json:"invitedCount"`
	InvitedPremiumCount int    `json:"invitedPremiumCount"`
}

type DrawInfo struct {
	PrizeType    int    `json:"prizeType"`
	PrizeId      uint64 `json:"prizeId"`
	Count        int    `json:"count"`
	Threshold    int    `json:"threshold"`
	LastDrawTime int64  `json:"lastDrawTime"`
}

type NokInfo struct {
	Stamina     int   `json:"stamina"`
	HitCount    int   `json:"hitCount"`
	LastHitTime int64 `json:"lastHitTime"`
}

type HammerInfo struct {
	EntitiesInfo
	EquipmentCount int `json:"equipmentCount"`
	Count          int `json:"count"`
	Quality        int `json:"quality"`
}

type MetaInfo struct {
}

type BoostInfo struct {
	Boost         int    `json:"boost"`
	NeededCoins   uint64 `json:"neededCoins"`
	ConsumeSamina int    `json:"consumeSamina"`
	CoinsRate     int    `json:"coinsRate"`
}

type TaskInfo struct {
	Id        uint64 `json:"id"`
	Season    int    `json:"season"`
	Type      int    `json:"type"`
	Threshold int    `json:"threshold"`
	Schedule  int    `json:"schedule"`
	Status    int    `json:"status"`
	Reward    uint64 `json:"reward"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type EntitiesInfo struct {
	Id          uint64 `json:"id"`
	Season      int    `json:"season"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	Probability int    `json:"probability"`
}

type LeaderboardInfo struct {
	Rank  int    `json:"rank"`
	Id    uint64 `json:"id"`
	Count uint64 `json:"count"`
}
