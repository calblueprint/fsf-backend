# == Schema Information
#
# Table notices
#
# id                  :string
# gs_user_id                  :string
# gs_user_name                  :string
# published                  :datetime
# content_text                  :string
# content_html                  :string
# url                 :string
class Notice < ApplicationRecord
  belongs_to :message, optional: true
  after_create_commit :create_message_if_new_alert
  before_save :update_message_object
  after_destroy :destroy_message_object

  private
  def create_message_if_new_alert
    if self.news_alert
      new_message = Message.create(content: self.content_text, title: self.gs_user_name, link: "fsf://fsf/gnu/social/" + self.id.to_s)
      self.message = new_message
      self.save
    end
  end

  private
  def update_message_object
    # Message.where(article_id: self.id).destroy_all
    if self.news_alert
      # TODO: multiple messages are created when a message is updated with the news_alert field checked
      new_message = Message.create(content: self.content_text, title: self.gs_user_name, link: "fsf://fsf/gnu/social/" + self.id.to_s)
      self.message = new_message
      # Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s, article_id: self.id)  
    else 
      if self.message
        self.message.destroy
      else
    end
  end

  private
  def destroy_message_object
    if self.news_alert
      self.message.destroy
      # Message.where(article_id: self.id).destroy_all
    end
  end

  # validates :id, presence:true
  validates :gs_user_id, presence:true
  validates :gs_user_name, presence:true
  validates :gs_user_handle, presence:true
  validates :published, presence:true
  validates :content_text, presence:true
  validates :content_html, presence:true
  validates :url, presence:true
end
