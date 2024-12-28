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
	"strings"
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
	filePaths, _ := s.ListFile(unzipDir, "img")

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

	data := s.getData()

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

func (s *KnowledgeBaseService) ListFile(dirName, exclude string) (ret []string, err error) {
	dirName = strings.TrimSuffix(dirName, string(os.PathSeparator))

	infos, err := os.ReadDir(dirName)
	if err != nil {
		return
	}

	for _, info := range infos {
		if info.Name() == exclude {
			continue
		}

		path := dirName + string(os.PathSeparator) + info.Name()
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if !realInfo.Mode().IsRegular() {
			continue
		}

		if realInfo.IsDir() {
			children, err := s.ListFile(path, exclude)
			if err != nil {
				return nil, err
			}
			ret = append(ret, children...)
			continue
		}

		ret = append(ret, path)

	}
	return
}

func (s *KnowledgeBaseService) getData() (ret domain.KbCreateReq) {
	rules := domain.Rules{Segmentation: domain.Segmentation{
		Separator: "###", MaxTokens: 500,
	}}
	rules.PreProcessingRules = append(rules.PreProcessingRules, domain.PreProcessingRule{
		Id:      "remove_extra_spaces",
		Enabled: true,
	})
	rules.PreProcessingRules = append(rules.PreProcessingRules, domain.PreProcessingRule{
		Id:      "remove_urls_emails",
		Enabled: true,
	})

	ret = domain.KbCreateReq{
		IndexingTechnique: "high_quality",
		ProcessRule: domain.ProcessRule{
			Rules: rules,
			Mode:  "custom",
		},
	}

	return
}
