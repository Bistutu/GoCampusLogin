package entity

// Result 结果
type Result struct {
	Code int // 1 表示成功，-1 表示失败
	Msg  status
	Data interface{}
}

type status string

const (
	SUCCESS status = "success"
	FAIL    status = "fail"
)
