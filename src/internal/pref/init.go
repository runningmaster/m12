package pref

import (
	"expvar"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sasbury/mini"
)

func init0(p ...pref) {
	// see above
}

func init1(p ...pref) {
	// not implemented
}

func init2(p ...pref) {
	if len(os.Args) < 2 {
		return
	}

	cfg, err := mini.LoadConfiguration(os.Args[1])
	if err != nil {
		return
	}

	for i := range p {
		setFromConfig(p[i], cfg)
	}
}

func init3(p ...pref) {
	for i := range p {
		setFromEvar(p[i])
	}
}

func init4(p ...pref) {
	for i := range p {
		setFromFlag(p[i])
	}
	flag.Parse()
}

func init5(key string) {
	expvar.NewInt(key).Set(1)
	d, _ := strconv.Atoi(expvar.Get(key).String())
	Debug = d == 1
}

func envVar(p pref) string {
	return fmt.Sprintf(envFormat, strings.ToUpper(strings.Replace(p.name, "-", "_", -1)))
}

func setFromConfig(p pref, cfg *mini.Config) {
	switch x := p.value.(type) {
	case *string:
		*x = cfg.String(p.name, *x)
	case *bool:
		*x = cfg.Boolean(p.name, *x)
	default:
		panic("pref: unreachable: config")
	}
}

func setFromEvar(p pref) {
	v := os.Getenv(envVar(p))
	if v == "" {
		return
	}

	switch x := p.value.(type) {
	case *string:
		*x = v
	case *bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		*x = b
	default:
		panic("pref: unreachable: evar")
	}
}

func setFromFlag(p pref) {
	switch x := p.value.(type) {
	case *string:
		flag.StringVar(x, p.name, *x, p.usage)
	case *bool:
		flag.BoolVar(x, p.name, *x, p.usage)
	default:
		panic("pref: unreachable: flag")
	}
}
