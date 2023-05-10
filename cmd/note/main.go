package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	apphttp "github.com/aasenknut/note/http"
	"github.com/aasenknut/note/redis"
	"github.com/aasenknut/note/sqlite"
	"github.com/pelletier/go-toml/v2"
)

var migrate = flag.Bool("migrate", false, "`true` to run migrations")
var confPath = flag.String("config", "./conf.toml", "path to .toml config")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Shutting down...")
		cancel()
	}()
	flag.Parse()
	app := NewApp()
	var err error
	if app.Config, err = ParseConfigFile(*confPath); err != nil {
		log.Println("read config file")
	}

	if err := app.Run(ctx); err != nil {
		log.Printf("closing app: %v", err)
		app.Close()
		os.Exit(1)
	}

	if *migrate {
		log.Println("migrating...")
		if err = app.DB.Migrate(); err != nil {
			log.Printf("migration: %v", err)
			cancel()
		}
	}

	<-ctx.Done()
	if err := app.Close(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type App struct {
	ConfigPath string
	Config     *Config
	DB         *sqlite.DB
	MemStore   *redis.DB
	HTTPServer *apphttp.Server
}

func NewApp() *App {
	return &App{
		DB:         sqlite.NewDB(),
		HTTPServer: apphttp.NewServer(),
		MemStore:   redis.NewDB(),
	}
}

func (a *App) Close() error {
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	// Set configs
	a.SetServerConf()
	a.SetDBConf()
	a.SetMemstoreConf()

	// Open services
	a.DB.DSN = a.Config.DB.DSN
	a.DB.Open()
	a.MemStore.Addr = a.Config.MemStore.Addr
	a.MemStore.Password = a.Config.MemStore.Password
	a.MemStore.DBIndex = a.Config.MemStore.DBIndex
	a.MemStore.Open()
	a.HTTPServer.NoteService = sqlite.NewNoteService(a.DB)
	a.HTTPServer.UserService = sqlite.NewUserService(a.DB)
	a.HTTPServer.AuthService = redis.NewAuthService(a.MemStore)
	err := a.HTTPServer.SetTmplCache()
	if err != nil {
		log.Printf("could not setup cache: %v", err)
	}
	a.HTTPServer.Open()
	return nil
}

func (a *App) SetServerConf() {
	a.HTTPServer.Addr = a.Config.HTTP.Addr
}

func (a *App) SetDBConf() {
	a.DB.DSN = a.Config.DB.DSN
}

func (a *App) SetMemstoreConf() {
	a.MemStore.Addr = a.Config.MemStore.Addr
	a.MemStore.Password = a.Config.MemStore.Password
	a.MemStore.DBIndex = a.Config.MemStore.DBIndex
}

type Config struct {
	DB       DBConfig       `toml:"db"`
	HTTP     HTTPConfig     `toml:"http"`
	MemStore MemStoreConfig `toml:"memstore"`
}
type DBConfig struct {
	DSN string `toml:"dsn"`
}
type HTTPConfig struct {
	Addr string `toml:"addr"`
}

type MemStoreConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DBIndex  int    `toml:"dbindex"`
}

func ParseConfigFile(fname string) (*Config, error) {
	var conf Config
	buf, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	if err := toml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
