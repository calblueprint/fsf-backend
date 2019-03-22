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
        parse_GNUsocial(source.GNU_social_url, 2)
      else
        puts "Unknown source type #{source.source_type}\n"
      end
    end

    puts DateTime.now
  end
end

def parse_GNUsocial(source, limit)
  puts "This is the GNU social user_timeline we are sourcing #{source}"
  puts "This will obtain #{limit} page(s) of results from the GNU social time line"
  response = RestClient.get source, { accept: :json }
  responseBody = JSON.parse(response.body)
  links = responseBody['links']
  notices = responseBody['items']
  parse_GNUsocial_pages(notices, links, limit)
end

def parse_GNUsocial_pages(notices, links, limit)
  limit_num = limit
  while (links.length > 1) && (links[1]['rel']['rel'] == 'next') && (limit_num > 0)
    notices.each do |notice|
      if (notice['object']['objectType'] == 'note')
        parse_notice(notice)
        puts "#{notice['url']}"
      else
        puts "\n#{notice['url']} is not a note but a #{notice['object']['objectType']} object\n"
      end
    end
    puts links[1]['url']
    puts 'next url'
    response = RestClient.get links[1]['url'], { accept: :json }
    puts "status code #{response.code}"
    responseBody = JSON.parse(response.body)
    links = responseBody['links']
    notices = responseBody['items']
    limit_num = limit_num - 1
    puts "There are #{limit_num} requests left"
    puts JSON.pretty_generate(links)
  end
end

def parse_notice(notice)
  notice_id = notice['object']['status_net']['notice_id']
  notice_detail_url = "https://status.fsf.org/api/statuses/show/#{notice_id}.json"
  response = RestClient.get(notice_detail_url) { |response, request, result| response }
  if response.code == 200
    notice_details = JSON.parse(response.body)
    gs_user_id = notice['actor']['status_net']['profile_info']['local_id']
    gs_user_name = notice['actor']['displayName']
    published = notice['published']
    content_text = notice_details['text']
    content_html = notice_details['statusnet_html']
    url = notice['url']
    unless Notice.exists?(notice_id)
      Notice.create(
        {
          id: notice_id,
          gs_user_id: gs_user_id,
          gs_user_name: gs_user_name,
          published: published,
          content_text: content_text,
          content_html: content_html,
          url: url
        }
      )
    else
      puts "Existing Notice #{notice_id}"
    end
  else
    puts "This is the url #{notice_detail_url}"
    puts "Notice #{notice_id} details do not exist"
    puts "the following url failed https://status.fsf.org/api/statuses/show/#{notice_id}.json"
    puts "Response code was #{response.code}\n"
  end
end

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
    Tweet.create({ id: tweet.id, date: tweet.created_at, url: tweet.uri, text: tweet.full_text })
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
    summary_from_content = summary_parser[2].split.join(' ') if (not summary_parser[2].empty?)
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
