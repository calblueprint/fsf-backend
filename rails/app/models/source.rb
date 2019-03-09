class Source < ApplicationRecord
  enum source_type: %i[rss]

  validates :url, presence:true
end
