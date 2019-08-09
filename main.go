package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	guerrilla "github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/log"
	"github.com/flashmob/go-guerrilla/mail"
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
				if task == backends.TaskValidateRcpt {
					backends.Log().Infof(
						"another funky sender [%s]",
						e.MailFrom.String())
					// log the last recipient appended to e.Rcpt
					backends.Log().Infof(
						"another funky recipient [%s]",
						e.RcptTo[len(e.RcptTo)-1].String())
					// if valid then forward call to the next processor in the chain
					return p.Process(e, task)
				} else if task == backends.TaskSaveMail {
					backends.Log().Infof("RemoteIP: %s", e.RemoteIP)
					domain := strings.Split(e.MailFrom.String(), "@")[1]
					backends.Log().Infof("From: %s", e.MailFrom.String())
					backends.Log().Infof("Domain: %s", domain)
					for _, rcpt := range e.RcptTo {
						backends.Log().Infof("To: %s", rcpt.String())
					}
					for k, v := range e.Header {
						backends.Log().Infof("%s: %s", k, v)
					}
					backends.Log().Info(e.String())
				}
				return p.Process(e, task)
			})
	}
}

func main() {
	config := &guerrilla.AppConfig{
		LogFile:      log.OutputStderr.String(),
		AllowedHosts: []string{"."},
		BackendConfig: backends.BackendConfig{
			"save_process":     "HeadersParser|Debugger|FunkyLogger",
			"validate_process": "FunkyLogger",
		},
	}
	d := guerrilla.Daemon{Config: config}
	d.AddProcessor("FunkyLogger", funkyLogger)
	err := d.Start()

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
