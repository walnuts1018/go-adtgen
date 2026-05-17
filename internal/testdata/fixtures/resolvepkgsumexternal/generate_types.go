//go:build goproducttype_generate

package resolvepkgsumexternal

import "time"

var _ time.Time

//goproducttype:sum Hoge time.Time
type HogeOrTime struct{}
