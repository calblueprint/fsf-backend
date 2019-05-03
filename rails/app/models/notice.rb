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
  belongs_to :message, optional: true, dependent: :destroy
  after_create_commit :create_message_if_new_alert
  before_save :update_message_object

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
    if self.news_alert
      if self.message
        # check if message contents match the new message we would like to have here, if not, update the message contents
        if self.message.title != self.gs_user_name
          m = self.message
          m.title = self.gs_user_name
          m.save
        elsif self.message.content != self.content_text
          m = self.message
          m.content = self.content_text
          m.save
        end
      else
        new_message = Message.create(content: self.content_text, title: self.gs_user_name, link: "fsf://fsf/gnu/social/" + self.id.to_s)
        self.message = new_message
      end
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
