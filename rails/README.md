# Rails backend for FSF's mobile app

Welcome to the rails backend for FSF's mobile application. It serves the news feed tab by storing articles from RSS and posts from Twitter or GNU social.

It serves three primary purposes

1. As an API access point for serving up news content for the News Feed Tab on the Mobile application
2. As an ETL pipeline (Extract, Transform and Load) for taking in data sources (FSF News RSS Feed, FSF GNU Social, Twitter) and transforming it into a format that can be used and keeping all data sources up to data
3. As an admin dashboard - allowing users who have administrator access to add/modify/delete new news sources.

---

## 💻️ Setting Up Your Environment Locally

In order to run the server locally you will need the following installed in your development environment

- Ruby 2.5.1
- Rails 5.2.0
- PostgresSQL 10.5

### 🔗 ️Installing Dependencies, Creating the Database and Running Migrations Locally

The following command will install the required gems for the application

```
$ bundle install
```

The following command will create the database locally

```
rails db:create
```

To set up the database run:

```
rails db:migrate
```

To start the server, run:

```
rails s
```

### 🔢 Seeding the Database with News Sources Locally

After running the server go to http://localhost:3000/admin to access the admin dashboard

### 🗞Adding News Sources

The following are instructions to adding different Sources to the database

Source types that are currently supported

- rss
- GNU Social Timeline
- Twitter Feed

#### 📰Adding FSF RSS News Feed Source

- Once on the admin dashboard navigate to the lefthand navigation panel and click on 'Sources'
- You will see a tab called "+ Add New" -- click on it and it will bring you to a form
- On the "Source type" dropdown -- select rss
- For the Rss url - specify the FSF rss news feed by providing the following RSS endpoint
  - https://static.fsf.org/fsforg/rss/news.xml
- Leave all fields blank
- Click "Save" at the bottom of the form

#### 📰Adding FSF GNU Social Timeline Source

- Once on the admin dashboard navigate to the lefthand navigation panel and click on 'Sources'
- You will see a tab called "+ Add New" -- click on it and it will bring you to a form
- On the "Source type" dropdown -- select GNUsocial
- For the FSF GNU Social url endpoint - specify the following endpoint in the "Gnu social url" field
  - https://status.fsf.org/api/statuses/user_timeline/440.as
- Leave all fields blank
- Click "Save" at the bottom of the form

#### 📑Modifying How Many Pages The GNU Social Parser will Parse

File to modify: rake.db

- Initially you will see parse_GNUsocial(source.GNU_social_url, 2), the second parameter represents how many page results we would like to parse -- each page results is 20 notices -- you may change 2 to 200.

#### 🔢Seeding Database Locally

Once you have added the Sources (RSS endpoint, GNU social endpoint) run the following command

```
rake db:updateNews
```

#### 📶API Access Points

Once you have your database setup, have the rails server running and have seeded the database with the sources specified in the previous step, if you go to the following routes in your browser, you will be able to access notices and articles endpoints

- http://localhost:3000/api/v1/notices
- http://localhost:3000/api/v1/articles

#### 📶Setting up a Cron Job Locally to Update Database with New Items from Sources

For automating the script that runs the job to update the database with the latest news items from the sources specified, we use the ruby gem 'whenever'
- https://github.com/javan/whenever

File to modify: config/schedule.rb

The following is the command that is included in the schedule.rb file -- the following command is for testing the whenever gem and cron capabilities in your local environment

- make sure to modify the path in the command to specify the abosolute path to the rails folder

```
every 1.day do # 1.minute 1.day 1.week 1.month 1.year is also supported
    command("/bin/bash -l -c 'cd /your/absolute/path/to/fsf-rails/rails/ && RAILS_ENV=development bundle exec rake db:updateNews --silent'")
end
```
- After modifying the file run the following command

```
whenever --update-crontab
```
- Run the following command to see the cron job

```
whenever --update-crontab
```
You should see something similar to
```
0 0 * * * /bin/bash -l -c '/bin/bash -l -c '\''cd /your/absolute/path/to/fsf-rails/rails/ && RAILS_ENV=development bundle exec rake db:updateNews --silent'\'''
```

### 🔨 Utils

#### 🔬Ruby Linter
The rails folder contains a ruby linter using prettier as a means to specify formatting for code

Files:
- lint.sh
  - Contains the linting command along with the list of files to apply formatting to
- .prettierrc
  - This file contains configurations for how to format ruby code

System Requirements:
- Install node 8.12
- Install npm 6.5.0

Installing Ruby Linter
```
npm i
```
Running the Ruby Linter
  - command: ./lint
  - https://github.com/prettier/plugin-ruby



#### 🔨 Configuration

Default database connection uses PostgreSQL connection. Change `config/database.yml` to use other db connections.

### Setup in Docker

`docker-compose` can be used to setup the ruby / rails environment and postgres server in docker containers. You need to install docker and docker-compose before you execute the scripts.

The following commands set up the containers:

```
cp config/database.yml.docker config/database.yml
./docker_setup.sh
```

You can use `docker-compose exec` commands to run the rails server. See docker-compose documentation for more details.
