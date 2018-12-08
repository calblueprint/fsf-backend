class Article < ApplicationRecord
  validates :title, presence: true
  validates :link, presence: true
  validates :pub_date, presence: true
  validates :content, presence: true
  validates :news_alert, inclusion: { in: [true, false] }
  validates :summary, presence: true
end
