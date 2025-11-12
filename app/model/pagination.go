package model

type PageInfo struct {
  StartCursor string `json:"startCursor"`
  EndCursor string `json:"endCursor"`
  HasNextPage bool `json:"hasNextPage"`
  HasPreviousPage bool `json:"hasPreviousPage"`
}