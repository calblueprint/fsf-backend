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
  belongs_to :message, optional: true
  after_create_commit :create_message_if_new_alert
  before_save :update_message_object
  after_destroy :destroy_message_object

  private
  def create_message_if_new_alert
    if self.news_alert
      new_message = Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s)
      self.message = new_message
      self.save
    end
  end

  private
  def update_message_object
    # Message.where(article_id: self.id).destroy_all
    if self.news_alert
      # TODO: multiple messages are created when a message is updated with the news_alert field checked
      new_message = Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s)
      self.message = new_message
      # Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s, article_id: self.id)  
    else 
      if self.message
        self.message.destroy
      end
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
  validates :link, presence: true
  validates :pub_date, presence: true
  validates :content, presence: true
  validates :news_alert, inclusion: { in: [true, false] }
  validates :summary, presence: true
end


