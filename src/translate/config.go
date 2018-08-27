package translate

type Config struct {
	Debug      bool
	ListenPort string
	Languages  map[string]string
	Accounts   []Account
}
