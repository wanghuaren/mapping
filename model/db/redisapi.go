package db

import (
	"lufergo/uts"
	"time"

	"github.com/go-redis/redis"
)

var _rdb *redis.Client

var _rdbSlave *redis.Client

func Rdb() *redis.Client {
	if checkRedis(_rdb) {
		return _rdb
	} else {
		Log("redis 异常")
		return nil
	}
}

func rdbSlave() *redis.Client {
	if checkRedis(_rdb) {
		if checkRedis(_rdbSlave) {
			return _rdbSlave
		} else {
			return _rdb
		}
	} else {
		Log("redis 异常")
		return nil
	}
}

func InitRedis() {
	redisConn()
}

func redisConn() {
	_rdb = redis.NewClient(&redis.Options{
		Addr:         Conf.String("redis_host") + ":" + Conf.String("redis_port"),
		Password:     Conf.String("redis_pwd"),
		Network:      "tcp", //网络类型，tcp or unix，默认tcp
		DB:           0,     // redis数据库index
		MinIdleConns: 10,    //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		MaxRetries:   0,     // 命令执行失败时，最多重试多少次，默认为0即不重试
	})
	if checkRedis(_rdb) {
		Log("创建Redis完成")
	} else {
		Log("Redis重试")
		time.Sleep(time.Second * 6)
		redisConn()
	}
}

const minInterval int64 = 5

var prevCheckTime int64 = 0
var prevCheckRedis *redis.Client

func checkRedis(r *redis.Client) bool {
	currCheckTime := time.Now().Unix()
	if r == prevCheckRedis && currCheckTime-prevCheckTime < minInterval {
		prevCheckTime = currCheckTime
		prevCheckRedis = r
		return true
	}
	prevCheckTime = currCheckTime
	prevCheckRedis = r

	// 测试连接
	_, err := r.Ping().Result()
	return !ChkErrNormal(err, "redis 检测异常")
}

func DB2Cache() {
	Rdb().FlushAll()
	Log("DB2Cache 导入")
	// Expire 命令用于设置 key 的过期时间，key 过期后将不再可用。单位以秒计。
	// PERSIST 命令用于移除给定 key 的过期时间，使得 key 永不过期。
	// TTL 命令以秒为单位返回 key 的剩余过期时间
	tableData := uts.FindAllTable()
	for _tabName, _ := range tableData {
		Log(_tabName)
		_tableData := uts.FindDBTableData(_tabName)
		for _, _map := range _tableData {
			//f458 h9 n150
			_mKeyIP := ""
			_mKeyTitle := ""
			for _k, _v := range _map {
				if _k == "ip" || _k == "host" {
					_mKeyIP = _v
				} else if _k == "title" || _k == "web_title" {
					_mKeyTitle = _v
				}
			}
			_rootKey := _tabName + ":" + _mKeyIP + ":" + _mKeyTitle
			for _fk, _fv := range _map {
				_, err := Rdb().HSet(_rootKey, _fk, _fv).Result()
				if !ChkErr(err) {
					// Rdb().Expire(_rootKey, KeyExpireTime)
				}
			}
			Rdb().RPush(_tabName, _rootKey)
		}
	}
	Log("DB2Cache 结束")

	go func() {
		defer uts.ChkRecover()
		RDDump()
	}()
}

func BackUp() {
	uts.BackupMysql()
}

func Cache2DB() {
	BackUp()
	var _tableNames = []string{"fofasearch", "huntersearch", "nmapssearch", "quakesearch"}
	for _, _tabName := range _tableNames {
		uts.ExecMySQLCommand("drop table if exists `" + _tabName + "`")
		_isCraedTable := false
		ret, err := Rdb().LRange(_tabName, 0, -1).Result()
		var _columnMap = map[string]bool{}
		if !ChkErr(err) {
			_isCraedTable = false
			count := 0
			for _, rootKey := range ret {
				_tableMapDat, err := Rdb().HGetAll(rootKey).Result()
				if !ChkErr(err) {
					if !_isCraedTable {
						_createTableStr := "create table `" + _tabName + "` ("
						for k, _ := range _tableMapDat {
							_columnMap[k] = true
							_createTableStr += "`" + k + "` longtext,"
						}
						_createTableStr = _createTableStr[:len(_createTableStr)-1]
						_createTableStr += ") default charset=utf8mb4"
						_isCraedTable = true
						Log("_createTableStr", _createTableStr)
						uts.ExecMySQLCommand(_createTableStr)
					}

					_sqlColumnStr := ""
					_sqlColumnValue := []any{}
					_sqlColumnValueNumStr := ""
					for k, v := range _tableMapDat {
						if _, ok := _columnMap[k]; ok {
							_sqlColumnStr += k + ","
							_sqlColumnValueNumStr += "?,"
							_sqlColumnValue = append(_sqlColumnValue, v)
						} else {
							Log(_tabName + " 中，字段 " + k + " 未被包含使用")
						}
					}
					_sqlColumnStr = _sqlColumnStr[:len(_sqlColumnStr)-1]
					_sqlColumnValueNumStr = _sqlColumnValueNumStr[:len(_sqlColumnValueNumStr)-1]
					_insertTableStr := "insert into " + _tabName + "(" + _sqlColumnStr + ") values(" + _sqlColumnValueNumStr + ");"
					// Log("_insertTableStr", _insertTableStr, _sqlColumnValue)
					Log(count)
					count++
					uts.ExecMySQLCommandInsert(_insertTableStr, _sqlColumnValue...)
				}
			}
		}
	}
	Log("redis to mysql success!!")
}

//http && status_code="200"

func RDAdd(data interface{}, fields ...[]string) interface{} {

	return nil
}

func RDDel(data interface{}, fields ...[]string) interface{} {

	return nil
}
func RDEdit(data interface{}, fields ...[]string) interface{} {

	return nil
}

func RDFind(data interface{}, fields ...[]string) interface{} {

	return nil
}
func RDDump(now ...bool) {
	if len(now) < 1 {
		time.Sleep(time.Hour * 2)
		Log("Redis 数据入库")
		Cache2DB()
		Rdb().BgSave()
		RDDump()
	} else {
		RDDump()
	}
}
