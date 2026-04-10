package example

const P = 100

const p = `1
2
3`

func q() string {
	a := "OK2"
	return a + "OK3"
}

type E struct{}

// simple New

func New() E {
	return E{}
}

func NewP() *E {
	return &E{}
}

func NewE() (*E, error) {
	return &E{}, nil
}

// factories

type T func() E

type I interface {
	F() E
}

//

func (E) A() {}

func (_ E) a() {}

//

type EE struct{}

func (EE) one() {}

//

type ee struct{}

func (ee) one() {}
