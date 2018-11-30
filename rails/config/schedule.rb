# Use this file to easily define all of your cron jobs.
#
# It's helpful, but not entirely necessary to understand cron before proceeding.
# http://en.wikipedia.org/wiki/Cron

####################### FSF NOTES
# the following command
# whenever --update-crontab
# updates the cron job
# running <whenever> in your command prompt will output the command it will run
# this cron job set by the whenever will run the following rake command that updates the
# news database 
every 1.day do # 1.minute 1.day 1.week 1.month 1.year is also supported
    rake "db:updateNews"
end

# Example
### For Development Purposes
# every 1.minute do # 1.minute 1.day 1.week 1.month 1.year is also supported
#     command("/bin/bash -l -c 'cd /Users/franco/Desktop/sandbox/fsf/fsf-rails/rails && RAILS_ENV=development bundle exec rake db:updateNews --silent'")
# end


#set :output, "/path/to/my/cron_log.log"

# every 2.hours do
#   commad "/usr/bin/some_great_command"
#   runner "MyModel.some_method"
#   rake "some:great:rake:task"
# end
#
# every 4.days do
#   runner "AnotherModel.prune_old_records"
# end

# Learn more: http://github.com/javan/whenever
