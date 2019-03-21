module Api
  module V1
    class NoticesController < ApplicationController
      def index
        notices = Notice.order('published DESC').limit(20)
        render json: { status: 'SUCCESS', message: 'Loaded FSF notices', data: notices },
               status: :ok
      end

      def show
        notice = Notice.find(params[:id])
        render json: { status: 'SUCCESS', message: 'Loaded FSF Notice', data: notice },
               status: :ok
      end

      def create
        notice = Notice.new(notice_params)
        if notice.save
          render json: { status: 'SUCCESS', message: 'Loaded FSF notice', data: notice },
                 status: :ok
        else
          render json: {
                   status: 'ERROR', message: 'Notice not saved', data: notice.errors
                 },
                 status: :unprocessable_entity
        end
      end

      private
      def notice_params
        params.permit(
          id,
          gs_user_id,
          gs_user_name,
          published,
          content_text,
          content_html,
          url
        )
      end
    end
  end
end
