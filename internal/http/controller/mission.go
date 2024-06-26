// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"
	"strconv"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/http/middleware"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/times"
)

type MissionController struct{}

// 注册路由
func (self MissionController) RegisterRoute(g *echo.Group) {
	g.GET("/mission/daily", self.Daily, middleware.NeedLogin())
	g.GET("/mission/daily/redeem", self.DailyRedeem, middleware.NeedLogin())
	g.GET("/mission/complete/:id", self.Complete, middleware.NeedLogin())
}

func (MissionController) Daily(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	userLoginMission := logic.DefaultMission.FindLoginMission(context.EchoContext(ctx), me)
	userLoginMission.Uid = me.Uid

	data := map[string]interface{}{"login_mission": userLoginMission}

	if userLoginMission != nil && times.Format("Ymd") == strconv.Itoa(userLoginMission.Date) {
		data["had_redeem"] = true
	} else {
		data["had_redeem"] = false
	}

	fr := ctx.QueryParam("fr")
	if fr == "redeem" {
		data["show_msg"] = true
	}
	return render(ctx, "mission/daily.html", data)
}

func (MissionController) DailyRedeem(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	logic.DefaultMission.RedeemLoginAward(context.EchoContext(ctx), me)

	return ctx.Redirect(http.StatusSeeOther, "/mission/daily?fr=redeem")
}

func (MissionController) Complete(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	id := ctx.Param("id")
	logic.DefaultMission.Complete(context.EchoContext(ctx), me, id)

	return ctx.Redirect(http.StatusSeeOther, "/balance")
}
