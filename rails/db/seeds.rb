# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: 'Star Wars' }, { name: 'Lord of the Rings' }])
#   Character.create(name: 'Luke', movie: movies.first)

#NUM_MESSAGES = 5
NUM_ARTICLES = 5 

# 0.upto(NUM_MESSAGES) do |i|
#   Message.create content: Faker::HarryPotter.quote
# end

Petition.create title: "Please Act Now!", description: "Your digital rights are dying", link: "https://my.fsf.org/civicrm/petition/sign?sid=8&reset=1"

#Creates 5 Test Articles
0.upto(NUM_ARTICLES) do |i|
    Article.create({
      title: "#{Faker::Kpop.ii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups} #{Faker::Kpop.iii_groups}",
      link: Faker::Internet.url("fsf.org"),
      pub_date: Faker::Date.backward(29),
      news_alert: false,
      content: Faker::Lorem.paragraphs(15),
      description: Faker::Lorem.sentence(88),
      summary: Faker::Lorem.sentence(88)
    })
end
