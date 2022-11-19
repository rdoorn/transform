package transform

import (
	"encoding/json"
	"log"
	"reflect"
)

type Wrapper struct {
	InterfaceType string
	Version       int
	Uuid          string
	I             json.RawMessage
}

type WrapperInterface interface {
	//New(body []byte) error
	Data(t interface{}) (interface{}, error)
	Type() string
	//WithId(s ...string) WrapperInterface
	//Version() int
	//Name() string
	//Id() string
	//Json() string
}

func Wrap(i interface{}) (WrapperInterface, error) {
	w := &Wrapper{
		Version: 1,
		//InterfaceType: fmt.Sprintf("%T", i),
		InterfaceType: TypeName(i),
	}
	customStructure := &customStruct{10, "ten"}
	log.Printf("i: %+v", i)
	j, err := json.Marshal(customStructure)
	if err != nil {
		return nil, err
	}
	log.Printf("j should be: %s", j)
	w.I = json.RawMessage(j)
	log.Printf("in I : %s", w.I)
	return w, nil
}

func (w *Wrapper) Data(t interface{}) (interface{}, error) {

	typi, ok := nameToConcreteType.Load(string(w.InterfaceType))
	if !ok {
		log.Printf("name not registered for interface: %q", w.InterfaceType)
	}
	//typ := &customStruct{}
	customStructure := &customStruct{}
	typ := typi.(reflect.Type)
	log.Printf("typ: %v (%T)", typ, typ)
	//newObj := reflect.New(reflect.TypeOf(typ).Elem())
	//newObj := reflect.New(reflect.TypeOf(typ).Elem())
	//newObj := reflect.New(typ).Interface()
	newObj := reflect.New(reflect.ValueOf(customStructure).Elem().Type()).Interface().(*customStruct)
	log.Printf("new: %v (%T)", newObj, newObj)

	log.Printf("in I : %s", w.I)

	err := json.Unmarshal(w.I, newObj)
	log.Printf("s: %+v", newObj)
	return newObj, err
}

func (w *Wrapper) Type() string {
	return w.InterfaceType
}
