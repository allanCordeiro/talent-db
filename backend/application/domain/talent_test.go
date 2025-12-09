package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateSuccess(t *testing.T) {
	talent, err := Create(
		"https://linkedin.com/in/john",
		"Senior Developer",
		"John Doe",
		"Experienced developer",
		"Tech Corp",
		"Developer",
		[]string{"Go", "Rust"},
		"Great candidate",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if talent == nil {
		t.Fatal("expected talent to be not nil")
	}

	if talent.ProfileURL != "https://linkedin.com/in/john" {
		t.Errorf("expected ProfileURL to be 'https://linkedin.com/in/john', got %s", talent.ProfileURL)
	}

	if talent.FullName != "John Doe" {
		t.Errorf("expected FullName to be 'John Doe', got %s", talent.FullName)
	}

	if talent.Id == uuid.Nil {
		t.Error("expected Id to be generated, got nil uuid")
	}

	if len(talent.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(talent.Tags))
	}
}

func TestCreateMissingProfileURL(t *testing.T) {
	talent, err := Create("", "Senior Developer", "John Doe", "Headline", "Company", "Role", []string{}, "Notes")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if talent != nil {
		t.Fatal("expected talent to be nil")
	}

	if err.Error() != "url is null" {
		t.Errorf("expected error message 'url is null', got %s", err.Error())
	}
}

func TestCreateMissingRole(t *testing.T) {
	talent, err := Create("https://linkedin.com/in/john", "", "John Doe", "Headline", "Company", "Role", []string{}, "Notes")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if talent != nil {
		t.Fatal("expected talent to be nil")
	}

	if err.Error() != "role is null" {
		t.Errorf("expected error message 'role is null', got %s", err.Error())
	}
}

func TestCreateMissingFullName(t *testing.T) {
	talent, err := Create("https://linkedin.com/in/john", "Senior Developer", "", "Headline", "Company", "Role", []string{}, "Notes")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if talent != nil {
		t.Fatal("expected talent to be nil")
	}

	if err.Error() != "name is null" {
		t.Errorf("expected error message 'name is null', got %s", err.Error())
	}
}

func TestCreateMissingHeadline(t *testing.T) {
	talent, err := Create("https://linkedin.com/in/john", "Senior Developer", "John Doe", "", "Company", "Role", []string{}, "Notes")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if talent != nil {
		t.Fatal("expected talent to be nil")
	}

	if err.Error() != "headline is null" {
		t.Errorf("expected error message 'headline is null', got %s", err.Error())
	}
}

func TestValidateSuccess(t *testing.T) {
	talent := &Talent{
		ProfileURL:   "https://linkedin.com/in/john",
		PossibleRole: "Developer",
		FullName:     "John Doe",
		Headline:     "Experienced",
	}

	err := talent.Validate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCapturedAtTimestamp(t *testing.T) {
	beforeCreate := time.Now().UTC()
	talent, _ := Create("https://test.com", "Dev", "Test", "Headline", "Company", "Role", []string{}, "Notes")
	afterCreate := time.Now().UTC()

	if talent.CapturedAt.Before(beforeCreate) || talent.CapturedAt.After(afterCreate) {
		t.Errorf("expected CapturedAt between %v and %v, got %v", beforeCreate, afterCreate, talent.CapturedAt)
	}
}

func TestCreateWithEmptyTags(t *testing.T) {
	talent, err := Create("https://test.com", "Dev", "Name", "Headline", "Company", "Role", []string{}, "Notes")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(talent.Tags) != 0 {
		t.Errorf("expected empty tags, got %v", talent.Tags)
	}
}
