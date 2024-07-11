package mobi

import (
	"strings"

	r "github.com/adamay909/AozoraConvert/mobi/records"
)

/*
	func chaptersToText(m Book) (string, []r.ChunkInfo, []r.ChapterInfo, error) {
		text := new(strings.Builder)
		chunks := make([]r.ChunkInfo, 0)
		chaps := make([]r.ChapterInfo, 0)
		if m.tpl == nil {
			m.tpl = defaultTemplate
		}

		chunkId := 0
		for chapId, chap := range m.Chapters {
			chapStart := text.Len()
			for _, chunk := range chap.Chunks {
				inv := newInventory(m, chap, chapId, chunkId)
				head, err := runTemplate(*m.tpl, inv)
				if err != nil {
					return "", nil, nil, err
				}
				chunks = append(chunks, r.ChunkInfo{
					PreStart:      text.Len(),
					PreLength:     len(head),
					ContentStart:  text.Len() + len(head),
					ContentLength: len(chunk.Body),
				})
				text.WriteString(head)
				text.WriteString(chunk.Body)
				chunkId++
			}
			chaps = append(chaps, r.ChapterInfo{
				Title:  chap.Title,
				Start:  chapStart,
				Length: text.Len() - chapStart,
			})
		}

		return text.String(), chunks, chaps, nil
	}
*/
func chaptersToText2(m Book) (html string, chunks []r.ChunkInfo, chaps []r.ChapterInfo, err error) {

	if m.tpl == nil {
		m.tpl = defaultTemplate
	}
	text := new(strings.Builder)

	inv := newInventory(m)

	head, err := runTemplate(*m.tpl, inv)
	if err != nil {
		return "", nil, nil, err
	}

	chunks = append(chunks, r.ChunkInfo{
		PreStart:      text.Len(),
		PreLength:     len(head),
		ContentStart:  text.Len() + len(head),
		ContentLength: len(m.Html),
	})

	text.WriteString(head)
	text.WriteString(m.Html)
	html = text.String()

	for _, chap := range m.Chapters {
		chaps = append(chaps, r.ChapterInfo{
			Title:  chap.Title,
			Start:  chap.Start + len(head),
			Length: chap.Length,
		})
	}

	return
}

func textToRecords(html string, chapters []r.ChapterInfo) []r.TextRecord {
	//		provider := r.NewTrailProvider(chapters)
	records := make([]r.TextRecord, 0)
	recordCount := len(html) / r.TextRecordMaxSize
	if len(html)%r.TextRecordMaxSize != 0 {
		recordCount++
	}

	for i := 0; i < recordCount; i++ {
		from := i * r.TextRecordMaxSize
		to := min(from+r.TextRecordMaxSize, len(html))
		trail := r.Get(len(html), from, to)
		records = append(records, r.NewTextRecord(html[from:to], trail))
	}

	return records
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
