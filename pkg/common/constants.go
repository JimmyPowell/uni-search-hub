package common

import (
	"sync"

	"github.com/google/uuid"
)

var CryptoSecret = uuid.New().String()
var SessionSecret = uuid.New().String()

var Version = "v0.0.0" // this hard coding will be replaced automatically when building, no need to manually change
var SystemName = "UniSearchHub"

var DebugEnabled bool

var ItemsPerPage = 10

var CommonGroupCol = `"group"`
var CommonKeyCol = `"key"`

var QuotaForNewUser = 0
var QuotaForInviter = 0
var QuotaForInvitee = 0

// TODO 批更新检查
var BatchUpdateEnabled = false
var BatchUpdateInterval int

var TurnstileCheckEnabled = false
var TurnstileSiteKey = ""
var TurnstileSecretKey = ""

var RegisterEnabled = true
var PasswordLoginEnabled = true
var PasswordRegisterEnabled = true
var EmailVerificationEnabled = false

const (
	RoleGuestUser  = 0
	RoleCommonUser = 1
	RoleAdminUser  = 10
	RoleRootUser   = 100
)

func IsValidateRole(role int) bool {
	return role == RoleGuestUser || role == RoleCommonUser || role == RoleAdminUser || role == RoleRootUser
}

const (
	UserStatusEnabled  = 1 // don't use 0, 0 is the default value!
	UserStatusDisabled = 2 // also don't use 0
)

var OptionMap map[string]string
var OptionMapRWMutex sync.RWMutex

const (
	TokenStatusEnabled   = 1 // don't use 0, 0 is the default value!
	TokenStatusDisabled  = 2 // also don't use 0
	TokenStatusExpired   = 3
	TokenStatusExhausted = 4
)

var QuotaPerUnit = 500 * 1000.0
