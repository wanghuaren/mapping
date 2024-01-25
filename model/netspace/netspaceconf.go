package netspace

import "lufergo/uts"

var Log = uts.Log
var ChkErr = uts.ChkErr
var ChkErrNormal = uts.ChkErrNormal

var isFindIP = false

type SearchDat struct {
	Sql_text           string
	Page_number        int
	Per_page           int
	Exclude_empty_data bool

	StartTime string
	EndTime   string
	IsWeb     string
}

type BackCurrSearchAPI struct {
	SearchAPI string
}

//==============nmaps

type NmapsResult struct {
	Status string          `json:"status"`
	Data   NmapsResultData `json:"data"`
}

type NmapsResultData struct {
	Info              NmapsResultDatInfo            `json:"info"`
	Page_dict         NmapsResultDatPageDict        `json:"page_dict"`
	Search_expression string                        `json:"search_expression"` //": "163.com",
	Export_fields     []NmapsResultDatExport_fields `json:"export_fields"`
}

type NmapsResultDatInfo struct {
	Took      int                  `json:"took"`      //": 760,
	Timed_out bool                 `json:"timed_out"` //": false,
	Shards    NmapsResultDatShards `json:"_shards"`
	Hits      NmapsResultDatHits   `json:"hits"`
}

type NmapsResultDatShards struct {
	Total      int `json:"total"`      //": 1,
	Successful int `json:"successful"` //": 1,
	Skipped    int `json:"skipped"`    //": 0,
	Failed     int `json:"failed"`     //": 0 }
}

type NmapsResultDatTotal struct {
	Value    int    `json:"value"`    //": 10668,
	Relation string `json:"relation"` //": "eq" }
}

type NmapsResultDatHits struct {
	Total     NmapsResultDatTotal      `json:"total"`     //": ,
	Max_score string                   `json:"max_score"` //": null,
	Hits      []NmapsResultDatHitsHits `json:"hits"`
}

type NmapsResultDatHitsHits struct {
	Index   string   `json:"_index"`   //": "cm_es_result_v2",
	Id      int      `json:"_id"`      //": "5060181",
	Score   string   `json:"_score"`   //": null,
	Ignored []string `json:"_ignored"` //": [
	//   "raw_header.keyword",
	//   "body.keyword",
	//   "tcpdata.keyword"
	// ],
	Source NmapsResultDatHitsHitsScore `json:"_source"`
	Sort   []int                       `json:"sort"`
}
type NmapsResultDatHitsHitsScore struct {
	Id               int      `json:"id"`               //": 5060181,
	Host             string   `json:"host"`             //": "43.198.13.235",
	Port             int      `json:"port"`             //": 443,
	Asn              string   `json:"asn"`              //": "16509",
	Chinaadmincode   string   `json:"chinaadmincode"`   //": "",
	Continent        string   `json:"continent"`        //": "亚洲",
	Country          string   `json:"country"`          //": "日本",
	Province         string   `json:"province"`         //": "日本",
	City             string   `json:"city"`             //": "",
	Cloudinnovation  string   `json:"cloudinnovation"`  //": "",
	Latitude         string   `json:"latitude"`         //": "22.257800",
	Longitude        string   `json:"longitude"`        //": "114.165700",
	Utcoffset        string   `json:"utcoffset"`        //": "AsiaHong_Kong",
	Isp              string   `json:"isp"`              //": "亚马逊",
	Asn_name         string   `json:"asn_name"`         //": "AMAZON-02",
	Port_insert_date string   `json:"port_insert_date"` //": "2024-01-17",
	Protocol         string   `json:"protocol"`         //": "http",
	Tcpdata          string   `json:"tcpdata"`          //": "HTTP1.1 400 Bad RequestServer: nginxDate: Tue, 16 Jan 2024 09:59:14 GMTContent-Type: texthtmlContent-Length: 150Connection: close<html<head<title400 Bad Request<title<head<body<center<h1400 Bad Request<h1<center<hr<centernginx<center<body<html",
	Asn_as_number    string   `json:"asn_as_number"`    //": "AS16509",
	Asn_as_name      string   `json:"asn_as_name"`      //": "amazon-02",
	Asn_as_country   string   `json:"asn_as_country"`   //": "US",
	Asn_as_range     []string `json:"asn_as_range"`     //": ["43.198.0.015", "43.200.0.013"],
	Url              string   `json:"url"`              //": "https:hydw-h5.com",
	Input            string   `json:"input"`            //": "hydw-h5.com",
	Title            string   `json:"title"`            //": "HY游戏",
	Method           string   `json:"method"`           //": "GET",
	Raw_header       string   `json:"raw_header"`       //": "",
	Status_code      float64  `json:"status_code"`      //": 200.0,
	Country_svg_url  string   `json:"country_svg_url"`  //": "staticcountry_imgUnited States_US.svg",
	Body_exists      bool     `json:"body_exists"`      //": true
	Favicon          string   `json:"favicon"`
	Favicon_md5      string   `json:"favicon_md5"`
	Favicon_base64   string   `json:"favicon_base64"`
	Csp              string   `json:"csp"`
	Location         string   `json:"location"`
	Cname            string   `json:"cname"`
	Cdn_name         string   `json:"cdn_name"`
	Cdn              string   `json:"cdn"`
	Websocket        string   `json:"websocket"`
}

