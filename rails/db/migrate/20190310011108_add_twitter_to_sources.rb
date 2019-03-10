class AddTwitterToSources < ActiveRecord::Migration[5.2]
  def change
    change_table :sources do |t|
      t.rename :url, :rss_url
      t.string :twitter_consumer_key
      t.string :twitter_consumer_secret
      t.string :twitter_access_token
      t.string :twitter_access_token_secret
      t.string :twitter_username
    end
  end
end
