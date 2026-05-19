package sum

//go:generate go run ../../../../../main.go

type Common struct {
	ID string `json:"id"`
}

type Hoge struct {
	Common
	Name string `json:"name"`
}

func (h *Hoge) String() string {
	if h == nil {
		return "<nil>"
	}
	return h.Name
}

type Fuga struct {
	Common
	Age int `json:"age"`
}

func (f *Fuga) String() string {
	if f == nil {
		return "<nil>"
	}
	return f.ID
}
