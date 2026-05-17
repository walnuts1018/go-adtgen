package sum

//go:generate go run ../../../../../main.go

type Common struct {
	ID string `json:"id"`
}

type Hoge struct {
	Common
	Name string `json:"name"`
}

type Fuga struct {
	Common
	Age int `json:"age"`
}
