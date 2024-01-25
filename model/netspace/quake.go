package netspace

import (
	"bytes"
	"encoding/json"
	"lufergo/model/db"
	"lufergo/uts"
	"net/http"
	"regexp"
	"strings"
)

func quakeSearch(dat *SearchDat) *QuakeResult {
	requestStr := db.Conf.String("quake_api")

	var sqlText = strings.ToLower(dat.Sql_text)
	sqlText = strings.ReplaceAll(sqlText, "domain=", "domain:")

	sqlText = strings.ReplaceAll(sqlText, "title=", "title:") //title="webui"	 从标题中搜索"webui"

	rege := regexp.MustCompile(`banner=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")            //banner='HTTP/1.1 400 Bad'	 数据包筛选
	sqlText = strings.ReplaceAll(sqlText, "ip=", "ip:")     //ip="18.166.73.142"	搜索ip为"18.166.73.142"的资产
	sqlText = strings.ReplaceAll(sqlText, "port=", "port:") //port="6379"	查找对应"6379"端口的资产
	sqlText = strings.ReplaceAll(sqlText, "host=", "host:") //host="18.166.73.142"	搜索host为"18.166.73.142"的资产
	sqlText = strings.ReplaceAll(sqlText, "net=", "ip:")    //net='16.163.128.0/17'网段信息
	sqlText = strings.ReplaceAll(sqlText, "body=", "body:") //body="网络空间测绘"	从html正文中搜索"网络空间测绘"

	rege = regexp.MustCompile(`body_md5=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //body_md5="a63df013125bb17528cc868347009c3f"	通过html正文计算的md5值搜索
	rege = regexp.MustCompile(`body_mmh3=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //body_mmh3="1178747436"	通过html正文计算的mmh3值搜索
	rege = regexp.MustCompile(`body_sha256=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //body_sha256="08d948ac7131b73e73ddbfe44d1a10857876909a3aa5c390746340c931ae83cd"	通过html正文计算的sha256值搜索
	rege = regexp.MustCompile(`body_simhash=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //body_simhash="15607296109983774326"	通过html正文计算的simhash值搜索

	rege = regexp.MustCompile(`header_md5=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //header_md5="f671390ecd477cab95181da75eea3b26"	通过http/https的响应头计算的md5值搜索
	rege = regexp.MustCompile(`header_mmh3=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //header_mmh3="-1132468473"	的header_mmh3值搜索
	rege = regexp.MustCompile(`header_sha256=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //header_sha256="afcf009865edf35af82f2252239a5e2279f03fd08c2b57c5aa089a5c979789af"通过http/https的响应头计算通过http/https的响应头计算的header_sha256值搜索

	sqlText = strings.ReplaceAll(sqlText, "status_code=", "status_code:") //status_code="402"	查询服务器状态为"402"的资产
	sqlText = strings.ReplaceAll(sqlText, "webserver=", "server:")        //webserver="nginx"搜索web服务器是"nginx"的资产
	sqlText = strings.ReplaceAll(sqlText, "scheme=", "service:")          //scheme="http"	搜索协议使用"http"的资产
	sqlText = strings.ReplaceAll(sqlText, "url=", "domain:")              //url="https://18.166.87.172:443"	和示例展示的那样，查询url路径
	rege = regexp.MustCompile(`content_length=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")           //content_length="18"	查询header中content_length
	sqlText = strings.ReplaceAll(sqlText, "input=", "ip:") //input="18.166.87.172:443"查询ip和端口的组合
	rege = regexp.MustCompile(`content_type=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //content_type="text/html"查询header中centent_type
	rege = regexp.MustCompile(`method=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //method="GET"	请求方法不会自动大小写转换
	rege = regexp.MustCompile(`path=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")
	sqlText = strings.ReplaceAll(sqlText, "raw_header=", "headers:") //raw_header="HTTP/1.1 200 "		raw_header是将所有header组成一行的文本字符串
	rege = regexp.MustCompile(`request=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //request="GET / HTTP/1.1"		请求头
	rege = regexp.MustCompile(`time=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")
	rege = regexp.MustCompile(`words=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //words="1652"	body内单词的个数
	rege = regexp.MustCompile(`lines=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //lines="160"	body内行的个数
	rege = regexp.MustCompile(`failed=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")
	sqlText = strings.ReplaceAll(sqlText, "asn==", "asn:")     //asn=='16509'	根据asn编码进行搜索
	sqlText = strings.ReplaceAll(sqlText, "asn_name=", "org:") //asn_name='AMAZON-02'	根据组织名称搜索
	rege = regexp.MustCompile(`continent==".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")                       //continent=='亚洲'	搜索大陆为亚洲的
	sqlText = strings.ReplaceAll(sqlText, "country==", "country:")     //country=='中国'	搜索国家为中国的
	sqlText = strings.ReplaceAll(sqlText, "province=", "ip.province=") //province='江西'	搜索省份为江西
	sqlText = strings.ReplaceAll(sqlText, "city=", "city:")            //city='南昌'	搜索城市为南昌
	rege = regexp.MustCompile(`latitude=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //latitude='22.257800'	经度
	rege = regexp.MustCompile(`longitude=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //longitude='114.165700'	纬度
	rege = regexp.MustCompile(`utc=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //utc='Asia/Hong_Kong'	utc时区
	rege = regexp.MustCompile(`favicon_path=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //favicon_path='/favicon.ico'	favicon的图标相对路径，需要使用url拼接才是绝对路径
	rege = regexp.MustCompile(`favicon_base64=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //搜索favicon图标的base64值
	rege = regexp.MustCompile(`favicon_md5=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "") //搜索favicon图标的md5值
	rege = regexp.MustCompile(`jarm=".*"`)
	sqlText = rege.ReplaceAllString(sqlText, "")          //jarm指纹
	sqlText = strings.ReplaceAll(sqlText, "isp=", "isp:") //isp='亚马逊'

	Log("Search", sqlText)

	param := map[string]interface{}{
		"query":      sqlText,
		"start":      dat.Page_number,
		"size":       dat.Per_page,
		"start_time": dat.StartTime,
		"include":    strings.Split(db.Conf.String("quake_fields"), ","),
		// "end_time": endtime
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
	req.Header.Set("X-QuakeToken", db.Conf.String("quake_token"))

	// Log(param)
	respBytes := callNetSpaceAPI(req, "")
	if respBytes == nil {
		return nil
	}
	_result := &QuakeResult{}
	json.Unmarshal(respBytes.Bytes(), _result)
	if _result.Code != "" {
		Log("调用错误", _result)
		return nil
	}
	saveQuakeToRedis(_result)

	return _result
}

func saveQuakeToRedis(result *QuakeResult) *NmapsResult {
	var _result = &NmapsResult{}
	for _, v := range result.Data {
		ipArr := strings.Split(v.Ip, ".")
		if len(ipArr) > 3 {
			if _, ok := matchCDNIP[ipArr[0]+"."+ipArr[1]+"."+ipArr[2]+"."]; ok {
				continue
			}
		}
		rootKey := "quakesearch:" + v.Ip + ":" + v.Service.Title
		if _ret, _ := db.Rdb().Exists(rootKey).Result(); _ret == 0 {
			db.Rdb().LPush("quakesearch", rootKey)
			_map := uts.Struct2Map(&v)
			saveMapToRedis(_map, rootKey)
		}
	}
	return _result
}
