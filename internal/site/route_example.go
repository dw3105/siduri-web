package site

// This is the worked registration example for later lanes. A lane-owned file
// can add another init/Register pair without editing main.go or the Makefile.
func init() {
	Register(Route{
		Name:   "about",
		Output: "about/index.html",
		Render: func(data PageData) interfaceComponent {
			return AboutPage(data.Preview)
		},
	})
}
