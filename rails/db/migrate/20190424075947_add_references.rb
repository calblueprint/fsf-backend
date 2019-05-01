class AddReferences < ActiveRecord::Migration[5.2]
  def change
    add_reference :articles, :message
    add_reference :petitions, :message
    add_reference :tweets, :message
    add_reference :notices, :message
  end
end
