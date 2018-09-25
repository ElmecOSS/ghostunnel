// +build windows

/*-
 * Copyright 2018 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/kardianos/service"
	"golang.org/x/sys/windows/svc"
	"os"
)

type program struct{}

var (
	shutdownSignals = []os.Signal{os.Interrupt}
	refreshSignals  = []os.Signal{ /* Not supported on Windows */ }
)

func useSyslog() bool {
	return false
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.runService()
	return nil
}
func (p *program) runService() error {
	err := run(os.Args[1:])
	return err
}
func (p *program) Stop(s service.Service) error {
	return nil
}

func IsInteractive() (bool, error) {
	return svc.IsAnInteractiveSession()
}

func Runner() error {
	svcConfig := &service.Config{
		Name:        "GhostTunnel",
		DisplayName: "Ghost Tunnel Service",
		Description: "SSL/TLS proxy with mutual authentication for securing non TLS services.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	interactive, err := IsInteractive()
	if err != nil {
		exitFunc(1)
	}
	if !interactive {
		err = s.Run()
	} else {
		err = run(os.Args[1:])
	}
	if err != nil {
		exitFunc(1)
	}
	if err != nil {
		return nil
	}
	return nil
}
