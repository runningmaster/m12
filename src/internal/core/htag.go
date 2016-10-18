package core

import (
	"fmt"
	"strings"
)

var (
	listHTag = map[string]struct{}{
		"geoapt.ru":           {},
		"geoapt.ua":           {},
		"sale-in.monthly.by":  {},
		"sale-in.monthly.kz":  {},
		"sale-in.monthly.ua":  {},
		"sale-in.weekly.ua":   {},
		"sale-in.daily.kz":    {},
		"sale-in.daily.ua":    {},
		"sale-out.monthly.by": {},
		"sale-out.monthly.kz": {},
		"sale-out.monthly.ru": {},
		"sale-out.monthly.ua": {},
		"sale-out.weekly.ru":  {},
		"sale-out.weekly.ua":  {},
		"sale-out.daily.by":   {},
		"sale-out.daily.kz":   {},
		"sale-out.daily.ua":   {},
	}

	convTags = map[string]string{
		// version 1 -> version 3
		"data.geostore":         "geoapt.ua",
		"data.sale-inp.monthly": "sale-in.monthly.ua",
		"data.sale-inp.weekly":  "sale-in.weekly.ua",
		"data.sale-inp.daily":   "sale-in.daily.ua",
		"data.sale-out.monthly": "sale-out.monthly.ua",
		"data.sale-out.weekly":  "sale-out.weekly.ua",
		"data.sale-out.daily":   "sale-out.daily.ua",
		// version 2 -> version 3
		"data.geoapt.ru":           "geoapt.ru",
		"data.geoapt.ua":           "geoapt.ua",
		"data.sale-inp.monthly.kz": "sale-out.monthly.kz",
		"data.sale-inp.monthly.ua": "sale-in.monthly.ua",
		"data.sale-inp.weekly.ua":  "sale-in.weekly.ua",
		"data.sale-inp.daily.kz":   "sale-in.daily.kz",
		"data.sale-inp.daily.ua":   "sale-in.daily.ua",
		"data.sale-out.monthly.kz": "sale-out.monthly.kz",
		"data.sale-out.monthly.ua": "sale-out.monthly.ua",
		"data.sale-out.weekly.ua":  "sale-out.weekly.ua",
		"data.sale-out.daily.by":   "sale-out.daily.by",
		"data.sale-out.daily.kz":   "sale-out.daily.kz",
		"data.sale-out.daily.ua":   "sale-out.daily.ua",
	}
)

func normHTag(t string) string {
	s, ok := convTags[t]
	if ok {
		return s
	}
	return t
}

func testHTag(t string) error {
	t = strings.ToLower(t)
	_, ok1 := convTags[t]
	_, ok2 := listHTag[t]

	if ok1 || ok2 {
		return nil
	}

	return fmt.Errorf("core: invalid htag %s", t)
}
