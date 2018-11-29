# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: 'Star Wars' }, { name: 'Lord of the Rings' }])
#   Character.create(name: 'Luke', movie: movies.first)

NUM_MESSAGES = 5
NUM_ARTICLES = 15


# 0.upto(NUM_MESSAGES) do |i|
#   Message.create content: Faker::HarryPotter.quote
# end

0.upto(NUM_ARTICLES) do |i|
  Article.create({
      headline: Faker::Book.title,
      lead: Faker::Lorem.sentence(9),
      pub_date: Faker::Date.backward(29),
      news_alert: Faker::Boolean.boolean(0.2),
      category: "news",
      author: Faker::Book.author,
      content: Faker::Lorem.paragraphs(15)
    }) 
end