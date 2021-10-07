package enums

type BACKEND_STAT string //每次请求记录后端节点的状态

var (
	BACKEND_GOOD  BACKEND_STAT = "succ"
	BACKEND_BAD   BACKEND_STAT = "fail"
	BACKEND_TOTAL BACKEND_STAT = "total"
)
