# Gator

Gator is a cli utility for aggregating RSS feeds.

## Requirements

To use Gator you will need the following:

- Go >= v1.24.2
- PostgreSQL >= v17
- [Goose](https://github.com/pressly/goose)

## Installation

Copy, paste and run the following command in your terminal to install gator

```
go install https://github.com/yashranjan1/gator
```

## Setup

You'll need to create a file called `.gatorconfig.json` in your `home`
directory. The contents of the file should look like this:

```
{
    "db_url": { ENTER YOU DB CONN STRING HERE },
    "current_user": ""
}
```

Once this is done, run the following command from the `./sql/schema/` directory
to apply your migrations to the database.

```
goose postgres {YOUR CONN STRING} up
```

## Usage

Gator has a list of commands you can use to aggregate your RSS feeds!

| command   | Description                                             | Argument                 | Example                                                        |
| --------- | ------------------------------------------------------- | ------------------------ | -------------------------------------------------------------- |
| register  | Creates a user by the given name                        | \<your-name>             | gator register "john go"                                       |
| login     | Logs into the account of the given name                 | \<your-name>             | gator login "john go"                                          |
| reset     | Resets your db                                          | None                     | gator reset                                                    |
| users     | Lists all users that have been registered               | None                     | gator users                                                    |
| addfeed   | Adds and follows a feed for the current logged in user  | \<feed-name> \<feed-url> | gator addfeed "Hacker News" "https://news.ycombinator.com/rss" |
| feeds     | Lists all feeds that have been added                    | None                     | gator feeds                                                    |
| follow    | Follow a feed that has already been added               | \<feed-url>              | gator follow "https://news.ycombinator.com/rss                 |
| following | Lists all feeds that the current logged in user follows | None                     | gator following                                                |
| unfollow  | Unfollow a feed                                         | \<feed-url>              | gator unfollow "https://news.ycombinator.com/rss               |
| browse    | Browse a list of posts from the feeds you follow        | OPTIONAL \<limit>        | gator browse 3                                                 |

Note: The browse command is the only command that has an optional argument. If
no argument has been specified it will use 2 by default.

## Issues

Found a bug? Report it using the issues tab. This was a pet project I used to
learn some Go. Chances are that this application is poorly written and if you
have any feedback, please open an issue! I will address it immediately.
