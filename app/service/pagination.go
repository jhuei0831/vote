package service

import (
	"reflect"
	"strconv"
	"vote/app/model"
	"vote/app/utils"
)

type PaginationService[T any, E any, C any] struct {
}

func NewPaginationService[T any, E any, C any]() PaginationService[T, E, C] {
	return PaginationService[T, E, C]{}
}

// BuildConnection creates edges and pageInfo, then returns a specific type of Connection structure.
func (p PaginationService[T, E, C]) BuildConnection(
	records []T,
	total int64,
	hasPreviousPage bool,
	hasNextPage bool,
	getID func(T) uint64,
) C {
	edges := make([]E, 0, len(records))
	var firstCursor string
	var lastCursor string

	for i, record := range records {
		cursor, _ := (&utils.Invitation{}).Encrypt(strconv.FormatUint(getID(record), 10))
		edges = append(edges, buildEdge[T, E](record, cursor))

		if i == 0 {
			firstCursor = cursor
		}
		lastCursor = cursor
	}

	pageInfo := buildPageInfo(firstCursor, lastCursor, hasPreviousPage, hasNextPage, len(edges) > 0)

	return buildConnectionValue[E, C](edges, pageInfo, total)
}

func buildEdge[T any, E any](record T, cursor string) E {
	val, isPointer := newValue[E]()
	target := val.Elem()

	nodeField := target.FieldByName("Node")
	if nodeField.IsValid() && nodeField.CanSet() {
		nodeField.Set(reflect.ValueOf(record))
	}

	cursorField := target.FieldByName("Cursor")
	if cursorField.IsValid() && cursorField.CanSet() {
		cursorField.SetString(cursor)
	}

	if isPointer {
		return val.Interface().(E)
	}

	return target.Interface().(E)
}

func buildConnectionValue[E any, C any](edges []E, pageInfo model.PageInfo, total int64) C {
	val, isPointer := newValue[C]()
	target := val.Elem()

	edgesField := target.FieldByName("Edges")
	if edgesField.IsValid() && edgesField.CanSet() {
		edgesField.Set(reflect.ValueOf(edges))
	}

	pageInfoField := target.FieldByName("PageInfo")
	if pageInfoField.IsValid() && pageInfoField.CanSet() {
		pageInfoField.Set(reflect.ValueOf(pageInfo))
	}

	totalField := target.FieldByName("TotalCount")
	if totalField.IsValid() && totalField.CanSet() {
		totalField.SetInt(total)
	}

	if isPointer {
		return val.Interface().(C)
	}

	return target.Interface().(C)
}

func buildPageInfo(firstCursor, lastCursor string, hasPreviousPage, hasNextPage, hasEdges bool) model.PageInfo {
	if !hasEdges {
		return model.PageInfo{
			HasNextPage:     false,
			HasPreviousPage: false,
		}
	}

	return model.PageInfo{
		StartCursor:     firstCursor,
		EndCursor:       lastCursor,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
}

func newValue[T any]() (reflect.Value, bool) {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ == nil {
		panic("generic type cannot be nil")
	}

	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()), true
	}

	return reflect.New(typ), false
}
