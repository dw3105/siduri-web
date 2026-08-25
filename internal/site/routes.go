package site

import (
	"fmt"
	"sort"
	"sync"

	"github.com/a-h/templ"
)

// PageData is the stable input passed to registered pages. New lanes can add a
// route in their own file and call Register from init without changing main.
type PageData struct {
	Posts   []Post
	Preview bool
}

type Route struct {
	Name   string
	Output string
	Render func(PageData) templ.Component
}

var (
	routeMu  sync.RWMutex
	routeSet = make(map[string]Route)
)

// Register adds a build route. Duplicate names or output paths are programmer
// errors and fail during package initialization rather than at publish time.
func Register(route Route) {
	if route.Name == "" || route.Output == "" || route.Render == nil {
		panic("site: incomplete route registration")
	}
	routeMu.Lock()
	defer routeMu.Unlock()
	if _, exists := routeSet[route.Name]; exists {
		panic(fmt.Sprintf("site: duplicate route %q", route.Name))
	}
	for _, existing := range routeSet {
		if existing.Output == route.Output {
			panic(fmt.Sprintf("site: duplicate route output %q", route.Output))
		}
	}
	routeSet[route.Name] = route
}

func registeredRoutes() []Route {
	routeMu.RLock()
	defer routeMu.RUnlock()
	routes := make([]Route, 0, len(routeSet))
	for _, route := range routeSet {
		routes = append(routes, route)
	}
	sort.Slice(routes, func(i, j int) bool { return routes[i].Output < routes[j].Output })
	return routes
}
