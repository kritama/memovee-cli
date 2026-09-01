package command

type Category string

const (
	CategoryUsage         Category = "usage"
	CategoryPrerequisite  Category = "prerequisite"
	CategoryContract      Category = "contract"
	CategoryConfiguration Category = "configuration"
	CategoryOwnership     Category = "ownership"
	CategorySecret        Category = "secret"
	CategoryProcess       Category = "process"
	CategoryNetwork       Category = "network"
	CategoryVerification  Category = "verification"
	CategoryActivation    Category = "activation"
	CategoryRollback      Category = "rollback"
	CategoryInternal      Category = "internal"
)

const (
	ExitSuccess       = 0
	ExitUsage         = 2
	ExitPrerequisite  = 10
	ExitContract      = 11
	ExitConfiguration = 12
	ExitOwnership     = 13
	ExitSecret        = 14
	ExitProcess       = 15
	ExitNetwork       = 16
	ExitVerification  = 17
	ExitActivation    = 18
	ExitRollback      = 19
	ExitInternal      = 20
)

type Problem struct {
	Category Category
	Message  string
	Next     string
}

func (p Problem) Error() string {
	return p.Message
}

func ExitCode(category Category) int {
	switch category {
	case CategoryUsage:
		return ExitUsage
	case CategoryPrerequisite:
		return ExitPrerequisite
	case CategoryContract:
		return ExitContract
	case CategoryConfiguration:
		return ExitConfiguration
	case CategoryOwnership:
		return ExitOwnership
	case CategorySecret:
		return ExitSecret
	case CategoryProcess:
		return ExitProcess
	case CategoryNetwork:
		return ExitNetwork
	case CategoryVerification:
		return ExitVerification
	case CategoryActivation:
		return ExitActivation
	case CategoryRollback:
		return ExitRollback
	case CategoryInternal:
		return ExitInternal
	default:
		return ExitInternal
	}
}
