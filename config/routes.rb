Rails.application.routes.draw do

  # Create all routes related to messages
  resources :messages, only: [:index]

  # Test route
  get 'home', to: 'home#home'

end
