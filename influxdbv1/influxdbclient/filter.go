package influxdbclient

import "fmt"

type Filter struct {
	Tag   string
	Value string
	Mode  string
}

type Filters []Filter

func (filters *Filters) Add(tag string, value string, mode string) {
	*filters = append(*filters, Filter{Tag: tag, Value: value, Mode: mode})
}

type FilterQuery struct {
	Content string
}

func (fQuery *FilterQuery) Append(text string) {
	if len(fQuery.Content) > 0 {
		fQuery.Content += " AND " + text
	} else {
		fQuery.Content = text
	}
}

func (fQuery *FilterQuery) AddFilters(filters *Filters) {
	for _, filter := range *filters {
		switch {
		case filter.Mode == "text":
			fQuery.Append(fmt.Sprintf("\"%s\" = '%s'", filter.Tag, filter.Value))
		case filter.Mode == "regexp":
			fQuery.Append(fmt.Sprintf("\"%s\" =~ /%s/", filter.Tag, filter.Value))
		}
	}
}
