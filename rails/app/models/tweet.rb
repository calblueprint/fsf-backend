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
  # NOT IMPLEMENTED - FUTURE EXTENSION
  # belongs_to :message, optional: true
  # after_create_commit :create_message_if_new_alert
  # before_save :update_message_object
  # after_destroy :destroy_message_object

  # private
  # def create_message_if_new_alert
  #   if self.news_alert
  #     new_message = Message.create(content: self.text, title: self.date, link: "fsf://fsf/TWEETS_NOT_USED/" + self.id.to_s)
  #     self.message = new_message
  #     self.save
  #   end
  # end

  # private
  # def update_message_object
  #   # Message.where(article_id: self.id).destroy_all
  #   if self.news_alert
  #     # TODO: multiple messages are created when a message is updated with the news_alert field checked
  #     new_message = Message.create(content: self.text, title: self.date, link: "fsf://fsf/TWEETS_NOT_USED/" + self.id.to_s)
  #     self.message = new_message
  #     # Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s, article_id: self.id)  
  #   end
  # end

  # private
  # def destroy_message_object
  #   if self.news_alert
  #     self.message.destroy
  #     # Message.where(article_id: self.id).destroy_all
  #   end
  # end
    
  # validates :id, presence:true
  validates :date, presence:true
  validates :url, presence:true
  validates :text, presence:true
end
