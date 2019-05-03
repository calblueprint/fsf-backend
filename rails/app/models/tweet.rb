# == Schema Information
#
# Table tweets
#
# t.string "url"
# t.string "text"
# t.datetime "date"
# t.datetime "created_at", null: false
# t.datetime "updated_at", null: false
# t.bigint "message_id"
# t.index ["message_id"], name: "index_tweets_on_message_id"

class Tweet < ApplicationRecord
  # IMPLEMENT functions to create/edit message objects when tweets are supported
  # same as that in petition.rb/notice.rb/article.rb

  validates :date, presence:true
  validates :url, presence:true
  validates :text, presence:true
end
