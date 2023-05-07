package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	apphttp "github.com/aasenknut/note/http"
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
	app := NewApp()
	var err error
	if app.Config, err = ParseConfigFile(*confPath); err != nil {
		log.Println("read config file")
	}

	if err := app.Run(ctx); err != nil {
		app.Close()
		log.Println("run app")
		os.Exit(1)
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
	HTTPServer *apphttp.Server
}

func NewApp() *App {
	return &App{
		DB:         sqlite.NewDB(),
		HTTPServer: apphttp.NewServer(),
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

	// Open services
	a.DB.DSN = a.Config.DB.DSN
	a.DB.Open()
	a.HTTPServer.NoteService = sqlite.NewNoteService(a.DB)
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

type Config struct {
	DB   DBConfig   `toml:"db"`
	HTTP HTTPConfig `toml:"http"`
}
type DBConfig struct {
	DSN string `toml:"dsn"`
}
type HTTPConfig struct {
	Addr string `toml:"addr"`
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
