class CreateSources < ActiveRecord::Migration[5.2]
  def change
    create_table :sources do |t|
      t.integer :source_type
      t.string :url

      t.timestamps
    end
  end
end
