class CreateArticles < ActiveRecord::Migration[5.2]
  def change
    create_table :articles do |t|
      t.string :title
      t.string :link
      t.datetime :pub_date
      t.text :content
      t.boolean :news_alert, default: false

      t.timestamps
    end
  end
end
