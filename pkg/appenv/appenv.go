package appenv

type AppEnv struct {
	Version string
	Env     string
	Port    string
}

func CreateTestAppEnv() AppEnv {
	testVersion := "0.0.0"

	return AppEnv{
		Version: testVersion,
		Env:     "LOCAL",
		Port:    "3001",
	}
}
