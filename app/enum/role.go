package enum

var Roles = newRoleRegistry()

func newRoleRegistry() *roleRegistry {
    return &roleRegistry{
        Admin:   "ADMIN",
        Creator: "CREATOR",
        Anon:    "ANON",
    }
}

type roleRegistry struct {
    Admin   string
    Creator string
    Anon    string
}