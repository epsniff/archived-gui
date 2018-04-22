package bleveimpl

type SearchResult struct {
	Hits   int
	DocIDs []string
}

func NewSearchResult() *SearchResult {
	return &SearchResult{
		Hits:   0,
		DocIDs: make([]string, 0, 0),
	}
}

func (sr *SearchResult) AddHit(docId string) {
	sr.Hits++
	sr.DocIDs = append(sr.DocIDs, docId)
}
