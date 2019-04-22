class AddBoolean < ActiveRecord::Migration[5.2]
  def change
    add_column :petitions, :news_alert, :boolean, :default => false
    add_column :notices, :news_alert, :boolean, :default => false
  end
end
