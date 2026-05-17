//go:build goproducttype_generate

package resolvepkg

import timex "time"

type LocalTime = timex.Time

//goproducttype:product Customer Address
type CustomerAddress struct{}

//goproducttype:product Customer Envelope[string]
type CustomerEnvelope struct{}

//goproducttype:product Customer timex.Time
type CustomerTime struct{}

//goproducttype:product Customer LocalTime
type CustomerLocalTime struct{}

//goproducttype:product Customer Envelope[T]
type CustomerEnvelopeForTypeParam[T any] struct{}
