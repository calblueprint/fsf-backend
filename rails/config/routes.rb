Rails.application.routes.draw do
  namespace 'api' do
    namespace 'v1' do
      resources :articles
    end
  end
  # Create all routes related to messages
  resources :messages, only: [:index, :show]

  # Test route
  get 'home', to: 'home#home'
  resources :petitions, only: [:index]

end
