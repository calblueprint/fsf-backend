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
  belongs_to :message, optional: true, dependent: :destroy
  after_create_commit :create_message_if_new_alert
  before_save :update_message_object

  private
  def create_message_if_new_alert
    if self.news_alert
      new_message = Message.create(content: self.description, title: self.title, link: "fsf://fsf/action/" + self.id.to_s)
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
        elsif self.message.content != self.description
          m = self.message
          m.content = self.description
          m.save
        end
      else
        new_message = Message.create(content: self.description, title: self.title, link: "fsf://fsf/action/" + self.id.to_s)
        self.message = new_message
      end
    end
  end

  validates :title, presence: true
  validates :description, presence: true
  validates :link, presence: true
end
