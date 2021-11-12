package postgres

import (
	"context"
	"fmt"
	"log"
	"site/internal/datastruct"
	"testing"
)

func TestSubmissionRepository(t *testing.T) {
	store := NewDB()
	if err := store.Connect(urlAddress); err != nil {
		log.Println(err)
		return
	}
	defer store.Close()
	ctx := context.Background()

	// CREATE
	err := store.Submissions().Create(ctx, &datastruct.Submission{
		ContestId:    1,
		AuthorHandle: "Kambar_Z",
		Verdict:      "Passed",
	})
	if err != nil {
		panic(err)
	}
	// CREATE

	// ALL
	submissions, err := store.Submissions().All(ctx, 0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(submissions)
	// ALL

	// BYAUTHORHANDLE
	submissions, err = store.Submissions().ByAuthorHandle(ctx, "Kambar_Z")
	if err != nil {
		panic(err)
	}
	fmt.Println(submissions)
	// BYAUTHORHANDLE

	// BYCONTESTID
	submissions, err = store.Submissions().ByContestId(ctx, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(submissions)
	// BYCONTESTID

	// BYID
	submission, err := store.Submissions().ById(ctx, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(submission)
	// BYID

	// DELETE
	if err := store.Submissions().Delete(ctx, 1); err != nil {
		panic(err)
	}
	// DELETE
}