package uuid

import (
	"fmt"
	"sync"

	"github.com/rogpeppe/fastuuid"
)

var uuidPool = sync.Pool{
	New: func() interface{} {
		g, err := fastuuid.NewGenerator()
		if err != nil {
			return err
		}
		return g
	},
}

func getGenerator() (*fastuuid.Generator, error) {
	switch g := uuidPool.Get().(type) {
	case *fastuuid.Generator:
		return g, nil
	case error:
		return nil, g
	}

	return nil, fmt.Errorf("uuid: unreachable")
}

func putGenerator(x *fastuuid.Generator) {
	uuidPool.Put(x)
}

// Next returns the next UUID from the generator.
func Next() string {
	g, err := getGenerator()
	if err != nil {
		// FIXME log err
	}
	defer putGenerator(g)

	return fmt.Sprintf("%x", g.Next())
}
