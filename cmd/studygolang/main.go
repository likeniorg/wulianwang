// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/studygolang/studygolang/cmd"
	"github.com/studygolang/studygolang/global"
	"github.com/studygolang/studygolang/internal/http/controller"
	"github.com/studygolang/studygolang/internal/http/controller/admin"
	"github.com/studygolang/studygolang/internal/http/controller/app"
	pwm "github.com/studygolang/studygolang/internal/http/middleware"
	"github.com/studygolang/studygolang/internal/logic"
	thirdmw "github.com/studygolang/studygolang/middleware"

	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	. "github.com/polaris1119/config"
	"github.com/polaris1119/keyword"
	"github.com/polaris1119/logger"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())

	structs.DefaultTagName = "json"
}

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "indexer":
			cmd.Indexer()
			return
		case "crawler":
			cmd.Crawler()
			return
		}
	}

	// 支持根据参数打印版本信息
	global.PrintVersion(os.Stdout)

	savePid()

	global.App.Init(logic.WebsiteSetting.Domain)

	logger.Init(ROOT+"/log", ConfigFile.MustValue("global", "log_level", "DEBUG"))

	go keyword.Extractor.Init(keyword.DefaultProps, true, ROOT+"/data/programming.txt,"+ROOT+"/data/dictionary.txt")

	go logic.Book.ClearRedisUser()

	go ServeBackGround()
	// go pprof
	Pprof(ConfigFile.MustValue("global", "pprof", "127.0.0.1:8096"))

	e := echo.New()

	serveStatic(e)

	e.Use(thirdmw.EchoLogger())
	e.Use(mw.Recover())
	e.Use(pwm.Installed(filterPrefixs))
	e.Use(pwm.HTTPError())
	e.Use(pwm.AutoLogin())

	// 评论后不会立马显示出来，暂时缓存去掉
	// frontG := e.Group("", thirdmw.EchoCache())
	frontG := e.Group("")
	controller.RegisterRoutes(frontG)

	adminG := e.Group("/admin", pwm.NeedLogin(), pwm.AdminAuth())
	admin.RegisterRoutes(adminG)

	// appG := e.Group("/app", thirdmw.EchoCache())
	appG := e.Group("/app")
	app.RegisterRoutes(appG)

	e.Server.Addr = getAddr()
	gracefulRun(e.Server)
}

func getAddr() string {
	host := ConfigFile.MustValue("listen", "host", "")
	if host == "" {
		global.App.Host = "localhost"
	} else {
		global.App.Host = host
	}
	global.App.Port = ConfigFile.MustValue("listen", "port", "8088")
	return host + ":" + global.App.Port
}

func savePid() {
	pidFilename := ROOT + "/pid/" + filepath.Base(os.Args[0]) + ".pid"
	pid := os.Getpid()

	ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0755)
}
