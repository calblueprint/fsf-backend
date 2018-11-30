#This module is the API for news Articles
module Api
    module V1
        class ArticlesController < ApplicationController
            def index
                articles = Article.order('pub_date DESC');
                render json: {status: 'SUCCESS', message:'Loaded FSF Articles', data: articles}, status: :ok
            end
            
            def show
                article = Article.find(params[:id])
                render json: {status: 'SUCCESS', message:'Loaded FSF Article', data: article}, status: :ok
            end

            def create
                article = Article.new(article_params)
                if article.save
                    render json: {status: 'SUCCESS', message:'Loaded FSF Article', data: article}, status: :ok
                else
                    render json: {status: 'ERROR', message:'Article not saved', data: article.errors}, status: :unprocessable_entity
                end
            end

            private
            def article_params
                params.permit(
                    :title,
                    :link,
                    :pub_date,
                    :content,
                    :news_alert,
                )
            end
        end
    end
end