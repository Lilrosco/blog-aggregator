# Gator
Boot.dev project building a blog aggregator using Go and Postgres

A CLI tool that enables end-users to do the following:

- Register users
- Add specified RSS feeds across the web to be tracked
- Save tracked feeds in a PostgresSQL database
- Run a `agg` process that will fetch followed feeds on an interval
- View posts that were aggregated in the terminal with links to full post

## Prerequistes

You will need to have the following tools installed on your local machine

- **Golang (Go)** 1.26 or newer - <https://go.dev/doc/install>
- **PostgresSQL** (Any latest version) - <https://www.postgresql.org/download/>

Setup your system user password and database user password

- **goose** (for running database migrations though you can use what works best for you) - you can install using:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

To run goose migrations use the follow commands

Using your own DB URL connection

rollback

```sh
goose postgres "postgres://<USERNAME>:<PASSWORD>@localhost:5432/gator" down
```

apply migrations

```sh
goose postgres "postgres://<USERNAME>:<PASSWORD>@localhost:5432/gator" up
```

- **sqlc** (if you plan to add to change SQL queries) - <https://sqlc.dev> 

you can install using:

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

The Go code this command generates will be located in the `internal/database/` directory. To use run the command

```sh
sqlc generate
```

If you modify anything in `sql/queries` be sure to run the command above to generate the needed Go code.

## Installation

Install the latest version of this repo

```sh
go install github.com/Lilrosco/gator@latest
```

## Config

Create a `.gatorconfig.json` file in your home directory with the following content:

```json
{
  "db_url": "postgres://<USERNAME>:<PASSWORD>@localhost:5432/database?sslmode=disable",
  "current_user_name": ""
}
```

Be sure to add your DB's username and password

## Usage

Register a user and add a RSS feed

```sh
gator register testuser
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

Other commands

```sh
gator login testuser
gator users
gator feeds
gator follow "https://example.com/feed.xml"
gator following
gator unfollow "https://example.com/feed.xml"
gator agg 30s
gator browse
gator browse 5
gator help
```

`agg` is a long-running process that fetches one feed on the interval specified (example - 30s, 1m, 15m). Leave it running in a separate terminal and stop it with `Ctrl / Cmd + C`. The optional number passed to `browse` limits the number of displayed posts; it defaults to 2.

For more help try running the command

```sh
gator help
```

for a list of available commands and their usage.
