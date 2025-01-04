package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dev79844/bitcask"
	"github.com/tidwall/redcon"
)

type App struct {
	l 		*slog.Logger
	bitcask *bitcask.Bitcask
}

func initLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}

func (app *App) ping(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

func (app *App) get(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("wrong number of arguments")
		return
	}

	key := string(cmd.Args[1])

	val, err := app.bitcask.Get(key)
	fmt.Println("value:", val)
	if err!=nil{
		conn.WriteError(fmt.Sprintf("ERR: %s", err))
		return
	}

	conn.WriteBulk(val)
}

func (app *App) set(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 3 {
		conn.WriteError("wrong number of arguments")
		return
	}

	var(
		key = string(cmd.Args[1])
		value = cmd.Args[2]
	)

	err := app.bitcask.Put(key, value)
	if err!=nil{
		conn.WriteError(fmt.Sprintf("ERR: %s", err))
		return
	}

	conn.WriteString("OK")
}

func (app *App) delete(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("wrong number of arguments")
		return
	}

	key := string(cmd.Args[1])

	err := app.bitcask.Delete(key)
	if err!=nil{
		conn.WriteError(fmt.Sprintf("ERR: %s", err))
		return
	}

	conn.WriteNull()
}

func (app *App) exit(conn redcon.Conn, cmd redcon.Command) {
	app.l.Info("closing the connection")
	conn.WriteString("OK")
	conn.Close()
}


func main() {
	app := &App{
		l: initLogger(),
	}

	cfg := []bitcask.Config{
		bitcask.WithDir("cmd/server/data"),
		bitcask.WithMaxActiveFileSize(4096),
	}

	//initialise a bitcask instance
	db, err := bitcask.Open(cfg...)
	if err!=nil{
		app.l.Error("error initialing bitcask instance", slog.Any("err", err))
	}

	app.bitcask = db

	mux := redcon.NewServeMux()

	mux.HandleFunc("ping", app.ping)
	mux.HandleFunc("get", app.get)
	mux.HandleFunc("set", app.set)
	mux.HandleFunc("del", app.delete)
	mux.HandleFunc("exit", app.exit)

	//for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	srv := redcon.NewServer(":6379", 
		mux.ServeRESP,
		func(conn redcon.Conn) bool {
			// use this function to accept or deny the connection.
			log.Printf("accept: %s", conn.RemoteAddr())
			return true 
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
			log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)

	go func(){
		err := srv.ListenAndServe()
		if err!=nil{
			app.l.Error("error starting the server")
		}
	}()

	<-ctx.Done()

	cancel()
	app.bitcask.Close()
	srv.Close()
}