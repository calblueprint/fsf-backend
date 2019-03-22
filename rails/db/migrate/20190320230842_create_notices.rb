class CreateNotices < ActiveRecord::Migration[5.2]
  def change
    create_table :notices do |t|
      t.string :gs_user_id
      t.string :gs_user_name
      t.string :notice_id
      t.datetime :published
      t.string :content_text
      t.string :content_html
      t.string :url
      t.timestamps
    end
  end
end
