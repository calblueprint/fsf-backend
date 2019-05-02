# == Schema Information
#
# Table tweets
#
# id                  :integer
# date                :datetime
# url                 :string
# text                :string
class Tweet < ApplicationRecord
  validates :date, presence:true
  validates :url, presence:true
  validates :text, presence:true
end
