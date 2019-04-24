# == Schema Information
#
#  t.string "title"
#  t.string "link"
#  t.datetime "pub_date"
#  t.text "content"
#  t.boolean "news_alert", default: false
#  t.datetime "created_at", null: false
#  t.datetime "updated_at", null: false
#  t.text "description"
#  t.text "summary"
class Article < ApplicationRecord
  has_one :message
  after_create_commit :create_message_if_new_alert
  after_update_commit :update_message_object
  after_destroy :destroy_message_object
    
  # TODO: need to figure out the linking stuff for deep linking

  private
  def create_message_if_new_alert
    if self.news_alert
      Message.create(content: self.content, title: self.title, link: "fsf://fsf/profile", article_id: self.id)
    end
  end

  private
  def update_message_object
    if self.news_alert
      Message.create(content: self.content, title: self.title, link: "fsf://fsf/profile", article_id: self.id)
    else
      Message.where(article_id: self.id).destroy_all
    end
  end

  private
  def destroy_message_object
    if self.news_alert
      Message.where(article_id: self.id).destroy_all
    end
  end

  validates :title, presence: true
  validates :link, presence: true
  validates :pub_date, presence: true
  validates :content, presence: true
  validates :news_alert, inclusion: { in: [true, false] }
  validates :summary, presence: true
end


