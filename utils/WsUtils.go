package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"strings"
)

func WriteWs(ws *websocket.Conn, bus string, data any) bool {
	if data == nil {
		return false
	}
	var out_data string
	switch value := data.(type) {
	case string:
		out_data = string(bus) + ":" + value
		break
	default:
		if o, err := json.Marshal(data); err != nil {
			return false
		} else {
			out_data = bus + ":" + string(o)
		}

		break
	}

	return ws.WriteMessage(1, []byte(out_data)) == nil

}

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
		outData = opCode
		return []byte(outData)
	}
	switch array.(type) {
	case []byte:
		v, _ := array.([]byte)
		outData = opCode + ":" + subOpCode + ":" + string(v)
		break
	default:
		tmp, _ := json.Marshal(array)
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
