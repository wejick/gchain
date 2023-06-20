package document

// Document is the data structure that will be the main building block of data processed by gochain
type Document struct {
	Text     string
	Metadata map[string]interface{}
}

// DocumentsToStrings convert array of document to array of string
func DocumentsToStrings(documents []Document) (output []string) {
	for _, doc := range documents {
		output = append(output, doc.Text)
	}
	return
}
