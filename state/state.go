package state

import (
	"context"
	"os"

	"clawmark/config"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/infinitybotlist/eureka/genconfig"
	"github.com/infinitybotlist/eureka/snippets"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	Pool      *pgxpool.Pool
	Redis     *redis.Client
	Logger    *zap.Logger
	Context   = context.Background()
	Validator = validator.New()
	Config    *config.Config
)

func Setup() {
	Validator.RegisterValidation("notblank", validators.NotBlank)
	Validator.RegisterValidation("nospaces", snippets.ValidatorNoSpaces)
	Validator.RegisterValidation("https", snippets.ValidatorIsHttps)
	Validator.RegisterValidation("httporhttps", snippets.ValidatorIsHttpOrHttps)

	genconfig.GenConfig(config.Config{})

	cfg, err := os.ReadFile("config.yaml")
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	err = yaml.Unmarshal(cfg, &Config)
	if err != nil {
		panic("Failed to parse config file: " + err.Error())
	}

	err = Validator.Struct(Config)
	if err != nil {
		panic("config validation error: " + err.Error())
	}

	// Initialize PostgreSQL connection
	Pool, err = pgxpool.New(Context, Config.Database.DatabaseURL)
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}

	// Initialize Redis connection
	rOptions, err := redis.ParseURL(Config.Database.RedisURL)
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}

	Redis = redis.NewClient(rOptions)
	if err := Redis.Ping(Context).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	// Initialize Logger
	Logger = snippets.CreateZap()
}
