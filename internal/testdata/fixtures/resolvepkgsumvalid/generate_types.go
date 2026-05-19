//go:build adtgen_generate

package resolvepkgsumvalid

import "io"

// +adtgen:sum=Hoge,Fuga
type HogeOrFuga interface {
	WriteTo(io.Writer) (int64, error)
}
