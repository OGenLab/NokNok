package entities

const (
	WOOD_HAMMER_ID = iota + 1
	STELL_HAMMER_ID
	SILVER_HAMMER_ID
	GOLD_HAMMER_ID
	DIAMOND_HAMMER_ID
	THOR_HAMMER_ID
)

var TypExtToHammerQuality = map[string]int{
	"1": WOOD_HAMMER_ID,
	"2": STELL_HAMMER_ID,
	"3": SILVER_HAMMER_ID,
	"4": GOLD_HAMMER_ID,
	"5": DIAMOND_HAMMER_ID,
	"6": THOR_HAMMER_ID,
}
