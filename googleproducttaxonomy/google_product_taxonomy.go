package googleproducttaxonomy

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

//go:embed taxonomy-with-ids.en-US.txt
var taxonomyWithIDs []byte

var (
	Categories       map[ProductCategory]string
	categoryFromName map[string]ProductCategory
)

func init() {
	var err error

	Categories, err = LoadTaxonomy(bytes.NewReader(taxonomyWithIDs))
	if err != nil {
		slog.Error("cannot load taxonomy", "error", err)
	}

	categoryFromName = make(map[string]ProductCategory, len(Categories))
	for category, name := range Categories {
		categoryFromName[name] = category
	}
}

// ProductCategory is node in taxonomy tree
type ProductCategory int

func (s ProductCategory) IsZero() bool { return s == 0 }

func (s ProductCategory) Validate() error {
	if _, ok := Categories[s]; !ok {
		return &ErrNotFoundGoogleProductCategory{ID: int(s)}
	}
	return nil
}

func (s ProductCategory) String() string { return Categories[s] }

var Unknown ProductCategory

func ProductCategoryFromString(s string) (ProductCategory, error) {
	if id, ok := strconv.Atoi(s); ok == nil {
		id := ProductCategory(id)

		if err := id.Validate(); err == nil {
			return id, nil
		}
	}

	if category, ok := categoryFromName[s]; ok {
		return category, nil
	}

	return 0, errors.New("unknown")
}

func LoadTaxonomy(r io.Reader) (map[ProductCategory]string, error) {
	categories := make(map[ProductCategory]string, 5000)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " - ", 2)
		if len(parts) != 2 {
			return nil, errors.New("wrong line")
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, errors.New("wrong line")

		}

		categories[ProductCategory(id)] = strings.TrimSpace(parts[1])
	}

	return categories, scanner.Err()
}

type ErrNotFoundGoogleProductCategory struct {
	ID int
}

func (s *ErrNotFoundGoogleProductCategory) Error() string {
	return "not found google product category id: " + strconv.Itoa(s.ID)
}
