package faker

import "strings"

// LoremGen — namespace faker.Lorem, аналог faker.lorem из faker.js.
type LoremGen struct{ f *Faker }

func (l *LoremGen) Word() string {
	return PickOne(l.f, loremWords)
}

func (l *LoremGen) Words(n int) string {
	if n <= 0 {
		return ""
	}
	words := make([]string, n)
	for i := range words {
		words[i] = l.Word()
	}
	return strings.Join(words, " ")
}

func (l *LoremGen) Sentence() string {
	n := l.f.IntRange(5, 12)
	s := l.Words(n)
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:] + "."
}

func (l *LoremGen) Sentences(n int) string {
	if n <= 0 {
		return ""
	}
	sents := make([]string, n)
	for i := range sents {
		sents[i] = l.Sentence()
	}
	return strings.Join(sents, " ")
}

func (l *LoremGen) Paragraph() string {
	return l.Sentences(l.f.IntRange(3, 6))
}

// Paragraphs генерирует n параграфов, соединённых sep (например "\n\n").
func (l *LoremGen) Paragraphs(n int, sep string) string {
	if n <= 0 {
		return ""
	}
	ps := make([]string, n)
	for i := range ps {
		ps[i] = l.Paragraph()
	}
	return strings.Join(ps, sep)
}
