# == Schema Information
#
# Table name: sources
#
# id              :integer
# source_type     :integer (enum)
# rss_url         :string  "url is for RSS feed"
# twitter_username               :string
# twitter_consumer_key           :string
# twitter_consumer_secret        :string
# twitter_access_token           :string
# twitter_access_token_secret    :string
class Source < ApplicationRecord
  enum source_type: %i[rss twitter GNUsocial]

  validates :source_type, presence:true

  def get_twitter_client
    require 'twitter'
    client = Twitter::REST::Client.new do |config|
      config.consumer_key        = twitter_consumer_key
      config.consumer_secret     = twitter_consumer_secret
      config.access_token        = twitter_access_token
      config.access_token_secret = twitter_access_token_secret
    end
  end

end
