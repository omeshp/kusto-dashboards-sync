package models

type Query struct {
	DataSource struct {
		Kind         string `json:"kind"`
		DataSourceId string `json:"dataSourceId"`
	} `json:"dataSource"`
	Kind          string   `json:"kind,omitempty"`
	Text          string   `json:"text"`
	Id            string   `json:"id,omitempty"`
	UsedVariables []string `json:"usedVariables"`
}

type Base struct {
	Id                string        `json:"id,omitempty"`
	Label             string        `json:"label,omitempty"`
	Columns           []interface{} `json:"columns,omitempty"`
	YAxisMaximumValue interface{}   `json:"yAxisMaximumValue,omitempty"`
	YAxisMinimumValue interface{}   `json:"yAxisMinimumValue,omitempty"`
	YAxisScale        string        `json:"yAxisScale,omitempty"`
	HorizontalLines   []interface{} `json:"horizontalLines,omitempty"`
}

type MultipleYAxes struct {
	Base               *Base         `json:"base,omitempty"`
	Additional         []interface{} `json:"additional,omitempty"`
	ShowMultiplePanels bool          `json:"showMultiplePanels,omitempty"`
}

type VisualOptions struct {
	MultipleYAxes          *MultipleYAxes `json:"multipleYAxes,omitempty"`
	HideLegend             bool           `json:"hideLegend,omitempty"`
	XColumnTitle           string         `json:"xColumnTitle,omitempty"`
	XColumn                interface{}    `json:"xColumn,omitempty"`
	YColumns               interface{}    `json:"yColumns,omitempty"`
	SeriesColumns          interface{}    `json:"seriesColumns,omitempty"`
	XAxisScale             string         `json:"xAxisScale,omitempty"`
	VerticalLine           string         `json:"verticalLine,omitempty"`
	CrossFilterDisabled    bool           `json:"crossFilterDisabled,omitempty"`
	DrillthroughDisabled   bool           `json:"drillthroughDisabled,omitempty"`
	CrossFilter            []interface{}  `json:"crossFilter,omitempty"`
	Drillthrough           []interface{}  `json:"drillthrough,omitempty"`
	TableEnableRenderLinks bool           `json:"table__enableRenderLinks,omitempty"`
	ColorRules             []struct {
		Id            string      `json:"id"`
		RuleType      string      `json:"ruleType"`
		ApplyToColumn interface{} `json:"applyToColumn"`
		HideText      bool        `json:"hideText"`
		ApplyTo       string      `json:"applyTo"`
		Conditions    []struct {
			Operator string   `json:"operator"`
			Column   string   `json:"column"`
			Values   []string `json:"values"`
		} `json:"conditions"`
		ChainingOperator string  `json:"chainingOperator"`
		VisualType       string  `json:"visualType"`
		ColorStyle       string  `json:"colorStyle"`
		Color            *string `json:"color"`
		Tag              string  `json:"tag"`
		Icon             *string `json:"icon"`
		RuleName         string  `json:"ruleName"`
	} `json:"colorRules,omitempty"`
	ColorRulesDisabled bool   `json:"colorRulesDisabled,omitempty"`
	ColorStyle         string `json:"colorStyle,omitempty"`
	TableRenderLinks   []struct {
		UrlColumn string `json:"urlColumn"`
		Disabled  bool   `json:"disabled"`
	} `json:"table__renderLinks,omitempty"`
	LabelDisabled               bool        `json:"labelDisabled,omitempty"`
	PieLabel                    []string    `json:"pie__label,omitempty"`
	TooltipDisabled             bool        `json:"tooltipDisabled,omitempty"`
	PieTooltip                  []string    `json:"pie__tooltip,omitempty"`
	PieOrderBy                  string      `json:"pie__orderBy,omitempty"`
	PieKind                     string      `json:"pie__kind,omitempty"`
	PieTopNSlices               interface{} `json:"pie__topNSlices"`
	MultiStatTextSize           string      `json:"multiStat__textSize,omitempty"`
	MultiStatValueColumn        *string     `json:"multiStat__valueColumn,omitempty"`
	MultiStatDisplayOrientation string      `json:"multiStat__displayOrientation,omitempty"`
	MultiStatLabelColumn        string      `json:"multiStat__labelColumn,omitempty"`
	MultiStatSlot               struct {
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
	} `json:"multiStat__slot,omitempty"`
}

type Dashboard struct {
	Schema            string `json:"$schema"`
	Id                string `json:"id"`
	IsDashboardEditor bool   `json:"isDashboardEditor"`
	ETag              string `json:"eTag"`
	SchemaVersion     string `json:"schema_version"`
	Title             string `json:"title"`
	Tiles             []struct {
		Title      string `json:"title"`
		PageId     string `json:"pageId,omitempty"`
		Id         string `json:"id"`
		VisualType string `json:"visualType"`
		Layout     struct {
			X      int `json:"x"`
			Y      int `json:"y"`
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"layout"`
		QueryRef struct {
			Kind    string `json:"kind,omitempty"`
			QueryId string `json:"queryId,omitempty"`
		} `json:"queryRef,omitempty"`
		Query         Query         `json:"query,omitempty"`
		VisualOptions VisualOptions `json:"visualOptions"`
		MarkdownText  string        `json:"markdownText,omitempty"`
	} `json:"tiles"`
	BaseQueries []interface{} `json:"baseQueries"`
	Parameters  []struct {
		Kind              string `json:"kind"`
		Id                string `json:"id"`
		DisplayName       string `json:"displayName"`
		Description       string `json:"description"`
		BeginVariableName string `json:"beginVariableName"`
		EndVariableName   string `json:"endVariableName"`
		VariableName      string `json:"variableName,omitempty"`
		SelectionType     string `json:"selectionType,omitempty"`
		IncludeAllOption  bool   `json:"includeAllOption,omitempty"`
		DefaultValue      struct {
			Kind   string   `json:"kind"`
			Count  int      `json:"count,omitempty"`
			Unit   string   `json:"unit,omitempty"`
			Value  string   `json:"value,omitempty"`
			Values []string `json:"values,omitempty"`
		} `json:"defaultValue"`
		DataSource struct {
			Kind   string `json:"kind,omitempty"`
			Values []struct {
				DisplayText string `json:"displayText,omitempty"`
				Value       string `json:"value,omitempty"`
			} `json:"values,omitempty"`
			Columns struct {
				Value string `json:"value,omitempty"`
			} `json:"columns,omitempty"`
			QueryRef struct {
				Kind    string `json:"kind,omitempty"`
				QueryId string `json:"queryId,omitempty"`
			} `json:"queryRef,omitempty"`
		} `json:"dataSource,omitempty"`
		ShowOnPages struct {
			Kind    string   `json:"kind"`
			PageIds []string `json:"pageIds,omitempty"`
		} `json:"showOnPages"`
	} `json:"parameters"`
	DataSources []struct {
		Id         string `json:"id"`
		Kind       string `json:"kind"`
		ScopeId    string `json:"scopeId"`
		Name       string `json:"name"`
		ClusterUri string `json:"clusterUri"`
		Database   string `json:"database"`
	} `json:"dataSources"`
	Pages []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"pages"`
	Queries []Query `json:"queries"`
}
