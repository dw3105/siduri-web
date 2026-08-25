package site

func init() {
	Register(Route{
		Name:   "home",
		Output: RouteOutput{Path: "index.html"},
		Render: func(data PageData) interfaceComponent {
			return homePage(data, renderHeadFragments(data))
		},
	})
	Register(Route{
		Name:   "journal",
		Output: RouteOutput{Path: "journal/index.html"},
		Render: func(data PageData) interfaceComponent {
			return journalPage(data, renderHeadFragments(data))
		},
	})
	Register(Route{
		Name:   "contact",
		Output: RouteOutput{Path: "contact/index.html"},
		Render: func(data PageData) interfaceComponent {
			return contactPage(data, renderHeadFragments(data))
		},
	})
}

// These wrappers keep the small rendering API used by golden tests and by
// callers that render a page outside a Build.
func HomePage(posts []Post, preview bool) interfaceComponent {
	return homePage(PageData{Posts: posts, Preview: preview}, nil)
}

func JournalPage(posts []Post, preview bool) interfaceComponent {
	return journalPage(PageData{Posts: posts, Preview: preview}, nil)
}

func ContactPage(preview bool) interfaceComponent {
	return contactPage(PageData{Preview: preview}, nil)
}
