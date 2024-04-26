package depositProvider

type MoneyMerchantExternalPayModel struct {
	Uid       int64 `json:"uid"`
	TargetId  int64 `json:"target_id"`
	Num       int64 `json:"num"`
	EventTime int64 `json:"event_time"`
}

type ResSupermeet struct {
	DmError  int         `json:"dm_error"`
	ErrorMsg string      `json:"error_msg"`
	Data     interface{} `json:"data"`
}

type ReqSupermeet struct {
	Message []byte `json:"message"`
	Uid     int64  `json:"uid"`
}
