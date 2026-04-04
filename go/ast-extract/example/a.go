package example

type E struct{}

func (E) A() {}

func (_ E) a() {}

//

type EE struct{}

func (EE) one() {}

//

type ee struct{}

func (ee) one() {}
