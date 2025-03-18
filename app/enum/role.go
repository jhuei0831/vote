package enum

var Roles = newRoleRegistry()

func newRoleRegistry() *roleRegistry {
    return &roleRegistry{
        Admin:   "ADMIN",
        Creator: "CREATOR",
    }
}

type roleRegistry struct {
    Admin   string
    Creator string
}