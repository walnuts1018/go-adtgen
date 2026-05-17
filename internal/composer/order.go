package composer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/walnuts1018/go-product-type/internal/model"
)

func OrderDeclarations(declarations []model.ResolvedDeclaration) ([]model.ResolvedDeclaration, error) {
	if len(declarations) == 0 {
		return nil, nil
	}

	localNames := make(map[string]struct{}, len(declarations))
	byName := make(map[string]model.ResolvedDeclaration, len(declarations))
	indegree := make(map[string]int, len(declarations))
	dependents := make(map[string][]string, len(declarations))

	for _, declaration := range declarations {
		name := declaration.Declaration.Name
		if _, ok := localNames[name]; ok {
			return nil, fmt.Errorf("composer: duplicate declaration name: %s", name)
		}
		localNames[name] = struct{}{}
		byName[name] = declaration
		indegree[name] = 0
	}

	for _, declaration := range declarations {
		name := declaration.Declaration.Name
		dependencies, err := declarationDependencies(declaration, localNames)
		if err != nil {
			return nil, fmt.Errorf("composer: %s: %w", declaration.Declaration.Name, err)
		}

		indegree[name] = len(dependencies)
		for dependency := range dependencies {
			dependents[dependency] = append(dependents[dependency], name)
		}
	}

	ready := make([]string, 0, len(declarations))
	for name, degree := range indegree {
		if degree == 0 {
			ready = append(ready, name)
		}
	}
	sort.Strings(ready)

	ordered := make([]model.ResolvedDeclaration, 0, len(declarations))
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		ordered = append(ordered, byName[name])

		next := dependents[name]
		sort.Strings(next)
		for _, dependent := range next {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
				sort.Strings(ready)
			}
		}
	}

	if len(ordered) != len(declarations) {
		return nil, fmt.Errorf("composer: declaration dependency cycle detected: %s", cycleNames(indegree))
	}

	return ordered, nil
}

func declarationDependencies(declaration model.ResolvedDeclaration, localNames map[string]struct{}) (map[string]struct{}, error) {
	dependencies := make(map[string]struct{})
	excludedNames := declarationTypeParameterNames(declaration.Declaration.TypeParameters)
	for _, input := range declaration.Inputs {
		names, err := input.ReferencedDeclarationNamesExcluding(localNames, excludedNames)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", input.Expr, err)
		}
		for _, name := range names {
			if name == declaration.Declaration.Name {
				continue
			}
			dependencies[name] = struct{}{}
		}
	}
	return dependencies, nil
}

func declarationTypeParameterNames(typeParameters []string) map[string]struct{} {
	if len(typeParameters) == 0 {
		return nil
	}

	names := make(map[string]struct{}, len(typeParameters))
	for _, typeParameter := range typeParameters {
		fields := strings.Fields(typeParameter)
		if len(fields) == 0 {
			continue
		}
		names[fields[0]] = struct{}{}
	}
	return names
}

func cycleNames(indegree map[string]int) string {
	names := make([]string, 0, len(indegree))
	for name, degree := range indegree {
		if degree > 0 {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