type NmapsResultDatPageDict struct {
	Page_number          int  `json:"page_number"`          //": 1,
	Per_page             int  `json:"per_page"`             //": 100,
	Num_pages            int  `json:"num_pages"`            //": 107,
	Count                int  `json:"count"`                //": 10668,
	Has_next             bool `json:"has_next"`             //": true,
	Has_previous         bool `json:"has_previous"`         //": false,
	Has_other_pages      bool `json:"has_other_pages"`      //": true,
	Next_page_number     int  `json:"next_page_number"`     //": 2,
	Previous_page_number int  `json:"previous_page_number"` //": 1
}

type NmapsResultDatExport_fields struct {
	Value string `json:"value"` //host"
	Label string `json:"label"` //主机地址"
}

//============================

// =====fofa
type FofaResult struct {
	Error            bool // 是否出现错误
	Errmsg           string
	Consumed_fpoint  int    // 应扣F点
	Required_fpoints int    // 实扣F点
	Size             int    // 查询总数量
	Page             int    // 当前页码
	Mode             string //"extended",
	Query            string //"title=\"bing\"", // 查询语句
	Results          [][]string
}

// ==========hunter
type HunterResult struct {
	Code    int
	Message string
	Data    HunterResultData
}

type HunterResultData struct {
	Account_type  string // "个人账号",
	Total         int    // 186112,
	Time          int    // 6712,
	Arr           []HunterResultDatArr
	Consume_quota string //": "消耗积分：50",
	Rest_quota    string //": "今日剩余积分：297",
	Syntax_prompt string
}

type HunterResultDatArr struct {
	Is_risk          string
	Url              string //": "https://gmall.dev.webapp.163.com:8072",
	Ip               string //": "42.186.175.33",
	Port             int    //": 8072,
	Web_title        string //": "5周年测试",
	Domain           string //": "gmall.dev.webapp.163.com",
	Is_risk_protocol string //": "",
	Protocol         string //": "https",
	Base_protocol    string //": "tcp",
	Status_code      int    //": 200,
	Component        []HunterResultDatArrCmp
	Os               string //": "",
	Company          string //": "广州网易计算机系统有限公司",
	Number           string //": "粤B2-20090191-18",
	Country          string //": "中国",
	Province         string //": "浙江省",
	City             string //": "杭州市",
	Updated_at       string //": "2024-01-19",
	Is_web           string //": "是",
	As_org           string //": "网易",
	Isp              string //": "网易",
	Banner           string //": "HTTP/1.1 403 Forbidden\\r\\nServer: nginx\\r\\nDate: Thu, 18 Jan 2024 19:41:07 GMT\\r\\nContent-Type: text/plain; charset=UTF-8\\r\\nContent-Length: 55\\r\\nConnection: keep-alive\\r\\nX-Trace-ID: 103d07839f581cfade049e70dcc85226\\r\\nntes-trace-id: 8afe5713ed69af0:8afe5713ed69af0:0:1\\r\\n\\r\\n403 Forbidden\\n\\nAccess was denied to this resource.\\n\\n   ",
	Vul_list         string //": ""
}

type HunterResultDatArrCmp struct {
	Name    string
	Version string
}

//===================

//==========quake
type QuakeResult struct {
	Code    string //": 0,
	Message string //Successful.",
	Data    []QuakeResultData
	Meta    QuakeResultDatMeta //":
}

type QuakeResultData struct {
	Components []QuakeResultDatComp   //": [
	Hostname   string                 //": "mx.noordsestern.de",
	Org        string                 //": "Hetzner Online GmbH",
	Port       int                    //": 9071,
	Service    QuakeResultDatService  //":
	Ip         string                 //": "121.41.89.91",
	Domain     string                 //": "www.cnbxfc.net",
	Location   QuakeResultDatLocation //":
	Transport  string                 //": "tcp",
	Time       string                 //": "2024-01-20T03:54:13.191Z",
	Asn        int                    //": 37963,
	Id         string                 //": "www.cnbxfc.net_80_tcp"
}

type QuakeResultDatComp struct {
	Product_level   string   //": "中间支撑层",
	Product_type    []string //": ["Web服务器"],
	Product_vendor  string   //": "F5Networks网络公司",
	Product_name_cn string   //": "Nginx Web服务器",
	Product_catalog []string //": ["Web系统与应用"],
	Version         string   //": ""
}

type QuakeResultDatService struct {
	Response string //": "",
	Title    string //": "玻璃纤维复合材料信息网--复合材料行业优质产品推荐平台s"
	Http     QuakeResultDatHttp
}

type QuakeResultDatLocation struct {
	Province_cn string //": "浙江省",
	Isp         string //": "阿里云",
	Province_en string //": "Zhejiang",
	Country_en  string //": "China",
	District_cn string //": "",
	City_en     string //": "Hangzhou City",
	District_en string //": "",
	Country_cn  string //": "中国",
	City_cn     string //": "杭州市"
}

type QuakeResultDatMeta struct {
	Pagination QuakeResultDatPagination
}

type QuakeResultDatPagination struct {
	Count      int //": 10,
	Page_index int //": 1,
	Page_size  int //": 10,
	Total      int //": 2386747185
}

type QuakeResultDatHttp struct {
	Server      string //": "",
	Status_code int    //": "玻璃纤维复合材料信息网--复合材料行业优质产品推荐平台s"
	Host        string
	Body        string
	Title       string
	Favicon     QuakeResultDatFavicon
}

type QuakeResultDatFavicon struct {
	Data   string //": "",
	Hash   string //": "玻璃纤维复合材料信息网--复合材料行业优质产品推荐平台s"
	S3_url string
}
