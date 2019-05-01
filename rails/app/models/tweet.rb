# == Schema Information
#
# Table tweets
#
# id                  :integer
# date                :datetime
# url                 :string
# text                :string
class Tweet < ApplicationRecord
  belongs_to :message, optional: true
  after_update_commit :create_message_object
    
  private
  def create_message_object
      Message.create(content: Faker::HarryPotter.quote, title: "I made this BABY 2", link: "fsf://fsf/profile")
  end
  validates :id, presence:true
  validates :date, presence:true
  validates :url, presence:true
  validates :text, presence:true
end
