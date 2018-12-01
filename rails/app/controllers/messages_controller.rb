class MessagesController < ApplicationController
  skip_before_action :verify_authenticity_token, only: [:create]
  def index
    @messages = Message.all
    render json: @messages
  end

  def show
    @message = Message.find(params[:id])
    render json: @message
  end

  def create
    link = params['link']
    if link.nil?
      link = 'fsf://fsf'
    end
    Message.create title: params['title'], content: params['content'], link: link
    render plain: "Successfully created notification"
  end
end
