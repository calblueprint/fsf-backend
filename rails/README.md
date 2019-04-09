# Rails backend for FSF's mobile app

This is the rails backend for FSF's mobile app. It serves the news feed tab by storing articles from RSS and posts from Twitter or GNU social.

### Environment

Ruby 2.5.1
Rails 5.2.0

- Ruby Linter
  - command: ./lint
  - https://github.com/prettier/plugin-ruby

Default database connection uses PostgreSQL connection. Change `config/database.yml` to use other db connections.

### Install dependencies and run on your native environment

Prerequisite: proper ruby version

First install rails and other dependencies.

```
bundle install
```

Then creates the database if necessary.

```
rails db:create
```

To set up the database, run:

```
rails db:migrate
```

To seed the database, use:

```
rails db:seed
```

To start the server, run:

```
rails s
```

### Setup in Docker

`docker-compose` can be used to setup the ruby / rails environment and postgres server in docker containers. You need to install docker and docker-compose before you execute the scripts.

The following commands set up the containers:

```
cp config/database.yml.docker config/database.yml
./docker_setup.sh
```

You can use `docker-compose exec` commands to run the rails server. See docker-compose documentation for more details.
