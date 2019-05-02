Rails.application.routes.draw do
  namespace 'api' do
    namespace 'v1' do
      resources :articles
      resources :notices
      resources :tweets
      resources :petitions
      resources :messages

      # resources :articles, only: [:show]
      # resources :notices, only: [:show]
      # resources :petitions, only: [:show]
      # resources :tweets, only: [:show]

      # message route
      get 'latestMessages', to: 'messages#get_latest_messages'
    end
  end

  mount RailsAdmin::Engine => '/admin', as: 'rails_admin'
  # Create all routes related to messages
  resources :messages, only: [:index, :show, :create]

  # Test route
  get 'home', to: 'home#home'
  # resources :petitions, only: [:index]

end
