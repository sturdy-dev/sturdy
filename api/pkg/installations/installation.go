package installations

type Type uint

const (
	TypeUndefined Type = iota
	TypeOSS
	TypeEnterprise
	TypeCloud
)

func (t Type) String() string {
	switch t {
	case TypeOSS:
		return "oss"
	case TypeEnterprise:
		return "enterprise"
	case TypeCloud:
		return "cloud"
	default:
		return "undefined"
	}
}

// Installation represents a selfhosted installation of Sturdy.
type Installation struct {
	ID      string `db:"id"`
	Type    Type   `db:"-"`
	Version string `db:"-"`
}
