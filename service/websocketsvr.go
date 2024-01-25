package service

import (
	"lufergo/uts"
	"time"

	"github.com/gorilla/websocket"
)

func websockerSvr(c *websocket.Conn, b []byte) []byte {
	var bt int64
	// if IsDebug {
	bt = time.Now().UnixMilli()
	// }

	// clientTrans := model.SendToServer(b, c.RemoteAddr().String())

	// _result := sendToWebsocketClient(clientTrans)

	// if IsDebug {
	var allT = time.Now().UnixMilli() - bt
	if allT > 1000 {
		uts.Log("pay time websocket 1000", allT)
	} else if allT > 500 {
		uts.Log("pay time websocket 500", allT)
	} else if allT > 100 {
		uts.Log("pay time websocket 100", allT)
	} else if allT > 50 {
		uts.Log("pay time websocket 50", allT)
	}
	// }

	return nil
}

// func sendToWebsocketClient(clientTrans *pbstruct.ClientTrans) []byte {
// 	if clientTrans.Id < 0 {
// 		return nil
// 	} else {
// 		if IsDebug {
// 			// clientTrans.Key = key
// 			_data := pbuts.DeCT2CPDat(clientTrans)
// 			_, _name := pbuts.GetProtoIDNameFromCP(pbuts.GetCPFromProtoID(clientTrans.Id))
// 			LogDebug("返回", _name, clientTrans.Token, clientTrans.Err, clientTrans.Msg, len(clientTrans.Protobuff), _data)
// 		}
// 		b, err := proto.Marshal(clientTrans)
// 		if !ChkErr(err) {
// 			return b
// 		}
// 	}
// 	return nil
// }
