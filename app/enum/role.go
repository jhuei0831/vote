package enum

type Role string

const (
    Admin   Role = "ADMIN"
    Creator Role = "CREATOR"
    Voter   Role = "VOTER"
)