
package model

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation"
)

type Tag struct {
	Name string
}

func (t Tag) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required))
}

type Article struct {
	Title string
	Body  string
	Tags  []Tag
}

func (a Article) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Body, validation.Required),
		validation.Field(&a.Tags, validation.Required))
}

// TestArticleValidate is a unit test for Article Validate method
func TestArticleValidate(t *testing.T) {
	type scenario struct {
		name            string
		article         Article
		expectingErrors bool
		expectedErrors  map[string]string
	}

	scenarios := []scenario{
		{
			name:            "Scenario 1: Test with a valid article",
			article:         Article{Title: "GoLang Unit Testing", Body: "Best practices and strategies", Tags: []Tag{{Name: "Go"}}},
			expectingErrors: false,
			expectedErrors:  nil,
		},
		{
			name:            "Scenario 2: Test with an article lacking title",
			article:         Article{Body: "Some article body", Tags: []Tag{{Name: "Go"}}},
			expectingErrors: true,
			expectedErrors:  map[string]string{"Title": "cannot be blank"},
		},
		{
			name:            "Scenario 3: Test with an article lacking body",
			article:         Article{Title: "Article", Tags: []Tag{{Name: "Go"}}},
			expectingErrors: true,
			expectedErrors:  map[string]string{"Body": "cannot be blank"},
		},
		{
			name:            "Scenario 4: Test with an article lacking tags",
			article:         Article{Title: "Article", Body: "Some article body"},
			expectingErrors: true,
			expectedErrors:  map[string]string{"Tags": "cannot be blank"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v", r)
					t.Fail()
				}
			}()
			err := s.article.Validate()
			if s.expectingErrors {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				for key, val := range s.expectedErrors {
					if err.Error() != key+" "+val {
						t.Fatalf("expected error %q but got %q", key+" "+val, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			}
		})
	}
}
