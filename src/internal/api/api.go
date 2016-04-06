package api

import (
	"strings"

	"internal/core"
	"internal/flag"
	"internal/server"
)

var (
	mapCoreHandlers map[string]core.Handler
	mapHTTPHandlers = map[string]bundle{
		"GET:/":     {use(pipeHead, pipeGzip, pipe(root), pipeFail, pipeTail), nil},
		"GET:/ping": {use(pipeHead, pipeGzip, pipe(exec), pipeFail, pipeTail), core.Ping},

		"POST:/system/get-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysAuth, flag.OpGet)},
		"POST:/system/set-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysAuth, flag.OpSet)},
		"POST:/system/del-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysAuth, flag.OpDel)},
		"POST:/system/get-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkAddr, flag.OpGet)},
		"POST:/system/set-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkAddr, flag.OpSet)},
		"POST:/system/del-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkAddr, flag.OpDel)},
		"POST:/system/get-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkDrug, flag.OpGet)},
		"POST:/system/set-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkDrug, flag.OpSet)},
		"POST:/system/del-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkDrug, flag.OpDel)},
		"POST:/system/get-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkStat, flag.OpGet)},
		"POST:/system/set-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkStat, flag.OpSet)},
		"POST:/system/del-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SysOp(core.SysLinkStat, flag.OpDel)},

		"POST:/upload/geoapt.ua":           {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/geoapt.ru":           {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-in.monthly.ua":  {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-in.monthly.kz":  {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-in.weekly.ua":   {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-in.daily.ua":    {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-in.daily.kz":    {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-out.monthly.ua": {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-out.monthly.kz": {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-out.weekly.ua":  {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-out.daily.ua":   {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-out.daily.kz":   {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},
		"POST:/upload/sale-out.daily.by":   {use(pipeHead, pipeAuth, pipe(exec), pipeFail, pipeTail), core.ToS3},

		// => Debug mode only, when flag.Debug == true
		"GET:/debug/info":               {use(pipeHead, pipeGzip, pipe(exec), pipeFail, pipeTail), core.Info}, // ?
		"GET:/debug/vars":               {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // expvar
		"GET:/debug/pprof/":             {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // net/http/pprof
		"GET:/debug/pprof/cmdline":      {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // net/http/pprof
		"GET:/debug/pprof/profile":      {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // net/http/pprof
		"GET:/debug/pprof/symbol":       {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // net/http/pprof
		"GET:/debug/pprof/trace":        {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // net/http/pprof
		"GET:/debug/pprof/goroutine":    {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // runtime/pprof
		"GET:/debug/pprof/threadcreate": {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // runtime/pprof
		"GET:/debug/pprof/heap":         {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // runtime/pprof
		"GET:/debug/pprof/block":        {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},       // runtime/pprof

		// => Workarounds for 404/405
		"GET:/error/404": {use(pipeHead, pipe(e404), pipeFail, pipeTail), nil},
		"GET:/error/405": {use(pipeHead, pipe(e405), pipeFail, pipeTail), nil},
	}
)

func init() {
	mapCoreHandlers = make(map[string]core.Handler, len(mapHTTPHandlers))
	for k, v := range mapHTTPHandlers {
		s := strings.Split(k, ":")
		server.RegHandler(s[0], s[1], v.h)
		if v.f != nil {
			mapCoreHandlers[s[1]] = v.f
		}
	}
}
