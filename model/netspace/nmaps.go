package netspace

import (
	"bytes"
	"encoding/json"
	"lufergo/model/db"
	"lufergo/uts"
	"net/http"
	"strconv"
	"strings"
)

func nmapsSearch(dat *SearchDat) *NmapsResult {
	requestStr := db.Conf.String("nmaps_api")

	var sqlText = strings.ToLower(dat.Sql_text)
	// sqlText = strings.ReplaceAll(sqlText, "ip=", "host=")
	Log("Search", sqlText)

	param := map[string]string{
		"sql_text":           sqlText,
		"page_number":        strconv.FormatInt(int64(dat.Page_number), 10),
		"per_page":           strconv.FormatInt(int64(dat.Per_page), 10),
		"exclude_empty_data": strconv.FormatBool(dat.Exclude_empty_data),
	}

	bytesData, err := json.Marshal(param)
	if uts.ChkErr(err) {
		return nil
	}
	req, err := http.NewRequest("POST", requestStr, bytes.NewBuffer(bytesData))
	if uts.ChkErr(err) {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Log(param)
	respBytes := callNetSpaceAPI(req, "")
	if respBytes == nil {
		return nil
	}
	// Log(respBytes.String())
	_result := &NmapsResult{}
	json.Unmarshal(respBytes.Bytes(), _result)
	if _result.Status == "error" {
		Log("调用错误", _result)
		return nil
	}
	saveNmapsToRedis(_result)

	return _result
}

func saveNmapsToRedis(result *NmapsResult) *NmapsResult {
	for _, v := range result.Data.Info.Hits.Hits {
		ipArr := strings.Split(v.Source.Host, ".")
		if len(ipArr) > 3 {
			if _, ok := matchCDNIP[ipArr[0]+"."+ipArr[1]+"."+ipArr[2]+"."]; ok {
				continue
			}
		}
		rootKey := "nmapssearch:" + v.Source.Host + ":" + v.Source.Title
		if _ret, _ := db.Rdb().Exists(rootKey).Result(); _ret == 0 {
			db.Rdb().LPush("nmapssearch", rootKey)
			_map := uts.Struct2Map(&v.Source)
			saveMapToRedis(_map, rootKey)
		}
	}

	return result
}
