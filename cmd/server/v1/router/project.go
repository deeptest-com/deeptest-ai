package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type ProjectModule struct {
	ProjectCtrl *handler.ProjectCtrl `inject:""`
}

// Party 项目
func (m *ProjectModule) Party() func(index iris.Party) {
	return func(party iris.Party) {
		party.Use(middleware.MultiHandler(), middleware.Casbin())

		party.Get("/", m.ProjectCtrl.List).Name = "项目列表"
		party.Get("/{id:uint}", m.ProjectCtrl.Get).Name = "项目详情"
		party.Post("/", m.ProjectCtrl.Create).Name = "新建项目"
		party.Put("/", m.ProjectCtrl.Update).Name = "更新项目"
		party.Delete("/{id:uint}", m.ProjectCtrl.Delete).Name = "删除项目"

		party.Post("/changeProject", m.ProjectCtrl.ChangeProject).Name = "切换用户默认项目"
		party.Get("/getByUser", m.ProjectCtrl.GetByUser).Name = "获取用户参与的项目"

		party.Get("/members", m.ProjectCtrl.Members).Name = "获取项目成员"
		party.Post("/removeMember", m.ProjectCtrl.RemoveMember).Name = "删除项目成员"
		party.Post("/changeUserRole", m.ProjectCtrl.ChangeUserRole).Name = "更新项目成员的角色"
	}
}
