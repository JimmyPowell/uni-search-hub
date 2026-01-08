package constant

var ItemsPerPage = 10

var CommonGroupCol = `"group"`
var CommonKeyCol = `"key"`

var QuotaForNewUser = 0
var QuotaForInviter = 0
var QuotaForInvitee = 0

var BatchUpdateEnabled = false
var BatchUpdateInterval int

const (
	RoleGuestUser  = 0
	RoleCommonUser = 1
	RoleAdminUser  = 10
	RoleRootUser   = 100
)

const (
	UserStatusEnabled  = 1 // don't use 0, 0 is the default value!
	UserStatusDisabled = 2 // also don't use 0
)
