package deps

import "fmt"

// Graph holds a dependency graph for services.
type Graph struct {
	edges map[string][]string
}

// New creates a new dependency Graph from a map of service -> dependencies.
func New(deps map[string][]string) *Graph {
	return &Graph{edges: deps}
}

// Order returns a topologically sorted list of services so that each service
// appears after all of its dependencies. Returns an error if a cycle is found.
func (g *Graph) Order(services []string) ([]string, error) {
	visited := make(map[string]int) // 0=unvisited,1=visiting,2=done
	var result []string

	var visit func(name string) error
	visit = func(name string) error {
		switch visited[name] {
		case 2:
			return nil
		case 1:
			return fmt.Errorf("cycle detected at service %q", name)
		}
		visited[name] = 1
		for _, dep := range g.edges[name] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visited[name] = 2
		result = append(result, name)
		return nil
	}

	for _, svc := range services {
		if err := visit(svc); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Deps returns the direct dependencies of a service.
func (g *Graph) Deps(service string) []string {
	return g.edges[service]
}
