# == Schema Information
#
# Table notices
#
# db schema:
# create_table "petitions", force: :cascade do |t|
#   t.string "title"
#   t.string "description"
#   t.string "link"
#   t.datetime "created_at", null: false
#   t.datetime "updated_at", null: false
# url                 :string
class Petition < ApplicationRecord
  belongs_to :message, optional: true
  after_create_commit :create_message_if_new_alert
  before_save :update_message_object
  after_destroy :destroy_message_object

  private
  def create_message_if_new_alert
    if self.news_alert
      new_message = Message.create(content: self.description, title: self.title, link: "fsf://fsf/news/CHANGE_THIS_HERE/" + self.id.to_s)
      self.message = new_message
      self.save
    end
  end

  private
  def update_message_object
    # Message.where(article_id: self.id).destroy_all
    if self.news_alert
      # TODO: multiple messages are created when a message is updated with the news_alert field checked
      new_message = Message.create(content: self.description, title: self.title, link: "fsf://fsf/news/CHANGE_THIS_HERE/" + self.id.to_s)
      self.message = new_message
      # Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s, article_id: self.id)  
    end
  end

  private
  def destroy_message_object
    if self.news_alert
      self.message.destroy
      # Message.where(article_id: self.id).destroy_all
    end
  end

  validates :title, presence: true
  validates :description, presence: true
  validates :link, presence: true
end
