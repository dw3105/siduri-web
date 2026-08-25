package site

import (
	"fmt"
	"sort"
	"sync"

	"github.com/a-h/templ"
)

// PageData is the site-wide input available to every route lane. A lane keeps
// any derived data in its own registration closure rather than adding fields
// here.
type PageData struct {
	Posts   []Post
	Tags    []string
	Preview bool
}

// Route describes one registration. A path is the one-page form. Expand is
// the fan-out form: it computes the outputs after content has been loaded.
// Render is the convenient component renderer for the one-page form.
type Route struct {
	Name   string
	Output RouteOutput
	Render func(PageData) templ.Component
}

type RouteOutput struct {
	Path   string
	Expand func(PageData) []Output
}

// Output is one file produced by a route. Exactly one of Render or Bytes must
// be set. Bytes are written as-is; components receive the normal HTML newline.
type Output struct {
	Path   string
	Render func() templ.Component
	Bytes  []byte
}

func PageOutput(path string, render func() templ.Component) Output {
	return Output{Path: path, Render: render}
}

func ByteOutput(path string, data []byte) Output {
	return Output{Path: path, Bytes: data}
}

// ContentRegistration is called once per build, after content and the tag
// vocabulary have been loaded. It registers routes into that build's local
// set, so a second build in the same process starts clean.
type ContentRegistration func(PageData, *RouteSet)

type contentRegistration struct {
	Name     string
	Register ContentRegistration
}

// RouteSet is deliberately build-scoped. Its Register method is the seam for
// lanes whose routes depend on loaded content.
type RouteSet struct {
	routes       map[string]Route
	staticOutput map[string]string
}

func newRouteSet() *RouteSet {
	return &RouteSet{
		routes:       make(map[string]Route),
		staticOutput: make(map[string]string),
	}
}

func (set *RouteSet) Register(route Route) {
	validateRoute(route)
	if _, exists := set.routes[route.Name]; exists {
		panic(fmt.Sprintf("site: duplicate route %q", route.Name))
	}
	if route.Output.Path != "" {
		if existing, exists := set.staticOutput[route.Output.Path]; exists {
			panic(fmt.Sprintf("site: duplicate route output %q (routes %q and %q)", route.Output.Path, existing, route.Name))
		}
		set.staticOutput[route.Output.Path] = route.Name
	}
	set.routes[route.Name] = route
}

func (set *RouteSet) sorted() []Route {
	routes := make([]Route, 0, len(set.routes))
	for _, route := range set.routes {
		routes = append(routes, route)
	}
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Output.Path != routes[j].Output.Path {
			return routes[i].Output.Path < routes[j].Output.Path
		}
		return routes[i].Name < routes[j].Name
	})
	return routes
}

var (
	routeMu           sync.RWMutex
	registeredRoutes  = make(map[string]Route)
	contentRegistries = make(map[string]ContentRegistration)
	headFragments     = make(map[string]HeadFragment)
)

// Register adds a build-independent route declaration. Duplicate route names
// and duplicate static output paths are programmer errors and panic here.
func Register(route Route) {
	routeMu.Lock()
	defer routeMu.Unlock()
	set := &RouteSet{routes: registeredRoutes, staticOutput: make(map[string]string)}
	for name, existing := range registeredRoutes {
		if existing.Output.Path != "" {
			set.staticOutput[existing.Output.Path] = name
		}
	}
	set.Register(route)
}

// RegisterContent lets a lane register routes after Build has loaded its
// typed site-wide input. The callback itself is declared from the lane's init.
func RegisterContent(name string, register ContentRegistration) {
	if name == "" || register == nil {
		panic("site: incomplete content route registration")
	}
	routeMu.Lock()
	defer routeMu.Unlock()
	if _, exists := contentRegistries[name]; exists {
		panic(fmt.Sprintf("site: duplicate content route registration %q", name))
	}
	contentRegistries[name] = register
}

func routesForBuild(data PageData) *RouteSet {
	set := newRouteSet()

	routeMu.RLock()
	static := make([]Route, 0, len(registeredRoutes))
	for _, route := range registeredRoutes {
		static = append(static, route)
	}
	content := make([]contentRegistration, 0, len(contentRegistries))
	for name, register := range contentRegistries {
		content = append(content, contentRegistration{Name: name, Register: register})
	}
	routeMu.RUnlock()

	sort.Slice(static, func(i, j int) bool { return static[i].Name < static[j].Name })
	for _, route := range static {
		set.Register(route)
	}
	sort.Slice(content, func(i, j int) bool { return content[i].Name < content[j].Name })
	for _, registration := range content {
		registration.Register(data, set)
	}
	return set
}

func outputsForBuild(set *RouteSet, data PageData) []routeOutput {
	used := make(map[string]string)
	var outputs []routeOutput
	for _, route := range set.sorted() {
		var routeOutputs []Output
		if route.Output.Path != "" {
			routeOutputs = []Output{PageOutput(route.Output.Path, func() templ.Component {
				return route.Render(data)
			})}
		} else {
			routeOutputs = route.Output.Expand(data)
		}
		for _, output := range routeOutputs {
			validateOutput(output)
			if existing, exists := used[output.Path]; exists {
				panic(fmt.Sprintf("site: duplicate route output %q (routes %q and %q)", output.Path, existing, route.Name))
			}
			used[output.Path] = route.Name
			outputs = append(outputs, routeOutput{Route: route.Name, Output: output})
		}
	}
	return outputs
}

type routeOutput struct {
	Route  string
	Output Output
}

func validateRoute(route Route) {
	if route.Name == "" {
		panic("site: incomplete route registration")
	}
	if (route.Output.Path == "") == (route.Output.Expand == nil) {
		panic(fmt.Sprintf("site: incomplete route output for %q", route.Name))
	}
	if route.Output.Path != "" && route.Render == nil {
		panic(fmt.Sprintf("site: missing component renderer for %q", route.Name))
	}
	if route.Output.Path == "" && route.Render != nil {
		panic(fmt.Sprintf("site: fan-out route %q must render each output", route.Name))
	}
}

func validateOutput(output Output) {
	if output.Path == "" || (output.Render == nil && output.Bytes == nil) || (output.Render != nil && output.Bytes != nil) {
		panic("site: incomplete route output")
	}
}

type HeadFragment struct {
	Name   string
	Render func(PageData) templ.Component
}

// RegisterHead contributes one typed component to every document. Fragments
// are sorted by name before rendering, making init order irrelevant.
func RegisterHead(fragment HeadFragment) {
	if fragment.Name == "" || fragment.Render == nil {
		panic("site: incomplete head fragment registration")
	}
	routeMu.Lock()
	defer routeMu.Unlock()
	if _, exists := headFragments[fragment.Name]; exists {
		panic(fmt.Sprintf("site: duplicate head fragment %q", fragment.Name))
	}
	headFragments[fragment.Name] = fragment
}

func renderHeadFragments(data PageData) []templ.Component {
	routeMu.RLock()
	fragments := make([]HeadFragment, 0, len(headFragments))
	for _, fragment := range headFragments {
		fragments = append(fragments, fragment)
	}
	routeMu.RUnlock()
	sort.Slice(fragments, func(i, j int) bool { return fragments[i].Name < fragments[j].Name })
	components := make([]templ.Component, 0, len(fragments))
	for _, fragment := range fragments {
		components = append(components, fragment.Render(data))
	}
	return components
}

type interfaceComponent = templ.Component
