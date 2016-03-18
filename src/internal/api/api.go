package api

import (
	"strings"

	"internal/core"
	"internal/server"
)

var (
	mapFuncs    map[string]coreFunc
	mapHandlers = map[string]bundle{
		"GET:/":     {use(pipeHead, pipeGzip, pipe(root), pipeFail, pipeTail), nil},
		"GET:/ping": {use(pipeHead, pipeGzip, pipe(exec), pipeFail, pipeTail), core.Ping},

		"POST:/get-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.GetAuth},
		"POST:/set-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SetAuth},
		"POST:/del-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.DelAuth},
		"POST:/get-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.GetLinkAddr},
		"POST:/set-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SetLinkAddr},
		"POST:/del-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.DelLinkAddr},
		"POST:/get-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.GetLinkDrug},
		"POST:/set-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SetLinkDrug},
		"POST:/del-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.DelLinkDrug},
		"POST:/get-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.GetLinkStat},
		"POST:/set-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.SetLinkStat},
		"POST:/del-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.DelLinkStat},

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
	mapFuncs = make(map[string]coreFunc, len(mapHandlers))
	for k, v := range mapHandlers {
		s := strings.Split(k, ":")
		server.RegHandler(s[0], s[1], v.h)
		if v.f != nil {
			mapFuncs[s[1]] = v.f
		}
	}
}
