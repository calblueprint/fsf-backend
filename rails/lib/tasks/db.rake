#This rake file updates the Article table with the most up to date FSF News Articles
#it accomplishes this by using Feedjira to parse the xml content from FSF RSS newsfeed
#it then iterates through each item from the FSF RSS feed and checks if it is in the table
#if its not in the table, the FSF News Article is added to Article table
namespace :db do
  desc "Updates FSF News Table by querying FSF News RSS Feed"
  task updateNews: :environment do
    url = "https://static.fsf.org/fsforg/rss/news.xml"
    Feedjira.configure do |config|
      config.parsers = [
         Feedjira::Parser::RSS
      ]
      config.strip_whitespace = true
      puts config.parsers
      puts config.logger
    end
    feed = Feedjira::Feed.fetch_and_parse url
    entries = feed.entries
    entries.each do |entry|
      articleInDB = Article.find_by title: entry.title
      puts "did we find an entry? here it is"
      puts entry
      article_params = {
          title: entry.title,
          link: entry.id, 
          pub_date: entry.published,
          content: entry.content,
          news_alert: false,
      }
      if(!articleInDB)
        begin
          savedArticle = Article.create(article_params)
          puts "we just created an article and"
        rescue ActiveRecord::StatementInvalid => invalid
          puts "#{article_params} invalid"
        end
        if savedArticle 
          article = Article.find(savedArticle.id)
          puts "success in creating article #{savedArticle}\n"
        else
          puts "failed to save article\n"
        end
        puts "we're saving!\n"
      else
        puts "This entry is already in the DB entry\n"
        puts articleInDB.title
      end
    end
    puts DateTime.now
  end
end