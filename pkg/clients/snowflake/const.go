package snowflake

const (
	AccountAdminRole  = "ACCOUNTADMIN"
	SecurityAdminRole = "SECURITYADMIN"
	SysAdminRole      = "SYSADMIN"
	// TODO: create dedicated dice warehouse. integrate creation of dice warehouse into dice init/bootstrapping for new snowflake account
	DiceWarehouse = "DEFAULT"
)
