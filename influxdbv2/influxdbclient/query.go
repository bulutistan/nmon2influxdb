package influxdbclient

import "fmt"

func buildQuery(fields string, filters *Filters, groupby string, serie string, from string, to string, function string) (query string) {
	if len(function) > 0 {
		query = fmt.Sprintf("SELECT %s(\"%s\") FROM \"%s\"", function, fields, serie)
	} else {
		query = fmt.Sprintf("SELECT \"%s\" FROM \"%s\"", fields, serie)
	}

	var filterQuery FilterQuery

	if len(from) > 0 {
		filterQuery.Append(fmt.Sprintf("time > '%s'", from))
	}

	if len(to) > 0 {
		filterQuery.Append(fmt.Sprintf("time < '%s'", to))
	}

	if len(*filters) > 0 {
		filterQuery.AddFilters(filters)
	}

	if len(filterQuery.Content) > 0 {
		query += " WHERE " + filterQuery.Content
	}

	if len(groupby) > 0 {
		query += fmt.Sprintf(" GROUP BY \"%s\"", groupby)
	}
	return
}
