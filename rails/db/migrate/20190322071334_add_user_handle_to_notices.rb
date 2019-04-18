class AddUserHandleToNotices < ActiveRecord::Migration[5.2]
  def change
    add_column :notices, :gs_user_handle, :string
  end
end
