package site

import "sort"

type relatedCandidateA2 struct {
	post  Post
	score int
}

// RelatedPosts returns the first three build-visible posts related to post.
// One point comes from each shared tag and one point from shared series
// membership. The score is local and deterministic; it is not a service.
// Ties use publish date and then slug, so map iteration cannot affect output.
func RelatedPosts(post Post, posts []Post) []Post {
	sharedTags := make(map[string]struct{}, len(post.Tags))
	for _, tag := range post.Tags {
		sharedTags[tag] = struct{}{}
	}

	candidates := make([]relatedCandidateA2, 0, len(posts))
	for _, candidate := range posts {
		if candidate.Slug == post.Slug {
			continue
		}
		score := 0
		countedTags := make(map[string]struct{}, len(candidate.Tags))
		for _, tag := range candidate.Tags {
			if _, counted := countedTags[tag]; counted {
				continue
			}
			countedTags[tag] = struct{}{}
			if _, shared := sharedTags[tag]; shared {
				score++
			}
		}
		if post.Series != "" && post.Series == candidate.Series {
			score++
		}
		if score > 0 {
			candidates = append(candidates, relatedCandidateA2{post: candidate, score: score})
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		if candidates[i].post.Date != candidates[j].post.Date {
			return candidates[i].post.Date > candidates[j].post.Date
		}
		return candidates[i].post.Slug < candidates[j].post.Slug
	})

	if len(candidates) > 3 {
		candidates = candidates[:3]
	}
	related := make([]Post, len(candidates))
	for i, candidate := range candidates {
		related[i] = candidate.post
	}
	return related
}

// relatedPosts is the package-local spelling for an article renderer in this
// package; RelatedPosts is exported for route and integration tests.
func relatedPosts(post Post, posts []Post) []Post {
	return RelatedPosts(post, posts)
}
