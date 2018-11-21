class MessagesController < ApplicationController
  def index
    @messages = Message.all
    render json: @messages
  end

  def show
    @message = Message.find(params[:id])
    render json: @message
  end

  def create
    Message.create title: "Brand new notification", content: "Notification created by #{params[:device_id]} at #{DateTime.now}"
  end
end
