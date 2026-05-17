package resolvepkgsumvalid

type Common struct {
	ID string
}

type Hoge struct {
	Common
	Name string
}

type Fuga struct {
	Common
	Age int
}
