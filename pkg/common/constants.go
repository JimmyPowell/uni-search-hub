package common

import (
	"github.com/google/uuid"
)

var CryptoSecret = uuid.New().String()

var Version = "v0.0.0" // this hard coding will be replaced automatically when building, no need to manually change
var SystemName = "UniSearchHub"

var DebugEnabled bool

var SyncFrequency int // unit is second
