package env

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

func (e Environment) IsDev() bool {
	return e == Dev || e == "development"
}

func (e Environment) IsProd() bool {
	return e == Prod || e == "production"
}

func (e Environment) NonProd() bool {
	return !e.IsProd()
}
