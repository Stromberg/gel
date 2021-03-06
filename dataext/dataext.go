// Package dataext provides functions to fill a store with values based on calculations that
// have dependencies on other data.
package dataext

import "fmt"

// Store is an interface for something storing values indexed by string ids.
type Store interface {
	Get(id string) (value interface{}, ok bool)
	Set(id string, value interface{})
}

// Extender is an interface for something that can calculate the value for an
// id based on values stored in a Store.
type Extender interface {
	ID() string
	String() string
	Missing(store Store) ([]string, error)
	Extend(store Store) error
}

// Extend extends a Store with new values from extenders.
func Extend(store Store, extenders ...Extender) (err error) {
	es := make(map[string]Extender, len(extenders))
	missing := make(map[string]bool)

	for _, e := range extenders {
		es[e.ID()] = e
		ms, err := e.Missing(store)
		if err != nil {
			return err
		}
		if len(ms) == 0 {
			err = e.Extend(store)
			if err != nil {
				return fmt.Errorf("Error extending %s: %s", e.ID(), err)
			}
			delete(missing, e.ID())
		} else {
			missing[e.ID()] = true
			for _, m := range ms {
				missing[m] = true
			}
		}
	}

	for len(missing) > 0 {
		missing2 := make(map[string]bool)

		for k := range missing {
			e, ok := es[k]
			if !ok {
				return fmt.Errorf("Missing extender for %s", k)
			}

			ms, err := e.Missing(store)
			if err != nil {
				return err
			}
			if len(ms) == 0 {
				err = e.Extend(store)
				if err != nil {
					return fmt.Errorf("Error extending %s: %s", e.ID(), err)
				}
				delete(missing2, e.ID())
			} else {
				missing2[e.ID()] = true
				for _, m := range ms {
					missing2[m] = true
				}
			}
		}

		missing = missing2
	}

	return nil
}
