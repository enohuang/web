package web

import "testing"

//import "dengming20240317/orm/internal/errs"

type StringValue3[T any] struct {
	val string
	err error
	a   T
}

func TestContext_BindJSON(t *testing.T) {
	var a StringValue3[int]
	a.a = 2
	t.Log(a.a)

	/*var e error
	if e.Error() ==errs.ErrPointerOnly.Error(){

	}*/
}
