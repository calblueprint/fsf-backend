class AddGnUsocialUrlToSources < ActiveRecord::Migration[5.2]
  def change
    add_column :sources, :GNU_social_url, :string
  end
end
