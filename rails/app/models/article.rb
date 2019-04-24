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
  after_create_commit :create_message_object
  after_update_commit :create_message_object
  after_destroy :destroy_message_object
    
  private
  def create_message_object
      Message.create(content: self.content, title: "I made this BABY 2", link: "fsf://fsf/profile")
  end

  private
  def destroy_message_object

  end

  validates :title, presence: true
  validates :link, presence: true
  validates :pub_date, presence: true
  validates :content, presence: true
  validates :news_alert, inclusion: { in: [true, false] }
  validates :summary, presence: true
end


