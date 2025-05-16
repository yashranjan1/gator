package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/database"
	"github.com/yashranjan1/gator/internal/rss"
	"github.com/yashranjan1/gator/internal/state"
)

func handleLogin(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}

	_, err := s.DataBase.GetUserByName(context.Background(), cmd.Args[0])

	if err == sql.ErrNoRows {
		return errors.New("user does not exist")
	} else if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func handleRegister(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the register handler expects a single argument, the username")
	}

	_, err := s.DataBase.GetUserByName(context.Background(), cmd.Args[0])

	if err == nil {
		return errors.New("a user with this name already exists")
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	userCreated, err := s.DataBase.CreateUser(context.Background(), params)

	err = s.Config.SetUser(userCreated.Name)

	if err != nil {
		return err
	}

	fmt.Printf("SUCCESS: user \"%s\" created, exiting....", userCreated.Name)
	return nil
}

func handleReset(s *state.State, cmd command.Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("the reset handler expects no arguments")
	}

	err := s.DataBase.Reset(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func handleList(s *state.State, cmd command.Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("the list handler expects no arguments")
	}

	users, err := s.DataBase.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user)
		if user == s.Config.CurrentUser {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func handleAggregate(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the aggregate handler expects 1 argument, a time interval")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every " + cmd.Args[0])

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handleAddFeed(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("the addfeed handler expects 2 arguments, name and url")
	}

	params := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	}

	feed, err := s.DataBase.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Feed added!")
	fmt.Printf("Name: %s\nUrl: %s\n", feed.Name, feed.Url)

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.DataBase.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}

	fmt.Println("Follow created for this feed")

	return nil
}

func handleListFeeds(s *state.State, cmd command.Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("the feeds handle expects no arguments")
	}

	feeds, err := s.DataBase.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Feed Name\tFeed URL\tUser Name")
	for _, feed := range feeds {
		fmt.Fprintf(w, "%s\t%s\t%s\n", feed.Feedname, feed.Url, feed.Username.String)
	}
	w.Flush()
	return nil
}

func handleFollow(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("the follow handler expects 1 argument, the url")
	}

	feed, err := s.DataBase.GetFeedByUrl(context.Background(), cmd.Args[0])

	if err == sql.ErrNoRows {
		return errors.New("this feed does not exist, please add the feed before you try to follow it")
	} else if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	added, err := s.DataBase.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Println("SUCCESS: The following feed has been added to the following user")
	fmt.Fprintln(w, "Feed Name\tUser Name")
	fmt.Fprintf(w, "%s\t%s\n", added.FeedName, added.UserName)
	w.Flush()

	return nil
}

func handleFollowing(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return errors.New("the following handler expects no arguments")
	}

	feeds, err := s.DataBase.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s follows the following feeds:\n", user.Name)
	for _, feed := range feeds {
		fmt.Println(feed.Name.String)
	}
	return nil
}

func handleUnfollow(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("the unfollow handler expects 1 argument, the url")
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	}

	err := s.DataBase.DeleteFeedFollow(context.Background(), params)

	if err != nil {
		return err
	}

	fmt.Println("feed successfully unfollowed!")

	return nil
}

func scrapeFeeds(s *state.State) {
	feed, err := s.DataBase.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("Error getting next feed: %v\n", err)
		return
	}

	markParams := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	}

	err = s.DataBase.MarkFeedFetched(context.Background(), markParams)
	if err != nil {
		fmt.Printf("Error in marking the feed as fetched: %v\n", err)
		return
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("Error in fetching the feed: %v\n", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		postParams := database.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: item.PubDate,
			FeedID:      feed.ID,
		}
		err := s.DataBase.CreatePosts(context.Background(), postParams)
		if err != nil {
			pqErr, ok := err.(*pq.Error)

			// this bit is specifically done so it ignore errors caused by the unique url constraint
			// because we are going to make calls that will cause duplicate insertions
			if ok {
				if pqErr.Code != "23505" {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
		fmt.Println(item.Title)
	}
}

func handleBrowse(s *state.State, cmd command.Command, user database.User) error {
	limit := 2

	if len(cmd.Args) > 0 {
		i, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return errors.New("the browse handler expects an integer argument")
		}

		limit = i
	}

	postsParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.DataBase.GetPostsForUser(context.Background(), postsParams)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published On: %s\n", post.PublishedAt)
		fmt.Println()
	}

	return nil
}
