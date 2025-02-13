// Copyright (C) 2023 Opsmate, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla
// Public License, v. 2.0. If a copy of the MPL was not distributed
// with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This software is distributed WITHOUT A WARRANTY OF ANY KIND.
// See the Mozilla Public License for details.

package monitor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stdoutMu sync.Mutex

type notification struct {
	environ []string
	summary string
	text    string
	json    []zap.Field
	level   zapcore.Level
	msg     string
}

func (s *FilesystemState) notify(ctx context.Context, notif *notification) error {
	if s.Stdout && !s.Json {
		writeToStdout(notif)
	} else if s.Json {
		writeJsonToStdout(notif)
	}

	if len(s.Email) > 0 {
		if err := sendEmail(ctx, s.Email, notif); err != nil {
			return err
		}
	}

	if s.Script != "" {
		if err := execScript(ctx, s.Script, notif); err != nil {
			return err
		}
	}

	if s.ScriptDir != "" {
		if err := execScriptDir(ctx, s.ScriptDir, notif); err != nil {
			return err
		}
	}

	return nil
}
func writeJsonToStdout(notif *notification) {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()
	switch notif.level {
	case zapcore.WarnLevel:
		zap.L().Warn(notif.msg, notif.json...)
	default:
		zap.L().Info(notif.msg, notif.json...)
	}

}

func writeToStdout(notif *notification) {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()
	os.Stdout.WriteString(notif.text + "\n")
}

func sendEmail(ctx context.Context, to []string, notif *notification) error {
	stdin := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	from := os.Getenv("EMAIL")

	if from != "" {
		fmt.Fprintf(stdin, "From: %s\n", from)
	}
	fmt.Fprintf(stdin, "To: %s\n", strings.Join(to, ", "))
	fmt.Fprintf(stdin, "Subject: [certspotter] %s\n", notif.summary)
	fmt.Fprintf(stdin, "Date: %s\n", time.Now().Format(mailDateFormat))
	fmt.Fprintf(stdin, "Message-ID: <%s>\n", generateMessageID())
	fmt.Fprintf(stdin, "Mime-Version: 1.0\n")
	fmt.Fprintf(stdin, "Content-Type: text/plain; charset=US-ASCII\n")
	fmt.Fprintf(stdin, "X-Mailer: certspotter\n")
	fmt.Fprintf(stdin, "\n")
	fmt.Fprint(stdin, notif.text)

	args := []string{"-i"}
	if from != "" {
		args = append(args, "-f", from)
	}
	args = append(args, "--")
	args = append(args, to...)

	sendmail := exec.CommandContext(ctx, sendmailPath(), args...)
	sendmail.Stdin = stdin
	sendmail.Stderr = stderr

	if err := sendmail.Run(); err == nil {
		return nil
	} else if ctx.Err() != nil {
		return ctx.Err()
	} else if exitErr, isExitError := err.(*exec.ExitError); isExitError && exitErr.Exited() {
		return fmt.Errorf("error sending email to %v: sendmail failed with exit code %d and error %q", to, exitErr.ExitCode(), strings.TrimSpace(stderr.String()))
	} else {
		return fmt.Errorf("error sending email to %v: %w", to, err)
	}
}

func execScript(ctx context.Context, scriptName string, notif *notification) error {
	stderr := new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, scriptName)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, notif.environ...)
	cmd.Stderr = stderr

	if err := cmd.Run(); err == nil {
		return nil
	} else if ctx.Err() != nil {
		return ctx.Err()
	} else if exitErr, isExitError := err.(*exec.ExitError); isExitError && exitErr.Exited() {
		return fmt.Errorf("script %q exited with code %d and error %q", scriptName, exitErr.ExitCode(), strings.TrimSpace(stderr.String()))
	} else if isExitError {
		return fmt.Errorf("script %q terminated by signal with error %q", scriptName, strings.TrimSpace(stderr.String()))
	} else {
		return fmt.Errorf("error executing script: %w", err)
	}
}

func execScriptDir(ctx context.Context, dirPath string, notif *notification) error {
	dirents, err := os.ReadDir(dirPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("error executing scripts in directory %q: %w", dirPath, err)
	}
	for _, dirent := range dirents {
		if strings.HasPrefix(dirent.Name(), ".") {
			continue
		}
		scriptPath := filepath.Join(dirPath, dirent.Name())
		info, err := os.Stat(scriptPath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return fmt.Errorf("error executing %q in directory %q: %w", dirent.Name(), dirPath, err)
		} else if info.Mode().IsRegular() && isExecutable(info.Mode()) {
			if err := execScript(ctx, scriptPath, notif); err != nil {
				return err
			}
		}
	}
	return nil
}

func isExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}
