package domain

import (
	"sort"
	"strings"
)

func CopyRoute(route NightTourRoute) NightTourRoute {
	copy := route
	copy.Stops = append([]string(nil), route.Stops...)
	return copy
}

func SortRoutes(routes []NightTourRoute) []NightTourRoute {
	result := append([]NightTourRoute(nil), routes...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].District == result[j].District {
			return result[i].Name < result[j].Name
		}
		return result[i].District < result[j].District

	})
	return result
}

func RouteMatches(route NightTourRoute, query RouteSearch) bool {
	if query.ActiveOnly && !route.Active {
		return false
	}
	if query.District != "" && NormalizeTerm(route.District) != NormalizeTerm(query.District) {
		return false
	}
	term := NormalizeTerm(query.Term)
	if term == "" {
		return true
	}
	return strings.Contains(strings.ToLower(route.Name), term) || strings.Contains(strings.ToLower(route.Summary), term) || strings.Contains(strings.ToLower(route.District), term)
}
