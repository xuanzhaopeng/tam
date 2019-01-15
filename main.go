package main

import (
	"flag"
	"os"
	"syscall"
	"os/signal"
	"net/http"
	"log"
	"context"
	"time"
	"tam/account"
	"io/ioutil"
)

var (
	accountJsonFilePath string
	accountKey string
	listen string
	accounts *account.Accounts
	accountTimeout time.Duration
)

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	flag.StringVar(&accountJsonFilePath, "accounts", "accounts.json", "The account json file path")
	flag.StringVar(&accountKey, "key", "username", "The key attribute name of accounts structure")
	flag.StringVar(&listen, "listen", "localhost:6666", "The host and port that listen to")
	flag.DurationVar(&accountTimeout, "timeout", 300*time.Second, "account timeout in time.Duration format, e.g. 300s or 500ms")
	flag.Parse()

	log.SetFlags(0)
	log.Printf("accounts: %s\n", accountJsonFilePath)
	log.Printf("key: %s\n", accountKey)
	log.Printf("timeout: %s\n", accountTimeout)

	loadAccounts()
}

func loadAccounts() {
	data, err := ioutil.ReadFile(accountJsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
	accounts, err = account.BuildAccounts(accountKey, accountTimeout, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	gracefulPeriod := 10 * time.Second

	server := &http.Server{
		Addr:    listen,
		Handler: mux(),
	}

	e := make(chan error)
	go func() {
		log.Printf("start server at %s...\n", listen)
		e <- server.ListenAndServe()
	}()
	select {
	case err := <-e:
		log.Fatalf("[-] [-] [INIT] [-] [-] [-] [-] [-] [-] [%v]\n", err)
	case <-stop:
	}

	log.Printf("[-] [%s] [SHUTTING_DOWN] [-] [-] [-] [-] [-] [-] [-]\n", gracefulPeriod)
	ctx, cancel := context.WithTimeout(context.Background(), gracefulPeriod)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[-] [-] [SHUTDOWN_FAILURE] [-] [-] [-] [-] [-] [-] [%v]\n", err)
	}
}
