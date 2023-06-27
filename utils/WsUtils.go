package utils

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/glog"
	"strings"
)

// EnCode 编码
func EnCode(opCode string, array []byte) []byte {
	var outData string
	if array == nil {
		outData = opCode
	} else {
		outData = opCode + ":" + string(array)
	}
	return []byte(outData)
}

// EnSubCode 编码
func EnSubCode(opCode string, subOpCode string, array any) []byte {
	var outData string
	if array == nil {
		outData = opCode + ":" + subOpCode
		return []byte(outData)
	}
	switch array.(type) {
	case []byte:
		v, _ := array.([]byte)
		outData = opCode + ":" + subOpCode + ":" + string(v)
		break
	case string:
		v, _ := array.(string)
		outData = opCode + ":" + subOpCode + ":" + v
		break
	default:
		tmp, err := gjson.Encode(array)
		if err != nil {
			glog.Info(context.Background(), "序列化错误", err.Error())
		}
		outData = opCode + ":" + subOpCode + ":" + string(tmp)
		break

	}
	return []byte(outData)
}

func SplitMsg(array []byte) (mType string, outArray []byte) {

	idx := strings.Index(string(array), ":")
	if idx < 1 {
		return "-1", nil
	}
	mType = string(array[0:idx])
	outArray = array[idx+1:]
	return mType, outArray
}
