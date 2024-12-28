package handler

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type KnowledgeBaseCtrl struct {
	BaseCtrl
	KnowledgeBaseService *service.KnowledgeBaseService `inject:""`
	FileService          *service.FileService          `inject:""`
}

var (
	kbStr = "b0b12d74-2f56-49a8-9fad-8f5c6919b85e"
)

func (c *KnowledgeBaseCtrl) UploadDoc(ctx iris.Context) {
	kb := ctx.URLParamDefault("kb", kbStr)

	f, fh, err := ctx.FormFile("file")
	if err != nil {
		_logs.Errorf("文件上传失败", zap.String("ctx.FormFile(\"file\")", err.Error()))
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	defer f.Close()

	name := fh.Filename
	pth, err := c.FileService.UploadFile(ctx, fh, kb)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	err = c.KnowledgeBaseService.UploadZipFile(pth, kb)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: iris.Map{"path": pth, "name": name}})
}

func (c *KnowledgeBaseCtrl) ClearAll(ctx iris.Context) {
	kb := ctx.URLParamDefault("kb", kbStr)

	err := c.KnowledgeBaseService.ClearAll(kb)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code})
}