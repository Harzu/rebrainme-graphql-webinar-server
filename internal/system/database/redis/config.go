package redis

type Config struct {
	URL  string
	Pass string `envconfig:"optional"`
}
