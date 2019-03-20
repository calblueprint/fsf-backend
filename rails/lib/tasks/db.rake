#  This rake file updates the Article table with the most up to date FSF News Articles
#  it accomplishes this by using Feedjira to parse the xml content from FSF RSS newsfeed
#  it then iterates through each item from the FSF RSS feed and checks if it is in the table
#  if its not in the table, the FSF News Article is added to Article table through the
#  parsing the feed item with feedjira and nokogiri

namespace :db do
  desc 'Updates FSF News Table by querying FSF News RSS Feed'
  task updateNews: :environment do
    Source.all.each do |source|
      if source.rss?
        parse_rss(source.rss_url)
      elsif source.twitter?
        parse_twitter(source)
      elsif source.GNUsocial?
        parse_GNUsocial(source)
      else
        puts "Unknown source type #{source.source_type}\n"
      end
    end

    puts DateTime.now
  end
end

def parse_GNUsocial(source, num)
  response = HTTParty.get(source)
  puts response
end

# def parse_notice(notice)
# end
# def parse_status(status)
# end

def parse_twitter(source)
  client = source.get_twitter_client

  # for now, just get the latest 20 tweets
  # TODO: long term solution
  tweets = client.user_timeline(source.twitter_username, count: 20)
  tweets.each do |tweet|
    parse_tweet(tweet)
    puts "#{tweet.id}"
  end
end

def parse_tweet(tweet)
  unless Tweet.exists?(tweet.id)
    Tweet.create(
      {
        id: tweet.id,
        date: tweet.created_at,
        url: tweet.uri,
        text: tweet.full_text
      }
    )
  else
    puts "Existing tweet #{tweet.id}"
  end
end

def parse_rss(url)
  Feedjira.configure do |config|
    config.parsers = [Feedjira::Parser::RSS]
    config.strip_whitespace = true
  end

  feed = Feedjira::Feed.fetch_and_parse url

  entries = feed.entries
  entries.each { |entry| parse_rss_entry(entry) }
end

def parse_rss_entry(entry)
  articleInDB = Article.find_by title: entry.title
  if (!articleInDB)
    html_doc = Nokogiri.HTML(entry.content)
    paragraph_list = html_doc.css('p')
    paragraph_item = paragraph_list.first.text

    #  This regex splits the incoming paragraph by periods
    sentence_list = paragraph_item.split(/(?<=(?<=[a-zA-Z])[a-zA-Z])\./)
    sentence_list.delete_if { |el| el.length == 0 }
    summary_from_content = sentence_list.first
    summary_format = summary_from_content.split.join(' ')
    summary_from_content = summary_format if summary_format
    #  This regex specifically partitions a sentence to exclude the first
    #  sentence of the content paragraph which has the general form of
    #  BOSTON, Massachusetts, USA -- Thursday, October 18, 2018 --
    #  This will be included until the description portion of the
    #  RSS newsfeed has a deterministic non-null value
    summary_parser = summary_from_content.partition(/\, \d{4} (?:[^A-Za-z]{2})/)
    if (not summary_parser[2].empty?)
      summary_from_content = summary_parser[2].split.join(' ')
    end
    article_params = {
      title: entry.title,
      link: entry.id,
      pub_date: entry.published,
      content: entry.content,
      news_alert: false,
      description: entry.summary,
      summary: summary_from_content
    }
    begin
      savedArticle = Article.create(article_params)
    rescue ActiveRecord::StatementInvalid => invalid
      puts "#{article_params} invalid"
    end
    if savedArticle
      article = Article.find(savedArticle.id)
      puts "success in creating article #{article}\n"
    else
      puts "failed to save article\n"
    end
  else
    puts "This entry is already in the DB entry\n"
    puts articleInDB.title
  end
end
