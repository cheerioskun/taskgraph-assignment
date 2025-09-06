package main

import (
	"errors"
)

var (
	ErrMissingStart    error = errors.New("missing_start")
	ErrMissingEnd      error = errors.New("missing_end")
	ErrMultipleStart   error = errors.New("multiple_start")
	ErrMultipleEnd     error = errors.New("multiple_end")
	ErrUnrunnableNodes error = errors.New("unrunnable_nodes")
	ErrIsolatedNodes   error = errors.New("isolated_nodes")
	ErrDuplicateEdge   error = errors.New("duplicate_edge")
	ErrOrphanedNodes   error = errors.New("orphaned_nodes")
)
