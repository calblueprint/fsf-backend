class CreateArticles < ActiveRecord::Migration[5.2]
  def change
    create_table :articles do |t|
      t.string :headline
      t.text :lead
      t.datetime :pub_date
      t.boolean :news_alert
      t.string :category
      t.string :author
      t.text :content

      t.timestamps
    end
  end
end
