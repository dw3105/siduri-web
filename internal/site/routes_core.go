package site

func init() {
	Register(Route{
		Name:   "home",
		Output: "index.html",
		Render: func(data PageData) interfaceComponent {
			return HomePage(data.Posts, data.Preview)
		},
	})
	Register(Route{
		Name:   "journal",
		Output: "journal/index.html",
		Render: func(data PageData) interfaceComponent {
			return JournalPage(data.Posts, data.Preview)
		},
	})
	Register(Route{
		Name:   "contact",
		Output: "contact/index.html",
		Render: func(data PageData) interfaceComponent {
			return ContactPage(data.Preview)
		},
	})
}
