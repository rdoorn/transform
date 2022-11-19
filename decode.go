package transform

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

type Decoder struct {
}

var (
	nameToConcreteType sync.Map // map[string]reflect.Type
	concreteTypeToName sync.Map // map[reflect.Type]string
)

// RegisterName is like Register but uses the provided name rather than the
// type's default.
func RegisterName(name string, value any) {
	if name == "" {
		// reserved for nil
		panic("attempt to register empty name")
	}

	ut := reflect.TypeOf(value)

	// Check for incompatible duplicates. The name must refer to the
	// same user type, and vice versa.

	// Store the name and type provided by the user....
	log.Printf("store1 %v (%T), %v (%T)", name, name, reflect.TypeOf(value), reflect.TypeOf(value))
	if t, dup := nameToConcreteType.LoadOrStore(name, reflect.TypeOf(value)); dup && t != ut {
		panic(fmt.Sprintf("gob: registering duplicate types for %q: %s != %s", name, t, ut))
	}

	// but the flattened type in the type table, since that's what decode needs.
	log.Printf("store2 %v (%T), %v (%T)", ut, ut, name, name)
	if n, dup := concreteTypeToName.LoadOrStore(ut, name); dup && n != name {
		nameToConcreteType.Delete(name)
		panic(fmt.Sprintf("gob: registering duplicate names for %s: %q != %q", ut, n, name))
	}
	log.Printf("ntt1: %+v", nameToConcreteType)
	log.Printf("ntt2: %+v", concreteTypeToName)
	log.Printf("name registered")
}

func TypeName(value any) string {
	rt := reflect.TypeOf(value)
	name := rt.String()

	// But for named types (or pointers to them), qualify with import path (but see inner comment).
	// Dereference one pointer looking for a named type.
	star := ""
	if rt.Name() == "" {
		if pt := rt; pt.Kind() == reflect.Pointer {
			star = "*"
			// NOTE: The following line should be rt = pt.Elem() to implement
			// what the comment above claims, but fixing it would break compatibility
			// with existing gobs.
			//
			// Given package p imported as "full/p" with these definitions:
			//     package p
			//     type T1 struct { ... }
			// this table shows the intended and actual strings used by gob to
			// name the types:
			//
			// Type      Correct string     Actual string
			//
			// T1        full/p.T1          full/p.T1
			// *T1       *full/p.T1         *p.T1
			//
			// The missing full path cannot be fixed without breaking existing gob decoders.
			rt = pt.Elem()
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			name = star + rt.Name()
		} else {
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}
	return name
}

func Register(value any) {
	// Default to printed representation for unnamed types
	rt := reflect.TypeOf(value)
	name := rt.String()

	// But for named types (or pointers to them), qualify with import path (but see inner comment).
	// Dereference one pointer looking for a named type.
	star := ""
	if rt.Name() == "" {
		if pt := rt; pt.Kind() == reflect.Pointer {
			star = "*"
			// NOTE: The following line should be rt = pt.Elem() to implement
			// what the comment above claims, but fixing it would break compatibility
			// with existing gobs.
			//
			// Given package p imported as "full/p" with these definitions:
			//     package p
			//     type T1 struct { ... }
			// this table shows the intended and actual strings used by gob to
			// name the types:
			//
			// Type      Correct string     Actual string
			//
			// T1        full/p.T1          full/p.T1
			// *T1       *full/p.T1         *p.T1
			//
			// The missing full path cannot be fixed without breaking existing gob decoders.
			rt = pt.Elem()
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			name = star + rt.Name()
		} else {
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}

	//log.Printf("name:%v (%T), value:%v (%T)", name, name, value, value)
	RegisterName(name, value)
}
