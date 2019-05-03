# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2019_04_24_075947) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "articles", force: :cascade do |t|
    t.string "title"
    t.string "link"
    t.datetime "pub_date"
    t.text "content"
    t.boolean "news_alert", default: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.text "description"
    t.text "summary"
    t.bigint "message_id"
    t.index ["message_id"], name: "index_articles_on_message_id"
  end

  create_table "messages", force: :cascade do |t|
    t.text "content"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "title"
    t.string "link"
  end

  create_table "notices", force: :cascade do |t|
    t.string "gs_user_id"
    t.string "gs_user_name"
    t.datetime "published"
    t.string "content_text"
    t.string "content_html"
    t.string "url"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "gs_user_handle"
    t.boolean "news_alert", default: false
    t.bigint "message_id"
    t.index ["message_id"], name: "index_notices_on_message_id"
  end

  create_table "petitions", force: :cascade do |t|
    t.string "title"
    t.string "description"
    t.string "link"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.boolean "news_alert", default: false
    t.bigint "message_id"
    t.index ["message_id"], name: "index_petitions_on_message_id"
  end

  create_table "sources", force: :cascade do |t|
    t.integer "source_type"
    t.string "rss_url"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "twitter_consumer_key"
    t.string "twitter_consumer_secret"
    t.string "twitter_access_token"
    t.string "twitter_access_token_secret"
    t.string "twitter_username"
    t.string "GNU_social_url"
  end

  create_table "tweets", force: :cascade do |t|
    t.string "url"
    t.string "text"
    t.datetime "date"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "message_id"
    t.index ["message_id"], name: "index_tweets_on_message_id"
  end

end
