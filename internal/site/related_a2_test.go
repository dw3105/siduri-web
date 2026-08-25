package site

import "testing"

func TestA2RelatedPostsAreThreeDeterministicAndSelfExcluding(t *testing.T) {
	current := Post{Slug: "current", Tags: []string{"method", "build-log"}, Series: "series-a"}
	posts := []Post{
		current,
		{Slug: "two-tags", Date: "2026-08-24", Tags: []string{"method", "build-log"}},
		{Slug: "one-tag-series", Date: "2026-08-23", Tags: []string{"method"}, Series: "series-a"},
		{Slug: "one-tag", Date: "2026-08-22", Tags: []string{"build-log"}},
		{Slug: "series-only", Date: "2026-08-21", Series: "series-a"},
		{Slug: "unrelated", Date: "2026-08-25", Tags: []string{"outcome"}},
	}

	related := RelatedPosts(current, posts)
	if len(related) != 3 {
		t.Fatalf("related count = %d, want 3", len(related))
	}
	want := []string{"two-tags", "one-tag-series", "one-tag"}
	for i, post := range related {
		if post.Slug != want[i] {
			t.Fatalf("related[%d] = %q, want %q", i, post.Slug, want[i])
		}
		if post.Slug == current.Slug {
			t.Fatal("post related to itself")
		}
	}

	for attempt := 0; attempt < 20; attempt++ {
		repeated := RelatedPosts(current, posts)
		for i := range repeated {
			if repeated[i].Slug != want[i] {
				t.Fatalf("run %d related[%d] = %q, want %q", attempt, i, repeated[i].Slug, want[i])
			}
		}
	}
}

func TestA2RelatedPostsReturnFewerThanThreeWithoutPadding(t *testing.T) {
	current := Post{Slug: "current", Tags: []string{"method"}}
	posts := []Post{
		current,
		{Slug: "same-tag", Date: "2026-08-24", Tags: []string{"method"}},
		{Slug: "same-series", Date: "2026-08-23", Series: "series-a"},
	}

	related := RelatedPosts(current, posts)
	if len(related) != 1 {
		t.Fatalf("related count = %d, want 1", len(related))
	}
	if related[0].Slug != "same-tag" {
		t.Fatalf("related[0] = %q, want same-tag", related[0].Slug)
	}
}

func TestA2RelatedPostsTieBreakDateThenSlug(t *testing.T) {
	current := Post{Slug: "current", Tags: []string{"method"}}
	posts := []Post{
		{Slug: "zeta", Date: "2026-08-20", Tags: []string{"method"}},
		{Slug: "alpha", Date: "2026-08-20", Tags: []string{"method"}},
		current,
	}

	related := relatedPosts(current, posts)
	if len(related) != 2 {
		t.Fatalf("related count = %d, want 2", len(related))
	}
	if related[0].Slug != "alpha" || related[1].Slug != "zeta" {
		t.Fatalf("tie order = %q, %q; want alpha, zeta", related[0].Slug, related[1].Slug)
	}
}
