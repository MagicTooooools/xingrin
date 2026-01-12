package scope

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// LogicalOp represents logical operator
type LogicalOp int

const (
	LogicalAnd LogicalOp = iota
	LogicalOr
)

// ParsedFilter represents a parsed filter condition
type ParsedFilter struct {
	Field    string
	Operator string // "=", "==", "!="
	Value    string
}

// FilterGroup represents a filter with logical operator
type FilterGroup struct {
	Filter    ParsedFilter
	LogicalOp LogicalOp
}

// FieldConfig represents field configuration for filtering
type FieldConfig struct {
	Column  string // Database column name
	IsArray bool   // Whether it's a PostgreSQL array field
}

// FilterMapping is a map of field name to field config
type FilterMapping map[string]FieldConfig

// filterPattern matches: field="value", field=="value", field!="value"
var filterPattern = regexp.MustCompile(`(\w+)(==|!=|=)"([^"]*)"`)

// orPattern matches || or "or" (case insensitive, not part of word)
var orPattern = regexp.MustCompile(`\s*(\|\||\bor\b)\s*`)

// andPattern matches && or "and" (case insensitive, not part of word)
var andPattern = regexp.MustCompile(`\s*(&&|\band\b)\s*`)

// WithFilter returns a GORM scope for smart filtering
// Supports:
//   - field="value"     fuzzy match (ILIKE)
//   - field=="value"    exact match
//   - field!="value"    not equal
//   - || or "or"        OR logic
//   - && or "and" or space   AND logic
//   - plain text        fuzzy match on default field (if defaultField provided)
func WithFilter(filterStr string, mapping FilterMapping) func(db *gorm.DB) *gorm.DB {
	return WithFilterDefault(filterStr, mapping, "")
}

// WithFilterDefault returns a GORM scope for smart filtering with a default field
// If filterStr is plain text (no operators), it will be treated as fuzzy search on defaultField
func WithFilterDefault(filterStr string, mapping FilterMapping, defaultField string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filterStr == "" || len(mapping) == 0 {
			return db
		}

		// Check if it's plain text (no filter syntax)
		if defaultField != "" && !hasFilterSyntax(filterStr) {
			// Treat as fuzzy search on default field
			config, ok := mapping[defaultField]
			if ok {
				return db.Where(config.Column+" ILIKE ?", "%"+filterStr+"%")
			}
		}

		groups := parseFilter(filterStr)
		if len(groups) == 0 {
			return db
		}

		return buildQuery(db, groups, mapping)
	}
}

// hasFilterSyntax checks if the string contains filter syntax (field="value")
func hasFilterSyntax(s string) bool {
	return filterPattern.MatchString(s)
}

// parseFilter parses filter string into FilterGroup list
func parseFilter(filterStr string) []FilterGroup {
	if filterStr == "" {
		return nil
	}

	// Step 1: Extract all filter conditions and replace with placeholders
	var filtersFound []string
	placeholder := "__FILTER_%d__"

	protected := filterPattern.ReplaceAllStringFunc(filterStr, func(match string) string {
		idx := len(filtersFound)
		filtersFound = append(filtersFound, match)
		return strings.Replace(placeholder, "%d", string(rune('0'+idx)), 1)
	})

	// Handle more than 10 filters
	for i := 0; i < len(filtersFound); i++ {
		oldPlaceholder := strings.Replace(placeholder, "%d", string(rune('0'+i)), 1)
		newPlaceholder := "__FILTER_" + string(rune('0'+i)) + "__"
		if i >= 10 {
			// For indices >= 10, use different approach
			continue
		}
		protected = strings.Replace(protected, oldPlaceholder, newPlaceholder, 1)
	}

	// Step 2: Normalize logical operators
	normalized := orPattern.ReplaceAllString(protected, " __OR__ ")
	normalized = andPattern.ReplaceAllString(normalized, " __AND__ ")

	// Step 3: Tokenize
	tokens := strings.Fields(normalized)

	var groups []FilterGroup
	pendingOp := LogicalAnd // Default AND

	for _, token := range tokens {
		switch token {
		case "__OR__":
			pendingOp = LogicalOr
		case "__AND__":
			pendingOp = LogicalAnd
		default:
			if strings.HasPrefix(token, "__FILTER_") && strings.HasSuffix(token, "__") {
				// Extract index
				idxStr := token[9 : len(token)-2]
				if len(idxStr) == 1 {
					idx := int(idxStr[0] - '0')
					if idx >= 0 && idx < len(filtersFound) {
						match := filterPattern.FindStringSubmatch(filtersFound[idx])
						if len(match) == 4 {
							logicalOp := LogicalAnd
							if len(groups) > 0 {
								logicalOp = pendingOp
							}
							groups = append(groups, FilterGroup{
								Filter: ParsedFilter{
									Field:    strings.ToLower(match[1]),
									Operator: match[2],
									Value:    match[3],
								},
								LogicalOp: logicalOp,
							})
							pendingOp = LogicalAnd // Reset to default
						}
					}
				}
			}
		}
	}

	return groups
}

