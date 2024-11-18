package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// documents
type Document struct {
	ID        uint `gorm:"primarykey"`
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DocID     string
	ChunkID   string
	DocText   string
	//Embedding []float64 `gorm:"type:float8[]"`
	Embedding string `gorm:"type:float8[]"` // Store as a string and ensure it's passed correctly
}

// documents_metadata
type DocumentMetadata struct {
	ID        int                    `json:"id" gorm:"primaryKey;autoIncrement"`
	DocID     string                 `json:"doc_id" gorm:"type:text;not null"` // Foreign key linking to documents
	Tags      pq.StringArray         `json:"tags" gorm:"type:text[];not null"` // Array of tags
	Metadata  map[string]interface{} `json:"metadata" gorm:"type:jsonb"`       // JSONB for additional metadata
	CreatedAt time.Time              `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// tag_embeddings
type TagEmbedding struct {
	TagName string `gorm:"column:tag_name"`
	//Embedding []float64 `gorm:"type:float8[];column:embedding"` // Postgres array type
	Embedding Float64Array `gorm:"type:float8[];column:embedding"`
}

// Float64Array is a custom type for handling float8 arrays in PostgreSQL.
type Float64Array []float64

// Scan implements the sql.Scanner interface to convert the database value to Float64Array.
func (f *Float64Array) Scan(src interface{}) error {
	if src == nil {
		*f = nil
		return nil
	}

	// Convert src to a string if necessary
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot convert %T to Float64Array", src)
	}

	// Trim the curly braces from the array string representation
	str = strings.Trim(str, "{}")

	// Split the string by commas
	elements := strings.Split(str, ",")
	*f = make(Float64Array, len(elements))

	// Parse each element as a float64
	for i, elem := range elements {
		val, err := strconv.ParseFloat(elem, 64)
		if err != nil {
			return fmt.Errorf("error parsing float64 from string %q: %w", elem, err)
		}
		(*f)[i] = val
	}

	return nil
}

// Value implements the driver.Valuer interface to convert Float64Array to a driver-compatible value.
func (f Float64Array) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}

	// Convert each float64 to a string
	strElems := make([]string, len(f))
	for i, val := range f {
		strElems[i] = strconv.FormatFloat(val, 'f', -1, 64)
	}

	// Join the elements with commas and wrap in curly braces
	return "{" + strings.Join(strElems, ",") + "}", nil
}
