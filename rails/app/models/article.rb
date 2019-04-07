class Article < ApplicationRecord
  #t.string "title"
  #  t.string "link"
  #  t.datetime "pub_date"
  #  t.text "content"
  #  t.boolean "news_alert", default: false
  #  t.datetime "created_at", null: false
  #  t.datetime "updated_at", null: false
  #  t.text "description"
  #  t.text "summary"

  validates :title, presence: true
  validates :link, presence: true
  validates :pub_date, presence: true
  validates :content, presence: true
  validates :news_alert, inclusion: { in: [true, false] }
  validates :summary, presence: true
end
