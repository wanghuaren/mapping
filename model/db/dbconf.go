package db

import (
	"lufergo/uts"
	"time"
)

var Conf = uts.Conf
var Log = uts.Log
var ChkErr = uts.ChkErr
var ChkErrNormal = uts.ChkErrNormal

var KeyExpireTime = 7 * 24 * time.Hour
