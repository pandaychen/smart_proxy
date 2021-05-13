package enums

type LB_TYPE string

var (
	LB_WEIGHT_RR       LB_TYPE = "wrr"
	LB_CONSISTENT_HASH LB_TYPE = "ch"
	LB_P2C             LB_TYPE = "p2c"
)
