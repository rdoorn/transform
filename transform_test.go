package transform

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type customStruct struct {
	Int    int
	String string
}

func TestEncodeDecode(t *testing.T) {
	Register(&customStruct{})
	a := assert.New(t)
	customStructure := &customStruct{10, "ten"}
	log.Printf("before: %+v", customStructure)

	wrapper, err := Wrap(customStructure)
	if err != nil {
		a.Error(err)
	}
	log.Printf("wrapped: %+v", wrapper)

	newStructType := wrapper.Type()
	log.Printf("type: %s", newStructType)
	newStruct := &customStruct{}
	z, err := wrapper.Data(newStruct)
	log.Printf("dta: %+v err:%s", z.(*customStruct).Int, err)
}
