// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package execute

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Command      string
	Args         []string
	Shell        bool
	StreamStdio  bool
	PrintCommand bool
}

type Response struct {
	Stdout   string
	Stderr 	 string
	ExitCode int
}

func (t Task) Execute() (Response, error) {
	var cmd *exec.Cmd

	if t.PrintCommand {
		fmt.Printf("[+] %s %s\n", t.Command, t.Args)
	}

	if t.Shell {
		var args []string
		args = append(args, t.Command)

		if len(t.Args) > 0 {
			args = append(args, t.Args...)
		}

		cmd = exec.Command("/bin/bash", args...)
	} else {
		cmd = exec.Command(t.Command, t.Args...)
	}

	stdoutBuff := bytes.Buffer{}
	stderrBuff := bytes.Buffer{}

	var stdoutWriters io.Writer
	var stderrWriters io.Writer

	if t.StreamStdio {
		stdoutWriters = io.MultiWriter(os.Stdout, &stdoutBuff)
		stderrWriters = io.MultiWriter(os.Stderr, &stderrBuff)
	} else {
		stdoutWriters = &stdoutBuff
		stderrWriters = &stderrBuff
	}

	cmd.Stdout = stdoutWriters
	cmd.Stderr = stderrWriters

	if err := cmd.Start(); err != nil {
		return Response{}, err
	}

	exitCode := 0
	if execErr := cmd.Wait(); execErr != nil {
		if exitError, ok := execErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	err := string(stderrBuff.Bytes())
	index := strings.LastIndex(err, "msg=")

	var msgError string
	if index != -1 {
		length := len(err)
		msgError = "\n" + err[index+5:length-2] + "\n"
	} else {
		msgError = err
	}

	return Response{
		Stdout: string(stdoutBuff.Bytes()),
		Stderr: msgError,
		ExitCode: exitCode,
	}, nil
}
