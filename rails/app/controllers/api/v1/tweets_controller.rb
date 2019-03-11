module Api
  module V1
    class TweetsController < ApplicationController
      def index
        tweets = Tweet.order('date DESC');
        render json: {status: 'SUCCESS', message:'Loaded FSF Tweets', data: tweets}, status: :ok
      end

      def show
        tweet = Tweet.find(params[:id])
        render json: {status: 'SUCCESS', message:'Loaded FSF Tweet', data: tweet}, status: :ok
      end

      def create
        tweet = Tweet.new(tweet_params)
        if tweet.save
          render json: {status: 'SUCCESS', message:'Loaded FSF Tweet', data: tweet}, status: :ok
        else
          render json: {status: 'ERROR', message:'Tweet not saved', data: tweet.errors}, status: :unprocessable_entity
        end
      end

      private
      def tweet_params
        params.permit(
          :id,
          :date,
          :url,
          :text
        )
      end
    end
  end
end
