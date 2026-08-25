package site

// A5 owns the stranger-facing pages. The first commit contains small fixed
// page examples for the route seam; replace those entries in the build-scoped
// set so this lane can ship the actual pages without editing another lane's
// registration files.
func init() {
	RegisterContent("a5-pages", func(data PageData, routes *RouteSet) {
		removeA5Placeholder(routes, "home")
		removeA5Placeholder(routes, "about")
		removeA5Placeholder(routes, "contact")

		routes.Register(Route{
			Name:   "a5-home",
			Output: RouteOutput{Path: "index.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5HomePage(buildData, renderHeadFragments(buildData))
			},
		})
		routes.Register(Route{
			Name:   "a5-about",
			Output: RouteOutput{Path: "about/index.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5AboutPage(buildData, renderHeadFragments(buildData))
			},
		})
		routes.Register(Route{
			Name:   "a5-contact",
			Output: RouteOutput{Path: "contact/index.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5ContactPage(buildData, renderHeadFragments(buildData))
			},
		})
		routes.Register(Route{
			Name:   "a5-stack",
			Output: RouteOutput{Path: "stack/index.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5StackPage(buildData, renderHeadFragments(buildData))
			},
		})
		routes.Register(Route{
			Name:   "a5-not-found",
			Output: RouteOutput{Path: "404.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5NotFoundPage(buildData, renderHeadFragments(buildData))
			},
		})
		routes.Register(Route{
			Name:   "a5-impressum",
			Output: RouteOutput{Path: "impressum/index.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5ImpressumPage(buildData, renderHeadFragments(buildData))
			},
		})
		routes.Register(Route{
			Name:   "a5-datenschutz",
			Output: RouteOutput{Path: "datenschutz/index.html"},
			Render: func(buildData PageData) interfaceComponent {
				return a5DatenschutzPage(buildData, renderHeadFragments(buildData))
			},
		})
	})
}

func removeA5Placeholder(routes *RouteSet, name string) {
	route, ok := routes.routes[name]
	if !ok {
		return
	}
	delete(routes.routes, name)
	if route.Output.Path != "" {
		delete(routes.staticOutput, route.Output.Path)
	}
}
