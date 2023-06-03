package env

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

func (e Environment) IsDev() bool {
	return e == Dev || e == "development"
}
