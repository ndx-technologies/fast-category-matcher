package fastcategorymatcher

import (
	"errors"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"

	"github.com/ndx-technologies/fast-category-matcher/googleproducttaxonomy"
	"github.com/ndx-technologies/sift4"
)

type ApproxStrGoogleTaxonomyMatcherConfig struct {
	MaxNodeDistance int     `json:"max_node_distance"`
	MinNodeLength   int     `json:"min_node_length"`
	MinScore        float32 `json:"min_score"`
}

func (s ApproxStrGoogleTaxonomyMatcherConfig) WithDefaults() ApproxStrGoogleTaxonomyMatcherConfig {
	if s.MaxNodeDistance == 0 {
		s.MaxNodeDistance = 5
	}
	if s.MinNodeLength == 0 {
		s.MinNodeLength = 6
	}
	if s.MinScore == 0 {
		s.MinScore = 0.75
	}
	return s
}

// ApproxStrGoogleTaxonomyMatcher finds closest Google Taxonomy Category based on its approximate string representation.
// It utilizes the fact that at least some nodes may be correct.
// This is effectivelly closest entry with very basic typo-correction.
// This approximator has very small memory footprint for invocations and latency ~5ms.
// This allows to run this analysis in-flight real time request.
// In turn, this allows user to see most likely parsed category in real time from noisy source (e.g. LLM).
type ApproxStrGoogleTaxonomyMatcher struct {
	config ApproxStrGoogleTaxonomyMatcherConfig

	// category part index
	partName      []string
	partNameStem  []string
	categoryParts map[googleproducttaxonomy.ProductCategory][]int
}

func NewApproxStrGoogleTaxonomyMatcher(config ApproxStrGoogleTaxonomyMatcherConfig) *ApproxStrGoogleTaxonomyMatcher {
	partFromName := make(map[string]int, len(googleproducttaxonomy.Categories)*2)
	partName := make([]string, 0, len(googleproducttaxonomy.Categories)*2)
	partNameStem := make([]string, 0, len(googleproducttaxonomy.Categories)*2)
	categoryParts := make(map[googleproducttaxonomy.ProductCategory][]int, len(googleproducttaxonomy.Categories))

	for category, path := range googleproducttaxonomy.Categories {
		for p := range strings.SplitSeq(path, ">") {
			p = strings.TrimSpace(p)
			if _, ok := partFromName[p]; !ok {
				partName = append(partName, p)
				partNameStem = append(partNameStem, normalizePhrase(p))
				partFromName[p] = len(partName) - 1
			}
			categoryParts[category] = append(categoryParts[category], partFromName[p])
		}
	}

	return &ApproxStrGoogleTaxonomyMatcher{
		config:        config,
		partName:      partName,
		partNameStem:  partNameStem,
		categoryParts: categoryParts,
	}
}

func removeNonLetters(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func normalizePhrase(v string) string {
	words := strings.Fields(v)
	for i, word := range words {
		word = strings.ToLower(word)
		word = removeNonLetters(word)

		if len(word) == 0 {
			continue
		}

		words[i] = english.Stem(word, false)
	}
	return strings.Join(words, " ")
}

func (s *ApproxStrGoogleTaxonomyMatcher) MatchGoogleProductCategory(v string) (googleproducttaxonomy.ProductCategory, error) {
	if v == "" {
		return 0, errors.New("empty input")
	}

	parts := strings.Split(v, ">")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	if len(parts) == 0 {
		return 0, errors.New("empty input")
	}

	scores := make([]float32, len(parts))
	var sift4Buffer sift4.Buffer

	var bestCategory googleproducttaxonomy.ProductCategory
	var maxScore float32

	queryStems := make([]string, len(parts))
	for i, part := range parts {
		queryStems[i] = normalizePhrase(part)
	}

	for category, categoryParts := range s.categoryParts {
		score := s.score(parts, categoryParts, queryStems, scores, &sift4Buffer)

		if score > maxScore {
			maxScore = score
			bestCategory = category
		}
	}

	if maxScore < 0 {
		return 0, errors.New("no matching category found")
	}

	return bestCategory, nil
}

func (s *ApproxStrGoogleTaxonomyMatcher) score(queryParts []string, categoryParts []int, queryStems []string, scores []float32, buf *sift4.Buffer) float32 {
	for i := range scores {
		scores[i] = 0
	}

	for iq, qpart := range queryParts {
		for _, partID := range categoryParts {
			// last entry query category in cannot be non-last entry in category
			if iq == len(queryParts)-1 {
				if partID != categoryParts[len(categoryParts)-1] {
					continue
				}
			}

			if qpart == s.partName[partID] {
				scores[iq] = 1
				continue
			}

			// stem will disregard slight variations of the same word (e.g. plurals)
			if queryStems[iq] == s.partNameStem[partID] {
				scores[iq] = 1
				continue
			}

			// target is too small, typos in small words are unlikely and they change meaning of word too much. chance of mistakes is high
			if len(qpart) < s.config.MinNodeLength || len(s.partName[partID]) < s.config.MinNodeLength {
				continue
			}

			// caching similarity is not effective, it requires 90KB of memory and does not increase performance by much
			// increment by 1 so that distance can breach max distance
			d := sift4.Distance(qpart, s.partName[partID], 256, s.config.MaxNodeDistance+1, buf)

			if d > s.config.MaxNodeDistance {
				continue
			}

			// ratio of matched chars to target works better than source
			score := 1 - float32(d)/float32(len(s.partName[partID]))

			if score < s.config.MinScore {
				continue
			}

			if score > scores[iq] {
				scores[iq] = score
			}
		}
	}

	var total float32
	var matchCount int
	for _, q := range scores {
		if q > 0 {
			total += q
			matchCount++
		}
	}
	if matchCount == 0 {
		return 0
	}

	// we cannot take average of matched parts, because this will lead to score 1 for queries and targets that do not even have many common parts
	// not taking len of target parts, because query often missing some parts

	return float32(total) / float32(len(queryParts))
}
