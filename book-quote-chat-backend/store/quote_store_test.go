package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestQuoteStore(t *testing.T) {
	_ = os.Remove("test_quotes.json")
	quoteFile = "test_quotes.json"

	q := model.Quote{ID: "id1", Text: "hello", Book: "test", User: "aric", Created: 1234567}
	if err := SaveQuotes([]model.Quote{q}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadQuotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "id1" {
		t.Fatal("not saved or loaded")
	}
	_ = os.Remove("test_quotes.json")
}
