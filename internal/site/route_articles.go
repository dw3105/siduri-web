package site

// Article pages are a content-registered fan-out: the route exists once, and
// its output set is computed from the posts available to this build.
func init() {
	RegisterContent("articles", func(data PageData, routes *RouteSet) {
		routes.Register(Route{
			Name: "articles",
			Output: RouteOutput{Expand: func(buildData PageData) []Output {
				outputs := make([]Output, 0, len(buildData.Posts))
				for _, post := range buildData.Posts {
					post := post
					body := renderArticleBody(post, buildData)
					sections := articleSectionsFor(post, buildData)
					outputs = append(outputs, PageOutput(
						"journal/"+post.Slug+"/index.html",
						func() interfaceComponent {
							return articlePage(post, body, buildData, renderHeadFragments(buildData), sections)
						},
					))
				}
				return outputs
			}},
		})
	})
}

func ArticlePage(post Post, body string, preview bool) interfaceComponent {
	data := PageData{Preview: preview}
	return articlePage(post, body, data, nil, articleSectionsFor(post, data))
}
