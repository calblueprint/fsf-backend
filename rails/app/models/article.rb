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
  # after_update_commit :set_created_at
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
    if self.news_alert
      if self.message
        # check if message contents match the new message we would like to have here, if not, update the message contents
        if self.message.title != self.title
          m = self.message
          m.title = self.title
          m.save
        elsif self.message.content != self.content
          m = self.message
          m.content = self.content
          m.save
        end
      else
        new_message = Message.create(content: self.content, title: self.title, link: "fsf://fsf/news/article/" + self.id.to_s)
        self.message = new_message
      end
    end
  end

  # private
  # def set_created_at
  #   if self.message
  #     self.message.created_at = self.pub_date
  #   end
  # end

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


