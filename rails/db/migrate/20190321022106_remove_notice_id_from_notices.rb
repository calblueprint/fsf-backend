class RemoveNoticeIdFromNotices < ActiveRecord::Migration[5.2]
  def change
    remove_column :notices, :notice_id, :string
  end
end
