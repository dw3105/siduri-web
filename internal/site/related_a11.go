package site

import "github.com/a-h/templ"

func init() {
	RegisterArticleSection(ArticleSection{
		Name: "a11-related-posts",
		Slot: ArticleSectionAfterBody,
		Render: func(post Post, data PageData) templ.Component {
			return relatedPostsSection(RelatedPosts(post, data.Posts))
		},
	})
}
