package common

type Role string

const (
	SuperAdmin Role = "superAdmin"
	Admin      Role = "admin"
	Demo       Role = "demo"
)

var AdminRoles = []Role{SuperAdmin, Admin}
