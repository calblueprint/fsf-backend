class Petition < ApplicationRecord
  # db schema:
  # create_table "petitions", force: :cascade do |t|
  #   t.string "title"
  #   t.string "description"
  #   t.string "link"
  #   t.datetime "created_at", null: false
  #   t.datetime "updated_at", null: false

  validates :title, presence: true
  validates :description, presence: true
  validates :link, presence: true
end
