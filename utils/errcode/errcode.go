package errcode

const (
	SUCCESS      = "SUCCESS"
	ERROR        = "ERROR"
	ERROR_PARAMS = "ERROR_PARAMS"
	ERROR_TOKEN  = "ERROR_TOKEN"
)

var MsgFlags = map[string]string{
	"SUCCESS":      "ok",
	"ERROR":        "failed",
	"ERROR_TOKEN":  "invalid user",
	"ERROR_PARAMS": "invalid params",
}

func GetMsg(name string) string {
	msg, ok := MsgFlags[name]
	if ok {
		return msg
	}
	return MsgFlags["ERROR"]
}
