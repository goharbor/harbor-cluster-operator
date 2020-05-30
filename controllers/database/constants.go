package database

const (
	CheckDatabaseHealthError          = "Check database health error"
	CreateDatabaseCrError             = "Create database CR error"
	UpdateDatabaseCrError             = "Update database CR error"
	GenerateRedisCrError              = "Generate database CR error"
	SetOwnerReferenceError            = "Set owner reference error"
	DefaultUnstructuredConverterError = "Default unstructured converter error"
)

const (
	DownScalingDatabase     = "DatabaseDownScaling"
	UpScalingDatabase       = "DatabaseUpScaling"
	RollingUpgradesDatabase = "DatabaseRollingUpgrades"

	MessageDatabaseCreate = "Database  %s already created."

	MessageDatabaseUpdate = "Database  %s already update."

	MessageDatabaseDownScaling     = "Database downscale from %d to %d"
	MessageDatabaseUpScaling       = "Database upscale from %d to %d"
	MessageDatabaseRollingUpgrades = "Database resource from %s to %s"
)

const (
	InClusterDatabasePort        = "5432"
	InClusterDatabaseUserName    = "postgres"
	InClusterDatabaseName        = "postgres"
	InClusterDatabasePasswordKey = "password"
)
