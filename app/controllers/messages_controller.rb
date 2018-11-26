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
    Message.create title: params['title'], content: params['content']
    render plain: "Successfully created notification"
  end
end
