package netspace

import (
	"bytes"
	"encoding/json"
	"io"
	"lufergo/model/db"
	"lufergo/uts"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func InitNetSpace() {
	for _, v := range cdnip {
		iparr := strings.Split(v, "-")
		matchCDNIP[iparr[0][:len(iparr[0])-1]] = true
	}
	//==========test

	_datRet := SearchDat{}
	// _datRet.Sql_text = "google.com"
	_datRet.Sql_text = "status_code=200"
	// _datRet.Sql_text = "ip=\"google.com\""
	// _datRet.Sql_text = "status_code=\"\""
	// _datRet.Sql_text = "server=\"\""
	_datRet.Page_number = 1
	_datRet.Per_page = 10

	// Search(nil, &_datRet)

	// Search_Fofa()
	// Search_Hunter()
	// Search_Nmaps()
	// Search_Quake()
	//=====================
}

var currSearchAPIName = sync.Map{}

func SearchCurrApiName(searchDat []byte) string {
	var dat = &SearchDat{}
	json.Unmarshal(searchDat, dat)
	var keyText = strings.ToLower(dat.Sql_text)
	keyText += strconv.FormatInt(int64(dat.Page_number), 10)
	keyText += strconv.FormatInt(int64(dat.Per_page), 10)
	keyText += strconv.FormatBool(dat.Exclude_empty_data)

	md5Key := uts.GetMd5(keyText)
	str, ok := currSearchAPIName.Load(md5Key)
	if !ok {
		str = "no search"
	}
	return str.(string)
}

func Search(searchDat []byte, searchDatST ...*SearchDat) *NmapsResult {
	var dat *SearchDat
	if len(searchDatST) < 1 {
		dat = &SearchDat{}
		json.Unmarshal(searchDat, dat)
	} else {
		dat = searchDatST[0]
	}

	if dat.Page_number == 0 {
		dat.Page_number = 1
	}

	if dat.Per_page == 0 {
		dat.Per_page = 20

	}

	var sqlText = strings.ToLower(dat.Sql_text)

	// sqlText = "host=\"1.1.1.1\"&&host=\"2.2.2.2\""
	// sqlText = "ip=\"1.1.1.1\"&&ip=\"2.2.2.2\""
	if net.ParseIP(sqlText) != nil {
		isFindIP = true
	} else {
		var domainRegex = regexp.MustCompile(`^[a-zA-Zd-]+(.[a-zA-Zd-]+)*.[a-zA-Z]{2,}$`)
		if domainRegex.MatchString(sqlText) {
			isFindIP = false
		} else {
			regFind := regexp.MustCompile(`ip="([0-9.]+)"`)
			var ipArr = regFind.FindAllStringSubmatch(sqlText, 2)
			for _, v := range ipArr {
				if net.ParseIP(v[1]) == nil {
					isFindIP = false
					break
				}
			}
		}
	}

	// Log("Search", sqlText)

	sqlText += strconv.FormatInt(int64(dat.Page_number), 10)
	sqlText += strconv.FormatInt(int64(dat.Per_page), 10)
	sqlText += strconv.FormatBool(dat.Exclude_empty_data)

	md5Key := uts.GetMd5(sqlText)
	var _nmaps *NmapsResult

	if ret, _ := db.Rdb().Exists(md5Key).Result(); ret == 0 {

		currSearchAPIName.Store(md5Key, "nmaps")
		_nmaps = nmapsSearch(dat)

		currSearchAPIName.Store(md5Key, "fofa")
		_fofa := foFaSearch(dat)

		currSearchAPIName.Store(md5Key, "hunter")
		_hunter := hunterSearch(dat)

		currSearchAPIName.Store(md5Key, "quake")
		_quake := quakeSearch(dat)

		currSearchAPIName.Delete(md5Key)

		if searchDat != nil {
			if _nmaps == nil {
				_nmaps = &NmapsResult{Data: NmapsResultData{Info: NmapsResultDatInfo{Hits: NmapsResultDatHits{Hits: []NmapsResultDatHitsHits{}}}}}
				_nmaps.Status = "error"
				_nmaps.Data.Search_expression = dat.Sql_text
				_nmaps.Data.Page_dict.Page_number = dat.Page_number
				_nmaps.Data.Page_dict.Per_page = dat.Per_page
			}

			megerFofa(_nmaps, _fofa)
			megerHunter(_nmaps, _hunter)
			megerQuake(_nmaps, _quake)

			if _fofa != nil || _hunter != nil || _quake != nil {
				_nmaps.Status = "success"
			}
			if _nmaps.Data.Info.Hits.Total.Value == 0 {
				_nmaps.Data.Info.Hits.Total.Value = len(_nmaps.Data.Info.Hits.Hits)
			}
			Log("Search over")
			jsonBuf, err := json.Marshal(_nmaps)
			if !ChkErr(err) {
				_, err := db.Rdb().Set(md5Key, string(jsonBuf), 24*time.Hour).Result()
				if !ChkErr(err) {
					db.Rdb().Expire(md5Key, db.KeyExpireTime)
				}
			}
			if !isFindIP {
				go getDomainToIP(_nmaps.Data.Info.Hits.Hits, dat)
			}
		}
	} else if searchDat != nil {
		Log("缓存读取")
		ret, _ := db.Rdb().Get(md5Key).Result()

		_nmaps = &NmapsResult{}
		err := json.Unmarshal([]byte(ret), _nmaps)
		ChkErr(err)
	} else {
		Log("数据已存在")
	}

	return _nmaps
}

func getDomainToIP(hits []NmapsResultDatHitsHits, dat *SearchDat) {
	for _, v := range hits {
		if net.ParseIP(v.Source.Host) != nil {
			Log("抓取域名查到的IP", v.Source.Host)
			_dat := &SearchDat{Sql_text: "ip=\"" + v.Source.Host + "\"", Page_number: dat.Page_number, Per_page: dat.Per_page, Exclude_empty_data: dat.Exclude_empty_data, StartTime: dat.StartTime, EndTime: dat.EndTime}
			Search(nil, _dat)
			Log("抓取域名查到的IP,数据完成")
		}
	}
}

// var _score = NmapsResultDatHitsHitsScore{}
//  _score.Id = 0                      //`json:"id"`               //": 5060181,
// _score.Host = "" //": "43.198.13.235",
// _score.Port = 0 //": 443,
// _score.Asn = ""                    //": "16509",
// _score.Chinaadmincode = ""         //": "",
// _score.Continent = ""              //": "亚洲",
// _score.Country = ""                //": "日本",
// _score.Province = ""               //": "日本",
// _score.City = ""                   //": "",
// _score.Cloudinnovation = ""        //": "",
// _score.Latitude = ""               //": "22.257800",
// _score.Longitude = ""              //": "114.165700",
// _score.Utcoffset = ""              //": "AsiaHong_Kong",
// _score.Isp = ""                    //": "亚马逊",
// _score.Asn_name = ""               //": "AMAZON-02",
// _score.Port_insert_date = ""       //": "2024-01-17",
// _score.Protocol = ""               //": "http",
// _score.Tcpdata = ""                //": "HTTP1.1 400 Bad RequestServer: nginxDate: Tue, 16 Jan 2024 09:59:14 GMTContent-Type: texthtmlContent-Length: 150Connection: close<html<head<title400 Bad Request<title<head<body<center<h1400 Bad Request<h1<center<hr<centernginx<center<body<html",
// _score.Asn_as_number = ""          //": "AS16509",
// _score.Asn_as_name = ""            //": "amazon-02",
// _score.Asn_as_country = ""         //": "US",
// _score.Asn_as_range = []string{""} //": ["43.198.0.015", "43.200.0.013"],
// _score.Url = ""                    //": "https:hydw-h5.com",
// _score.Input = ""                  //": "hydw-h5.com",
// _score.Title = ""                  //": "HY游戏",
// _score.Method = ""                 //": "GET",
// _score.Raw_header = ""             //": "",
// _score.Status_code = 0             //": 200,
// _score.Country_svg_url = ""        //": "staticcountry_imgUnited States_US.svg",
// _score.Body_exists = false         //": true

// _score.Favicon = ""                //"error"
// _score.Favicon_md5 = ""            //"d41d8cd98f00b204e9800998ecf8427e"
// _score.Favicon_base64 = ""         // "data:image/x-icon;base64,"
// _score.Csp = ""
// _score.Location = ""
// _score.Cname = ""
// _score.Cdn_name = ""
// _score.Cdn = ""
// _score.Websocket = ""

func megerFofa(nmap *NmapsResult, fofa *FofaResult) {
	if fofa == nil {
		return
	}
	for _, v := range fofa.Results {
		var _score = NmapsResultDatHitsHitsScore{}
		// _score.Id = 0                      //`json:"id"`               //": 5060181,
		_score.Host = v[0] //": "43.198.13.235",
		var port int64 = 0
		port, _ = strconv.ParseInt(v[1], 10, 64)
		_score.Port = int(port) //": 443,

		_score.Asn = v[10]                 //": "16509",
		_score.Chinaadmincode = ""         //": "",
		_score.Continent = ""              //": "亚洲",
		_score.Country = v[3]              //": "日本",
		_score.Province = ""               //": "日本",
		_score.City = v[6]                 //": "",
		_score.Cloudinnovation = ""        //": "",
		_score.Latitude = v[8]             //": "22.257800",
		_score.Longitude = v[7]            //": "114.165700",
		_score.Utcoffset = ""              //": "AsiaHong_Kong",
		_score.Isp = v[10]                 //": "亚马逊",
		_score.Asn_name = v[10]            //": "AMAZON-02",
		_score.Port_insert_date = ""       //": "2024-01-17",
		_score.Protocol = v[2]             //": "http",
		_score.Tcpdata = ""                //": "HTTP1.1 400 Bad RequestServer: nginxDate: Tue, 16 Jan 2024 09:59:14 GMTContent-Type: texthtmlContent-Length: 150Connection: close<html<head<title400 Bad Request<title<head<body<center<h1400 Bad Request<h1<center<hr<centernginx<center<body<html",
		_score.Asn_as_number = v[9]        //": "AS16509",
		_score.Asn_as_name = v[10]         //": "amazon-02",
		_score.Asn_as_country = ""         //": "US",
		_score.Asn_as_range = []string{""} //": ["43.198.0.015", "43.200.0.013"],
		_score.Url = v[12]                 //": "https:hydw-h5.com",
		_score.Input = v[12]               //": "hydw-h5.com",
		_score.Title = v[16]               //": "HY游戏",
		_score.Method = ""                 //": "GET",
		_score.Raw_header = v[18]          //": "",
		_score.Status_code = 0             //": 200.0,
		_score.Country_svg_url = ""        //": "staticcountry_imgUnited States_US.svg",
		_score.Body_exists = false         //": true
		_score.Favicon = ""                //"error"
		_score.Favicon_md5 = ""            //"d41d8cd98f00b204e9800998ecf8427e"
		_score.Favicon_base64 = ""         // "data:image/x-icon;base64,"
		_score.Csp = ""
		_score.Location = ""
		_score.Cname = v[33]
		_score.Cdn_name = ""
		_score.Cdn = ""
		_score.Websocket = ""

		var hits = NmapsResultDatHitsHits{}
		hits.Source = _score
		nmap.Data.Info.Hits.Hits = append(nmap.Data.Info.Hits.Hits, hits)
	}
}

func megerHunter(nmap *NmapsResult, hunter *HunterResult) {
	if hunter == nil {
		return
	}
	for _, v := range hunter.Data.Arr {
		var _score = NmapsResultDatHitsHitsScore{}
		// _score.Id = 0                      //`json:"id"`               //": 5060181,
		_score.Host = v.Ip                          //": "43.198.13.235",
		_score.Port = v.Port                        //": 443,
		_score.Asn = v.Number                       //": "16509",
		_score.Chinaadmincode = ""                  //": "",
		_score.Continent = ""                       //": "亚洲",
		_score.Country = v.Country                  //": "日本",
		_score.Province = v.Province                //": "日本",
		_score.City = v.City                        //": "",
		_score.Cloudinnovation = ""                 //": "",
		_score.Latitude = ""                        //": "22.257800",
		_score.Longitude = ""                       //": "114.165700",
		_score.Utcoffset = ""                       //": "AsiaHong_Kong",
		_score.Isp = v.Isp                          //": "亚马逊",
		_score.Asn_name = ""                        //": "AMAZON-02",
		_score.Port_insert_date = v.Updated_at      //": "2024-01-17",
		_score.Protocol = v.Protocol                //": "http",
		_score.Tcpdata = ""                         //": "HTTP1.1 400 Bad RequestServer: nginxDate: Tue, 16 Jan 2024 09:59:14 GMTContent-Type: texthtmlContent-Length: 150Connection: close<html<head<title400 Bad Request<title<head<body<center<h1400 Bad Request<h1<center<hr<centernginx<center<body<html",
		_score.Asn_as_number = ""                   //": "AS16509",
		_score.Asn_as_name = ""                     //": "amazon-02",
		_score.Asn_as_country = ""                  //": "US",
		_score.Asn_as_range = []string{""}          //": ["43.198.0.015", "43.200.0.013"],
		_score.Url = v.Url                          //": "https:hydw-h5.com",
		_score.Input = v.Domain                     //": "hydw-h5.com",
		_score.Title = v.Web_title                  //": "HY游戏",
		_score.Method = ""                          //": "GET",
		_score.Raw_header = v.Banner                //": "",
		_score.Status_code = float64(v.Status_code) //": 200.0,
		_score.Country_svg_url = ""                 //": "staticcountry_imgUnited States_US.svg",
		_score.Body_exists = false                  //": true

		_score.Favicon = ""        //"error"
		_score.Favicon_md5 = ""    //"d41d8cd98f00b204e9800998ecf8427e"
		_score.Favicon_base64 = "" // "data:image/x-icon;base64,"
		_score.Csp = ""
		_score.Location = ""
		_score.Cname = ""
		_score.Cdn_name = ""
		_score.Cdn = ""
		_score.Websocket = ""

		var hits = NmapsResultDatHitsHits{}
		hits.Source = _score
		nmap.Data.Info.Hits.Hits = append(nmap.Data.Info.Hits.Hits, hits)
	}
}

func megerQuake(nmap *NmapsResult, quake *QuakeResult) {
	if quake == nil {
		return
	}
	for _, v := range quake.Data {
		var _score = NmapsResultDatHitsHitsScore{}
		_score.Id = 0                                    //`json:"id"`               //": 5060181,
		_score.Host = v.Ip                               //": "43.198.13.235",
		_score.Port = v.Port                             //": 443,
		_score.Asn = strconv.FormatInt(int64(v.Asn), 10) //": "16509",
		_score.Chinaadmincode = ""                       //": "",
		_score.Continent = ""                            //": "亚洲",
		_score.Country = v.Location.Country_cn           //": "日本",
		_score.Province = v.Location.Province_cn         //": "日本",
		_score.City = v.Location.City_cn                 //": "",
		_score.Cloudinnovation = ""                      //": "",
		_score.Latitude = ""                             //": "22.257800",
		_score.Longitude = ""                            //": "114.165700",
		_score.Utcoffset = ""                            //": "AsiaHong_Kong",
		_score.Isp = v.Location.Isp                      //": "亚马逊",
		_score.Asn_name = ""                             //": "AMAZON-02",
		_score.Port_insert_date = v.Time                 //": "2024-01-17",
		_score.Protocol = v.Transport                    //": "http",
		_score.Tcpdata = ""                              //": "HTTP1.1 400 Bad RequestServer: nginxDate: Tue, 16 Jan 2024 09:59:14 GMTContent-Type: texthtmlContent-Length: 150Connection: close<html<head<title400 Bad Request<title<head<body<center<h1400 Bad Request<h1<center<hr<centernginx<center<body<html",
		_score.Asn_as_number = ""                        //": "AS16509",
		_score.Asn_as_name = ""                          //": "amazon-02",
		_score.Asn_as_country = ""                       //": "US",
		_score.Asn_as_range = []string{""}               //": ["43.198.0.015", "43.200.0.013"],
		_score.Url = v.Hostname                          //": "https:hydw-h5.com",
		_score.Input = v.Domain                          //": "hydw-h5.com",
		_score.Title = v.Service.Title                   //": "HY游戏",
		_score.Method = ""                               //": "GET",
		_score.Raw_header = v.Service.Response           //": "",
		_score.Status_code = 0                           //": 200,
		_score.Country_svg_url = ""                      //": "staticcountry_imgUnited States_US.svg",
		_score.Body_exists = false                       //": true

		_score.Favicon = v.Service.Http.Favicon.S3_url //"error"
		_score.Favicon_md5 = ""                        //"d41d8cd98f00b204e9800998ecf8427e"
		_score.Favicon_base64 = ""                     // "data:image/x-icon;base64,"
		_score.Csp = ""
		_score.Location = ""
		_score.Cname = ""
		_score.Cdn_name = ""
		_score.Cdn = ""
		_score.Websocket = ""

		var hits = NmapsResultDatHitsHits{}
		hits.Source = _score
		nmap.Data.Info.Hits.Hits = append(nmap.Data.Info.Hits.Hits, hits)
	}
}

// ===========test=============
func Search_Fofa() {
	var dat = &SearchDat{Sql_text: "google.com", Page_number: 1, Per_page: 100}
	foFaSearch(dat)
}

func Search_Hunter() {
	var dat = &SearchDat{Sql_text: "163.com", Page_number: 1, Per_page: 100, IsWeb: "true", StartTime: "2019-01-01 00:00:00"}
	hunterSearch(dat)
}

func Search_Nmaps() {
	var dat = &SearchDat{Sql_text: "163.com", Page_number: 1, Per_page: 100, Exclude_empty_data: true}
	nmapsSearch(dat)
}

func Search_Quake() {
	var dat = &SearchDat{Sql_text: "163.com", Page_number: 1, Per_page: 10, StartTime: "2019-01-01 00:00:00"}
	quakeSearch(dat)
}

//==============================

func callNetSpaceAPI(req *http.Request, reqStr string) *bytes.Buffer {
	requestTime := time.Now().UnixMilli()
	var httpClient = http.Client{Timeout: 10 * time.Second}
	var response *http.Response
	var err error
	if req == nil {
		response, err = httpClient.Get(reqStr)
	} else {
		response, err = httpClient.Do(req)
	}
	if ChkErrNormal(err) {
		return nil
	}
	defer response.Body.Close()
	Log("调用接口用时", time.Now().UnixMilli()-requestTime, response.Request.URL)

	requestTime = time.Now().UnixMilli()
	data := make([]byte, 0, 1024000000)
	respBytes := bytes.NewBuffer(data)
	io.Copy(respBytes, response.Body)
	Log("读取接口数据用时", time.Now().UnixMilli()-requestTime)
	return respBytes
}

func saveMapToRedis(_map map[string]interface{}, rootKey string) {
	for k, v := range _map {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			_mapv := v.(map[string]interface{})
			saveMapToRedis(_mapv, rootKey)
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			saveSlicToRedis(v, rootKey)
		} else {
			_, err := db.Rdb().HSet(rootKey, strings.ToLower(k), v).Result()
			if !ChkErr(err) {
				db.Rdb().Expire(rootKey, db.KeyExpireTime)
			}
		}
	}
}

func saveSlicToRedis(_slice interface{}, rootKey string) {
	_sliceValue := reflect.ValueOf(_slice)
	for i := 0; i < _sliceValue.Len(); i++ {
		v := _sliceValue.Index(i)
		var realVale = v.Interface()
		var realValeRel = reflect.ValueOf(realVale)
		if realValeRel.Kind() == reflect.Map {
			saveMapToRedis(realVale.(map[string]interface{}), rootKey)
		} else if realValeRel.Kind() == reflect.Slice {
			saveSlicToRedis(realVale, rootKey)
		}
	}
}

func checkList(lisName string, key string) bool {
	ret, err := db.Rdb().LRange(lisName, 0, -1).Result()
	if !ChkErr(err) {
		for _, v := range ret {
			if key == v {
				return true
			}
		}
	}
	return false
}
