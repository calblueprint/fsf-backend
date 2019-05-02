# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: 'Star Wars' }, { name: 'Lord of the Rings' }])
#   Character.create(name: 'Luke', movie: movies.first)

NUM_ITEMS = 5 

# #Creates 5 Test Messages
# 0.upto(NUM_ITEMS) do |i|
#   Message.create({
#       title: "#{Faker::Kpop.ii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups}",
#       content: Faker::HarryPotter.quote,
#       link: "fsf://fsf/donate",
#     })
# end

# Create source urls
Source.create source_type: :rss, rss_url: "https://static.fsf.org/fsforg/rss/news.xml"
Petition.create title: "Please Act Now!", description: "Your digital rights are dying", link: "https://my.fsf.org/civicrm/petition/sign?sid=8&reset=1"

#Creates 5 Test Articles
# t.string "gs_user_id"
# t.string "gs_user_name"
# t.datetime "published"
# t.string "content_text"
# t.string "content_html"
# t.string "url"
# t.datetime "created_at", null: false
# t.datetime "updated_at", null: false
# t.string "gs_user_handle"
# t.boolean "news_alert", default: false
# t.bigint "message_id"
# t.index ["message_id"], name: "index_notices_on_message_id"
0.upto(NUM_ITEMS) do |i|
    Article.create({
      title: "this is my title",
      link: Faker::Internet.url("fsf.org"),
      pub_date: Faker::Date.backward(29),
      news_alert: false,
      content: Faker::Lorem.paragraphs(15),
      description: Faker::Lorem.sentence(88),
      summary: Faker::Lorem.sentence(88)
    })
end

#Creates 5 Test Notices
# t.string "gs_user_id"
# t.string "gs_user_name"
# t.datetime "published"
# t.string "content_text"
# t.string "content_html"
# t.string "url"
# t.datetime "created_at", null: false
# t.datetime "updated_at", null: false
# t.string "gs_user_handle"
# t.boolean "news_alert", default: false
# t.bigint "message_id"
# t.index ["message_id"], name: "index_notices_on_message_id"
# validates :id, presence:true
# validates :gs_user_id, presence:true
# validates :gs_user_name, presence:true
# validates :gs_user_handle, presence:true
# validates :published, presence:true
# validates :content_text, presence:true
# validates :content_html, presence:true
# validates :url, presence:true
0.upto(NUM_ITEMS) do |i|
  Notice.create({
    gs_user_id: Faker::Kpop.ii_groups,
    gs_user_name: Faker::Kpop.ii_groups,
    gs_user_handle: "555",
    published: Faker::Date.backward(29),
    content_text: Faker::Lorem.paragraphs(15),
    content_html: Faker::Lorem.paragraphs(15),
    url: Faker::Internet.url("fsf.org"),
    news_alert: false
  })
end

#Creates 5 Test Tweets
# t.string "url"
# t.string "text"
# t.datetime "date"
# t.datetime "created_at", null: false
# t.datetime "updated_at", null: false
# t.bigint "message_id"
# t.index ["message_id"], name: "index_tweets_on_message_id"
0.upto(NUM_ITEMS) do |i|
  Tweet.create({
    url: "#{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")} #{Faker::Internet.url("fsf.org")}",
    text: Faker::Kpop.ii_groups,
    date: Faker::Date.backward(29),
  })
end