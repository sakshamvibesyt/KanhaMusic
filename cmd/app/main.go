/*
 * ● KanhaMusic
 * ○ A high-performance engine for streaming music in Telegram voicechats.
 *
 * Copyright (C) 2026 Kanha
 *
 * This program is free software: you can redistribute it and/or modify it under the
 * terms of the GNU General Public License as published by the Free Software Foundation,
 * either version 3 of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
 * PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * Repository: https://github.com/Oyekanhaa/KanhaMusic
 */

package main

/*
#cgo CFLAGS: -I../../
#cgo linux LDFLAGS: -L ../../ -lntgcalls -lm -lz
#cgo darwin LDFLAGS: -L ../../ -lntgcalls -lc++ -lz -lbz2 -liconv -framework AVFoundation -framework AudioToolbox -framework CoreAudio -framework QuartzCore -framework CoreMedia -framework VideoToolbox -framework AppKit -framework Metal -framework MetalKit -framework OpenGL -framework IOSurface -framework ScreenCaptureKit

// Currently is supported only dynamically linked library on Windows due to
// https://github.com/golang/go/issues/63903
#cgo windows LDFLAGS: -L../../ -lntgcalls
#include "ntgcalls/ntgcalls.h"
#include "glibc_compatibility.h"
*/
import "C"

import (
	"io"
	"net/http"
	"os"

	"KanhaMusic/config"
	"KanhaMusic/kanha/core"
	"KanhaMusic/kanha/database"
	"KanhaMusic/kanha/locales"
	"KanhaMusic/kanha/logger"
	"KanhaMusic/kanha/modules"
)

func main() {
	f, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		logger.Fatalf("failed to open logs.txt: %v", err)
	}
	defer f.Close()

	logger.SetOutput(io.MultiWriter(os.Stderr, f))

	if err := config.Load(); err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	if err := checkDependencies(); err != nil {
		logger.Fatalf("dependency check failed: %v", err)
	}

	if err := refreshDirs(); err != nil {
		logger.Fatalf("failed to refresh directories: %v", err)
	}

	logger.Debug("Initializing MongoDB...")

	closeDB, err := database.Init(config.MongoURI)
	if err != nil {
		logger.Fatalf("failed to initialize database: %v", err)
	}
	defer closeDB()

	logger.Info("Database connected successfully")

	if err := locales.Load(); err != nil {
		logger.Fatalf("failed to load locales: %v", err)
	}

	logger.Debug("Initializing clients...")

	shutdownCore, err := core.Init()
	if err != nil {
		logger.Fatalf("failed to initialize core: %v", err)
	}
	defer shutdownCore()

	core.F = modules.F

	if err := database.InitAssistantIndexes(core.Assistants.Count()); err != nil {
		logger.Fatalf("failed to init assistants: %v", err)
	}

	modules.Init(core.Bot, core.Assistants)

	startHTTPServer()

	core.Bot.Idle()
}


func startHTTPServer() {
	go func() {
		addr := "0.0.0.0:" + config.Port
		mux := http.NewServeMux()

		// UptimeRobot health-check endpoint
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "KanhaMusic is running")
		})

		// Homepage
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}

			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			_, _ = io.WriteString(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>KanhaMusic</title>
<style>
body {
	margin: 0;
	min-height: 100vh;
	display: grid;
	place-items: center;
	background: #100d1c;
	color: white;
	font-family: Arial, sans-serif;
}
.card {
	text-align: center;
	padding: 35px;
	border: 1px solid #39304e;
	border-radius: 20px;
	background: #1c172b;
	max-width: 340px;
	margin: 20px;
}
h1 { color: #c4a7ff; }
.status { color: #64e89a; }
</style>
</head>
<body>
<div class="card">
<h1>🎵 KanhaMusic</h1>
<p>Telegram Music Bot</p>
<p class="status">● Web Server Online</p>
<p>Music · Voice Chats · Assistant</p>
</div>
</body>
</html>`)
		})

		logger.Infof("HTTP server listening on %s", addr)

		if err := http.ListenAndServe(addr, mux); err != nil {
			logger.Infof("HTTP server failed: %v", err)
		}
	}()
}


func refreshDirs() error {
	dirs := []string{
		"./cache",
		"./downloads",
	}

	for _, dir := range dirs {

		if err := os.RemoveAll(dir); err != nil {
			return err
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}
