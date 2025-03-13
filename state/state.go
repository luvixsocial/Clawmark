package state

import (
	"context"
	"os"

	"clawmark/config"
	"clawmark/types"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/infinitybotlist/eureka/genconfig"
	"github.com/infinitybotlist/eureka/snippets"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Pool      *gorm.DB
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

	// Initalize Gorm connection
	Pool, err = gorm.Open(postgres.Open(Config.Database.DatabaseURL), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database: %v" + err.Error())
    }

	Pool.AutoMigrate(&types.User{})
	Pool.AutoMigrate(&types.Post{})
	Pool.AutoMigrate(&types.PostPlugin{})
	Pool.AutoMigrate(&types.Comment{})
	Pool.AutoMigrate(&types.Like{})
	Pool.AutoMigrate(&types.Dislike{})
	Pool.AutoMigrate(&types.Follow{})
	
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
