class AddArticleIdToMessages < ActiveRecord::Migration[5.2]
  def change
    add_column :messages, :article_id, :integer
  end
end
