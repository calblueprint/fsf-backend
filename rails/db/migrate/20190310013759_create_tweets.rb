class CreateTweets < ActiveRecord::Migration[5.2]
  def change
    create_table :tweets do |t|
      t.string :url
      t.string :text
      t.datetime :date

      t.timestamps
    end
  end
end
