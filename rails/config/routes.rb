Rails.application.routes.draw do

  # Create all routes related to messages
  resources :messages, only: [:index, :show]

  # Test route
  get 'home', to: 'home#home'
  resources :petitions, only: [:index]

end
