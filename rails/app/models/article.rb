class Article < ApplicationRecord
    validates :headline, presence: true
    validates :lead, presence: true
    validates :pub_date, presence: true
    validates :category, presence:true
    validates :author, presence:true
    validates :content, presence: true
end