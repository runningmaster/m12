package api

import (
	"net/http"
	"strings"

	"internal/core"
)

/*
curl -v -k -X 'POST' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJza2V5IjoibWFzdGVya2V5IiwiaHRhZyI6Imdlb2FwdC51YSIsIm5hbWUiOiLQkNC_0YLQtdC60LAgMyIsImhlYWQiOiLQkdCG0JvQkCDQoNCe0JzQkNCo0JrQkCIsImFkZHIiOiLQkdC-0YDQuNGB0L_QvtC70Ywg0YPQuy4g0JrQuNC10LLRgdC60LjQuSDQqNC70Y_RhSwgOTgiLCJjb2RlIjoiMTIzNDU2In0.EaVgSUK2oYWyEunPkc9HhmOtjNOyTjdlyQMe81slBHM' \
-H 'Content-Encoding: gzip' \
-H 'Content-Type: application/json; charset=utf-8' \
-T 'data.json.gz' \
http://localhost:8080/upload
*/

var (
	mapCoreHandlers map[string]core.Handler
	mapHTTPHandlers = map[string]bundle{
		"GET:/":     {use(pipeHead, pipeGzip, pipe(root), pipeFail, pipeTail), nil},
		"GET:/ping": {use(pipeHead, pipeGzip, pipe(exec), pipeFail, pipeTail), core.Ping},

		"POST:/system/get-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("get", "auth")},
		"POST:/system/set-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("set", "auth")},
		"POST:/system/del-auth":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("del", "auth")},
		"POST:/system/get-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("get", "addr")},
		"POST:/system/set-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("set", "addr")},
		"POST:/system/del-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("del", "addr")},
		"POST:/system/get-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("get", "drug")},
		"POST:/system/set-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("set", "drug")},
		"POST:/system/del-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("del", "drug")},
		"POST:/system/get-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("get", "stat")},
		"POST:/system/set-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("set", "stat")},
		"POST:/system/del-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.RunC("del", "stat")},
		"POST:/system/pop-data":      {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.Popd},

		"POST:/upload": {use(pipeHead, pipeAuth, pipeMeta, pipe(exec), pipeFail, pipeTail), core.Upld},

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

// Init is caled from other package for manually initialization
func Init(regFunc func(string, string, http.Handler)) error {
	err := core.Init()
	if err != nil {
		return err
	}

	regHTTPHandlers(regFunc)

	return nil
}

func regHTTPHandlers(regFunc func(string, string, http.Handler)) {
	if regFunc == nil {
		return
	}

	mapCoreHandlers = make(map[string]core.Handler, len(mapHTTPHandlers))
	for k, v := range mapHTTPHandlers {
		s := strings.Split(k, ":")
		regFunc(s[0], s[1], v.h)
		if v.f != nil {
			mapCoreHandlers[s[1]] = v.f
		}
	}
}
