// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

// 喜欢系统

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/http/middleware"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type LikeController struct{}

// 注册路由
func (self LikeController) RegisterRoute(g *echo.Group) {
	g.POST("/like/:objid", self.Like, middleware.NeedLogin())
}

// Like 喜欢（或取消喜欢）
func (LikeController) Like(ctx echo.Context) error {
	form, _ := ctx.FormParams()
	if !util.CheckInt(form, "objtype") || !util.CheckInt(form, "flag") {
		return fail(ctx, 1, "参数错误")
	}

	user := ctx.Get("user").(*model.Me)
	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.FormValue("objtype"))
	likeFlag := goutils.MustInt(ctx.FormValue("flag"))

	err := logic.DefaultLike.LikeObject(context.EchoContext(ctx), user.Uid, objid, objtype, likeFlag)
	if err != nil {
		return fail(ctx, 2, "服务器内部错误")
	}

	return success(ctx, nil)
}
