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
  validates :id, presence:true
  validates :gs_user_id, presence:true
  validates :gs_user_name, presence:true
  validates :published, presence:true
  validates :content_text, presence:true
  validates :content_html, presence:true
  validates :url, presence:true
end