// buildQuery builds GORM query from filter groups
func buildQuery(db *gorm.DB, groups []FilterGroup, mapping FilterMapping) *gorm.DB {
	if len(groups) == 0 {
		return db
	}

	// Build combined condition
	var conditions []string
	var args []interface{}
	var orGroups [][]int // Track which conditions are OR'd together

	currentOrGroup := []int{}

	for i, group := range groups {
		config, ok := mapping[group.Filter.Field]
		if !ok {
			continue
		}

		cond, arg := buildSingleCondition(config, group.Filter)
		if cond == "" {
			continue
		}

		conditions = append(conditions, cond)
		args = append(args, arg...)

		if group.LogicalOp == LogicalOr && i > 0 {
			// This condition is OR'd with previous
			if len(currentOrGroup) == 0 {
				currentOrGroup = append(currentOrGroup, len(conditions)-2)
			}
			currentOrGroup = append(currentOrGroup, len(conditions)-1)
		} else if len(currentOrGroup) > 0 {
			orGroups = append(orGroups, currentOrGroup)
			currentOrGroup = []int{}
		}
	}

	if len(currentOrGroup) > 0 {
		orGroups = append(orGroups, currentOrGroup)
	}

	if len(conditions) == 0 {
		return db
	}

	// Simple case: no OR groups, just AND everything
	if len(orGroups) == 0 {
		for i, cond := range conditions {
			db = db.Where(cond, args[i])
		}
		return db
	}

	// Complex case: handle OR groups
	// Build the query with proper grouping
	return buildComplexQuery(db, groups, mapping)
}

// buildComplexQuery handles queries with OR logic
func buildComplexQuery(db *gorm.DB, groups []FilterGroup, mapping FilterMapping) *gorm.DB {
	if len(groups) == 0 {
		return db
	}

	// Process groups sequentially, combining with appropriate logic
	var currentExpr *gorm.DB
	isFirst := true

	for _, group := range groups {
		config, ok := mapping[group.Filter.Field]
		if !ok {
			continue
		}

		cond, args := buildSingleCondition(config, group.Filter)
		if cond == "" {
			continue
		}

		if isFirst {
			currentExpr = db.Where(cond, args...)
			isFirst = false
		} else if group.LogicalOp == LogicalOr {
			currentExpr = currentExpr.Or(cond, args...)
		} else {
			currentExpr = currentExpr.Where(cond, args...)
		}
	}

	if currentExpr == nil {
		return db
	}

	return currentExpr
}

// buildSingleCondition builds a single WHERE condition
func buildSingleCondition(config FieldConfig, filter ParsedFilter) (string, []interface{}) {
	column := config.Column

	if config.IsArray {
		return buildArrayCondition(column, filter)
	}

	switch filter.Operator {
	case "==":
		// Exact match
		return column + " = ?", []interface{}{filter.Value}
	case "!=":
		// Not equal
		return column + " != ?", []interface{}{filter.Value}
	default: // "="
		// Fuzzy match (ILIKE for case-insensitive)
		return column + " ILIKE ?", []interface{}{"%" + filter.Value + "%"}
	}
}

// buildArrayCondition builds condition for PostgreSQL array fields
func buildArrayCondition(column string, filter ParsedFilter) (string, []interface{}) {
	switch filter.Operator {
	case "==":
		// Exact match: array contains this exact element
		return "? = ANY(" + column + ")", []interface{}{filter.Value}
	case "!=":
		// Not contains
		return "NOT (? = ANY(" + column + "))", []interface{}{filter.Value}
	default: // "="
		// Fuzzy match: convert array to string and search
		return "array_to_string(" + column + ", ',') ILIKE ?", []interface{}{"%" + filter.Value + "%"}
	}
}
