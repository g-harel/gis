// TODO rename to type_index
package source_index

import (
	"sort"
	"strings"

	"github.com/g-harel/gothrough/internal/string_index"
	"github.com/g-harel/gothrough/internal/types"
)

// Confidence values for info items.
const (
	confidenceHigh = 120
	confidenceMed  = 80
	confidenceLow  = 20
)

type Result struct {
	Confidence        float64
	Name              string
	PackageName       string
	PackageImportPath string
	Value             types.Type
}

type Index struct {
	textIndex         *string_index.Index
	results           []*Result
	computed_packages *[][]string
}

func NewIndex() *Index {
	return &Index{
		textIndex: string_index.NewIndex(),
		results:   []*Result{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (si *Index) Search(query string) ([]*Result, error) {
	matches := si.textIndex.Search(query)
	if len(matches) == 0 {
		return []*Result{}, nil
	}

	results := make([]*Result, len(matches))
	for i, match := range matches {
		result := si.results[match.ID]
		result.Confidence = match.Confidence
		results[i] = result
	}

	return results, nil
}

func (si *Index) Packages() [][]string {
	if si.computed_packages != nil {
		return *si.computed_packages
	}

	// Collect list of unique packages, separating the standard library vs. hosted ones.
	seenPackages := map[string]bool{}
	stdPackages := []string{}
	hostedPackages := map[string][]string{}

	// Add package names.
	for _, result := range si.results {
		packageName := result.PackageImportPath

		if seenPackages[packageName] {
			continue
		}
		seenPackages[packageName] = true

		firstNamePart := strings.Split(packageName, "/")[0]
		if !strings.Contains(firstNamePart, ".") {
			stdPackages = append(stdPackages, packageName)
			continue
		}
		if _, ok := hostedPackages[firstNamePart]; !ok {
			hostedPackages[firstNamePart] = []string{}
		}
		hostedPackages[firstNamePart] = append(hostedPackages[firstNamePart], packageName)
	}

	// Create sorted list of hosts.
	hosts := []string{}
	for host := range hostedPackages {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)

	// Created nested array of packages grouped by host and in sorted host order.
	// Standard library packages are added to the front.
	packages := [][]string{stdPackages}
	for _, host := range hosts {
		packages = append(packages, hostedPackages[host])
	}

	// Sort packages within each host's list.
	for i := range packages {
		sort.Strings(packages[i])
	}

	si.computed_packages = &packages
	return packages
}
