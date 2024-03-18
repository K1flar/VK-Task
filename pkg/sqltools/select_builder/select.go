package selectbuilder

import (
	"film_library/pkg/pagination"
	"fmt"
	"strings"
)

type SelectQueryBuilder struct {
	query      *strings.Builder
	joins      []string
	conditions []string
	orders     []string
	pagination *pagination.Pagination
}

func New(selectQuery string) *SelectQueryBuilder {
	query := &strings.Builder{}
	query.WriteString(selectQuery + " ")

	return &SelectQueryBuilder{
		query:      query,
		joins:      []string{},
		conditions: []string{},
		orders:     []string{},
	}
}

func (b *SelectQueryBuilder) Join(joinQuery string) *SelectQueryBuilder {
	b.joins = append(b.joins, " JOIN "+joinQuery+" ")
	return b
}

func (b *SelectQueryBuilder) LeftJoin(joinQuery string) *SelectQueryBuilder {
	b.joins = append(b.joins, "LEFT JOIN "+joinQuery+" ")
	return b
}

func (b *SelectQueryBuilder) Where(conditionFormat string, params ...any) *SelectQueryBuilder {
	b.conditions = append(b.conditions, fmt.Sprintf(conditionFormat, params...))
	return b
}

func (b *SelectQueryBuilder) OrderBy(orderBy, direction string) *SelectQueryBuilder {
	if strings.ToLower(direction) != "asc" && strings.ToLower(direction) != "desc" {
		direction = "asc"
	}
	b.orders = append(b.orders, fmt.Sprintf("%s %s", orderBy, direction))
	return b
}

func (b *SelectQueryBuilder) AddPagination(pagination *pagination.Pagination) *SelectQueryBuilder {
	b.pagination = pagination
	return b
}

func (b *SelectQueryBuilder) Build() string {
	for _, join := range b.joins {
		b.query.WriteString(join)
	}

	if len(b.conditions) != 0 {
		b.query.WriteString(" WHERE " + b.conditions[0])
	}
	for i := 1; i < len(b.conditions); i++ {
		b.query.WriteString(" AND " + b.conditions[i])
	}

	if len(b.orders) != 0 {
		b.query.WriteString(" ORDER BY " + b.orders[0])
	}
	for i := 1; i < len(b.orders); i++ {
		b.query.WriteString(", " + b.orders[i])
	}

	b.query.WriteString(fmt.Sprintf(" LIMIT %d", b.pagination.GetLimit()))
	b.query.WriteString(fmt.Sprintf(" OFFSET %d", b.pagination.GetOffset()))

	return b.query.String()
}
