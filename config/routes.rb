Rails.application.routes.draw do

  mount RailsAdmin::Engine => '/admin', as: 'rails_admin'
  # Create all routes related to messages
  resources :messages, only: [:index, :show]

  # Test route
  get 'home', to: 'home#home'

end
