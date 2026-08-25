package site

import "sort"

func init() {
	RegisterContent("tags-a2", func(data PageData, routes *RouteSet) {
		routes.Register(Route{
			Name: "tag-pages-a2",
			Output: RouteOutput{Expand: func(buildData PageData) []Output {
				postsByTag := make(map[string][]Post, len(buildData.Tags))
				for _, tag := range buildData.Tags {
					postsByTag[tag] = nil
				}
				for _, post := range buildData.Posts {
					for _, tag := range post.Tags {
						if _, known := postsByTag[tag]; known {
							postsByTag[tag] = append(postsByTag[tag], post)
						}
					}
				}

				outputs := make([]Output, 0, len(buildData.Tags))
				for _, tag := range buildData.Tags {
					tag := tag
					posts := sortedTagPosts(postsByTag[tag])
					outputs = append(outputs, PageOutput(
						"tags/"+tag+"/index.html",
						func() interfaceComponent {
							return tagPage(tag, posts, buildData, renderHeadFragments(buildData))
						},
					))
				}
				return outputs
			}},
		})
	})
}

// TagPage renders one tag page for callers that need a component outside Build.
func TagPage(tag string, posts []Post, preview bool) interfaceComponent {
	return tagPage(tag, sortedTagPosts(posts), PageData{Preview: preview}, nil)
}

func sortedTagPosts(posts []Post) []Post {
	sorted := append([]Post(nil), posts...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Date != sorted[j].Date {
			return sorted[i].Date > sorted[j].Date
		}
		return sorted[i].Slug < sorted[j].Slug
	})
	return sorted
}

func a2PostURL(slug string) string {
	return "/journal/" + slug + "/"
}
