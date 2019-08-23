package main

import (
	"flag"
	"fmt"
	l "log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	guerrilla "github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/log"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

const (
	InsertIP      = "INSERT INTO ips (ip) values ($1) ON CONFLICT DO NOTHING;"
	InsertAddress = "INSERT INTO addresses (address) values ($1) ON CONFLICT DO NOTHING;"
	InsertDomain  = "INSERT INTO domains (domain) values ($1) ON CONFLICT DO NOTHING;"
)

var funkyLogger = func() backends.Decorator {
	backends.Svc.AddInitializer(
		backends.InitializeWith(
			func(backendConfig backends.BackendConfig) error {
				backends.Log().Info("Funky logger is up & down to funk!")
				return nil
			}),
	)

	backends.Svc.AddShutdowner(
		backends.ShutdownWith(
			func() error {
				backends.Log().Info("The funk has been stopped!")
				return nil
			}),
	)

	return func(p backends.Processor) backends.Processor {
		return backends.ProcessWith(
			func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
				if task == backends.TaskSaveMail {
					_, err := db.Exec(InsertIP, e.RemoteIP)
					if err != nil {
						backends.Log().Error(err)
					}

					_, err = db.Exec(InsertAddress, e.MailFrom.String())
					if err != nil {
						backends.Log().Error(err)
					}

					domain := strings.Split(e.MailFrom.String(), "@")[1]
					_, err = db.Exec(InsertDomain, domain)
					if err != nil {
						backends.Log().Error(err)
					}
				}
				return p.Process(e, task)
			})
	}
}

var (
	connectionString string
)

func main() {
	flag.StringVar(&connectionString, "db", os.Getenv("COCKROACHDB_URL"), "connection string postgres like database (required)")
	flag.Parse()

	if strings.TrimSpace(connectionString) == "" {
		flag.Usage()
		return
	}

	var err error
	db, err = sqlx.Open("postgres", connectionString)
	if err != nil {
		l.Fatal(err)
	}

	config := &guerrilla.AppConfig{
		LogFile:      log.OutputStderr.String(),
		AllowedHosts: []string{"."},
		BackendConfig: backends.BackendConfig{
			"save_process": "HeadersParser|Debugger|FunkyLogger",
		},
	}

	sc := guerrilla.ServerConfig{
		ListenInterface: "0.0.0.0:2525",
		IsEnabled:       true,
	}

	config.Servers = append(config.Servers, sc)

	d := guerrilla.Daemon{Config: config}
	d.AddProcessor("FunkyLogger", funkyLogger)
	err = d.Start()

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGUSR1,
		os.Kill,
	)
	_ = <-signalChannel
	go func() {
		select {
		// exit if graceful shutdown not finished in 60 sec.
		case <-time.After(time.Second * 60):
			fmt.Println("graceful shutdown timed out")
			os.Exit(1)
		}
	}()
	d.Shutdown()
}
