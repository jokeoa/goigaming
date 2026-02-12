type Config struct {
    
    Redis RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host" default:"localhost"`
    Port     string `mapstructure:"port" default:"6379"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db" default:"0"`
}
