# Gator RSS Aggregator

## What is Gator

**Gator** is a RSS feed Aggregator for the terminal to read blog posts in the terminal

## Dependencies

**Gator** requires both [goose](https://github.com/pressly/goose) and [postgres](https://www.postgresql.org/)

### To install goose:

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### To install postgres:

```
apt install postgresql
```

## Setup Postgres

### Ensure that postgres is installed correctly:

```
psql --version
```

### For Linux systems, update the password for postgres:

```
sudo passwd postgres
```

### Start server for postgres:

- Mac: `brew services start postgresql`
- Linux: `sudo service postgresql start`

### Enter psql shell:

- Mac: `psql postgres`
- Linux `sudo -u postgres psql`

### Create a database called `gator`:

```
CREATE DATABASE gator;
```

### Connect to the new database:

```
\c gator
```

### For Linux systems only setup the user password:

```
ALTER USER postgres PASSWORD 'postgres';

```

Exit with `exit`

## Setup goose

### Ensure goose is installed properly:

```
goose --version
```

### Enter the sql/schema folder in the repo:

```
cd sql/schema
```

### Migrate up with goose:

Your postgres connection is dependent on how you setup the postgres server

- Mac: `postgres://<username>:@localhost:5432/gator`
- Linux: `postgres://postgres:postgres@localhost/gator`

### Migrate:

```
goose postgres <connection-string> up
```

### Create the config file:

Create a config file in your home directory with

```
{
  "db_url": <connection-string>
}
```

## How to use

### User Controls

Create a user with the `register` command

```
gatorcli register <name>
```

List all registered users with the `users` command

```
gatorcli users
```

Switch to another registered user with the `login` command

```
gatorcli login <name>
```

### Feed Controls

Add a feed by using the `addfeed` command, the current user with automatically subscribed to the feed

```
gatorcli addfeed <feed-name> <feed-url>
```

List all the added feeds with the `feeds` command

```
gatorcli feeds
```

### Following Controls

Follow an added feed with the `follow` command for the current user

```
gatorcli follow <feed-url>
```

Unfollow a feed with the `follow` command for the current user

```
gatorcli unfollow <feed-url>
```

List the feeds the current user is following with the `following` command

```
gatorcli following
```

### Aggregator

Aggregate feeds with the `agg` command by a time interval

The time interval is formated like `eg: 1s, 1m, 1h`

```
gatorcli agg <time-interval
```

### Browse

Browse posts of followed feeds with the `browse` command

```
gatorcli browse
```

### Help

Get help of all commands with `help`

```
gatorcli help
```
