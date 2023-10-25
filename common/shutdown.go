// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/osintami/fingerprintz/log"
)

type IShutdown interface {
	Listen()
	AddListener(f func())
}

type ShutdownHandler struct {
	listeners     []func()
	signalChannel chan os.Signal
}

func NewShutdownHandler() *ShutdownHandler {
	return &ShutdownHandler{}
}

func (x *ShutdownHandler) AddListener(f func()) {
	x.listeners = append(x.listeners, f)
}

func (x *ShutdownHandler) Listen(Exit func()) {
	x.signalChannel = make(chan os.Signal, 1)
	signal.Notify(x.signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-x.signalChannel
		log.Info().Str("component", "shutdown").Str("item", s.String()).Msg("item nofify")
		for _, ShutItDown := range x.listeners {
			ShutItDown()
		}
		Exit()
	}()
}

func (x *ShutdownHandler) Interrupt() {
	x.signalChannel <- syscall.SIGINT
}
