package service

import (
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	httpUtils "github.com/deeptest-com/deeptest-next/internal/pkg/libs/http"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	"github.com/deeptest-com/deeptest-next/pkg/libs/http"
	"github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"os"
	"path/filepath"
)

type KnowledgeBaseService struct {
}

var (
	defaultDb = "b0b12d74-2f56-49a8-9fad-8f5c6919b85e"

	kbCreateDocUri = "datasets/%s/document/create-by-file"
	kbQueryDocUri  = "datasets/%s/documents"
	kbRemoveDocUri = "datasets/%s/documents/%s"
)

func (s *KnowledgeBaseService) UploadZipFile(zipPath, kb string) (err error) {
	unzipDir, err := _file.Unzip(zipPath, filepath.Join(consts.WorkDir, "_temp"))
	if err != nil {
		return
	}
	filePaths := s.ListFile(unzipDir)

	for _, filePath := range filePaths {
		err := s.UploadDoc(filePath, kb)
		if err != nil {
			continue
		}
	}

	return
}

func (s *KnowledgeBaseService) UploadDoc(pth, kb string) (err error) {
	if kb == "" {
		kb = defaultDb
	}

	url := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		url = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbCreateDocUri, kb)
	}
	_logs.Infof("%s url = %s", config.CONFIG.Ai.PlatformType, url)

	data := domain.KbCreateReq{
		IndexingTechnique: "high_quality",
	}

	bts, err := httpUtils.PostFile(url, data, pth, s.getHeaders())
	if err != nil {
		return
	}
	_logs.Infof("create doc resp %s", string(bts))

	return
}

func (s *KnowledgeBaseService) ClearAll(kb string) (err error) {
	if kb == "" {
		kb = defaultDb
	}

	queryUrl := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		queryUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbQueryDocUri, kb)
	}
	_logs.Infof("%s queryUrl = %s", config.CONFIG.Ai.PlatformType, queryUrl)

	headers := s.getHeaders()
	bts, err := httpUtils.Get(queryUrl, headers)
	if err != nil {
		return
	}

	docs := domain.KbQueryResult{}
	json.Unmarshal(bts, &docs)

	for _, doc := range docs.Data {
		removeUrl := ""
		if config.CONFIG.Ai.PlatformType == consts.Dify {
			removeUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
				fmt.Sprintf(kbRemoveDocUri, kb, doc.Id)
		}
		_logs.Infof("%s removeUrl = %s", config.CONFIG.Ai.PlatformType, removeUrl)

		bts, err = httpUtils.Delete(removeUrl, headers)
		if err != nil {
			continue
		}
	}

	return
}

func (s *KnowledgeBaseService) getHeaders() (ret map[string]string) {
	ret = map[string]string{"Authorization": "Bearer " + os.Getenv("AI_DATASET_API_KEY")}
	return
}

func (s *KnowledgeBaseService) ListFile(pth string) (ret []string) {
	return
}