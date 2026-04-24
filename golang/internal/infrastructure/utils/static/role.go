package static

type RoleName string

const (
	Admin RoleName = "admin"
	User  RoleName = "user"
)

var Role = []RoleName{
	Admin,
	User,
}

func IsRole(role RoleName) bool {
	for _, r := range Role {
		if r == role {
			return true
		}
	}
	return false
}