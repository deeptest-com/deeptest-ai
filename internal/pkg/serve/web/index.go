package web

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/joho/godotenv"
	"os"
)

// init
func init() {
	viper_server.Init(config.GetViperConfig())

	err := godotenv.Load()
	if err == nil {
		config.CONFIG.Ai.ApiKey = os.Getenv("AI_PLATFORM_API_KEY")
	}
}

type WebBaseFunc interface {
	AddWebStatic(staticAbsPath, webPrefix string, paths ...string)
	AddUploadStatic(staticAbsPath, webPrefix string)
	InitRouter() error
	Run()
	GetSources() ([]map[string]string, []map[string]string)
}

type WebFunc interface {
	WebBaseFunc
}

var PermRoutes = make([]map[string]string, 0)

// Start
func Start(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}

	PermRoutes, _ = wf.GetSources()
	wf.Run()
}

func StartTest(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		_logUtils.Error(err.Error())
	}
}
