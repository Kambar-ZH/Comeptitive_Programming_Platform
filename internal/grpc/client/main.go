package main

// import (
// 	"context"
// 	"log"
// 	"site/internal/grpc/api"

// 	"google.golang.org/grpc"
// )

// const (
// 	port = ":8081"
// )

// func TestUser(ctx context.Context, userRepositoryClient api.UserRepositoryClient) {
// 	newUser := &api.User{
// 		Handle:    "Kambar",
// 		Country:   "Kazakhstan",
// 		City:      "Atyrau",
// 		Rating:    1703,
// 		MaxRating: 1703,
// 		Avatar:    "url-link",
// 	}

// 	createResponseUser, err := userRepositoryClient.Create(ctx, newUser)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("user successfully created to db: %v", createResponseUser)

// 	users, err := userRepositoryClient.All(ctx, &api.Pagination{})
// 	if err != nil {
// 		log.Fatalf("could not get users: %v", err)
// 	}
// 	log.Printf("got list of users: %v", users.Users)

// 	validHandle, invalidHandle := "Kambar", "Rasul"
// 	user, err := userRepositoryClient.ByHandle(ctx, &api.UserRequestHandle{Handle: validHandle})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("got user with handle %s: %v", validHandle, user)

// 	_, err = userRepositoryClient.ByHandle(ctx, &api.UserRequestHandle{Handle: invalidHandle})
// 	if err != nil {
// 		log.Printf("got error: %v", err)
// 	}

// 	newUser = &api.User{
// 		Handle:    "Yergeldi",
// 		Country:   "Kazakhstan",
// 		City:      "Ust Kamenogorsk",
// 		Rating:    2500,
// 		MaxRating: 3000,
// 		Avatar:    "url-link",
// 	}

// 	createResponseUser, err = userRepositoryClient.Create(ctx, newUser)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("user successfully created to db: %v", createResponseUser)

// 	updatedUser := &api.User{
// 		Handle:    "Aldiyar",
// 		Country:   "Kazakhstan",
// 		City:      "Atyrau",
// 		Rating:    1800,
// 		MaxRating: 1800,
// 		Avatar:    "url-link",
// 	}
// 	updateResponseUser, err := userRepositoryClient.Update(ctx, updatedUser)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("user successfully updated: %v", updateResponseUser)

// 	// _, err = userRepositoryClient.Delete(ctx, &api.UserRequestId{Id: 1})
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// }

// func TestSubmission(ctx context.Context, submissionRepositoryClient api.SubmissionRepositoryClient) {
// 	newSubmission := &api.Submission{
// 		Id:           1,
// 		Date:         "12.12.2021",
// 		AuthorHandle: "Kambar",
// 		ProblemId:    1,
// 		SubmissionResult: &api.SubmissionResult{
// 			Verdict:    api.VERDICT_PASSED,
// 			FailedTest: -1,
// 		},
// 	}
// 	createResponseSubmission, err := submissionRepositoryClient.Create(ctx, newSubmission)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("submission successfully created to db: %v", createResponseSubmission)
// 	submissions, err := submissionRepositoryClient.All(ctx, &api.Pagination{})
// 	if err != nil {
// 		log.Fatalf("could not get submissions: %v", err)
// 	}
// 	log.Printf("got list of submissions: %v", submissions.Submissions)

// 	validId, invalidId := int32(1), int32(3)
// 	submission, err := submissionRepositoryClient.ById(ctx, &api.SubmissionRequestId{Id: validId})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("got submission with id %d: %v", validId, submission)

// 	_, err = submissionRepositoryClient.ById(ctx, &api.SubmissionRequestId{Id: invalidId})
// 	if err != nil {
// 		log.Printf("got error: %v", err)
// 	}

// 	updatedSubmission := &api.Submission{
// 		Id:           1,
// 		Date:         "13.12.2021",
// 		AuthorHandle: "Yergeldi",
// 		ProblemId:    1,
// 		SubmissionResult: &api.SubmissionResult{
// 			Verdict:    api.VERDICT_FAILED,
// 			FailedTest: 3,
// 		},
// 	}
// 	updateResponseUser, err := submissionRepositoryClient.Update(ctx, updatedSubmission)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("submission successfully updated: %v", updateResponseUser)
// }

// func main() {
// 	ctx := context.Background()

// 	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
// 	if err != nil {
// 		log.Fatalf("could not connect to %s: %v", port, err)
// 	}

// 	// userRepositoryClient := api.NewUserRepositoryClient(conn)
// 	submissionRepositoryClient := api.NewSubmissionRepositoryClient(conn)

// 	// TestUser(ctx, userRepositoryClient)
// 	TestSubmission(ctx, submissionRepositoryClient)
// }
