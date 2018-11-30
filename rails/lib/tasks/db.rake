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
      puts entry
      article_params = {
          title: entry.title, 
          link: entry.id, 
          pub_date: entry.published,
          content: entry.content,
          news_alert: :false,
      }
      if(!articleInDB)
        article = Article.new(article_params)
        begin
          saved = article.save!
        rescue ActiveRecord::StatementInvalid => invalid
          puts "invalid"
        end
        if saved
          article = Article.find(article.id)
          puts "success in creating article #{article}\n"
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