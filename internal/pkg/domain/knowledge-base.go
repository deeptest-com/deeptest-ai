package domain

type KbCreateReq struct {
	IndexingTechnique string `json:"indexing_technique"`
	ProcessRule       struct {
		Rules struct {
			PreProcessingRules []struct {
				Id      string `json:"id"`
				Enabled bool   `json:"enabled"`
			} `json:"pre_processing_rules"`
			Segmentation struct {
				Separator string `json:"separator"`
				MaxTokens int    `json:"max_tokens"`
			} `json:"segmentation"`
		} `json:"rules"`
		Mode string `json:"mode"`
	} `json:"process_rule"`
}

type KbQueryResult struct {
	Data []struct {
		Id                   string      `json:"id"`
		Position             int         `json:"position"`
		DataSourceType       string      `json:"data_source_type"`
		DataSourceInfo       interface{} `json:"data_source_info"`
		DatasetProcessRuleId interface{} `json:"dataset_process_rule_id"`
		Name                 string      `json:"name"`
		CreatedFrom          string      `json:"created_from"`
		CreatedBy            string      `json:"created_by"`
		CreatedAt            int         `json:"created_at"`
		Tokens               int         `json:"tokens"`
		IndexingStatus       string      `json:"indexing_status"`
		Error                interface{} `json:"error"`
		Enabled              bool        `json:"enabled"`
		DisabledAt           interface{} `json:"disabled_at"`
		DisabledBy           interface{} `json:"disabled_by"`
		Archived             bool        `json:"archived"`
	} `json:"data"`
	HasMore bool `json:"has_more"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
}