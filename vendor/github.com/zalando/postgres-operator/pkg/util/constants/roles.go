package constants

// Roles specific constants
const (
	PasswordLength            = 64
	SuperuserKeyName          = "superuser"
	ConnectionPoolerUserKeyName = "pooler"
	ReplicationUserKeyName    = "replication"
	RoleFlagSuperuser         = "SUPERUSER"
	RoleFlagInherit           = "INHERIT"
	RoleFlagLogin             = "LOGIN"
	RoleFlagNoLogin           = "NOLOGIN"
	RoleFlagCreateRole        = "CREATEROLE"
	RoleFlagCreateDB          = "CREATEDB"
	RoleFlagReplication       = "REPLICATION"
	RoleFlagByPassRLS         = "BYPASSRLS"
	OwnerRoleNameSuffix       = "_owner"
	ReaderRoleNameSuffix      = "_reader"
	WriterRoleNameSuffix      = "_writer"
	UserRoleNameSuffix        = "_user"
)
