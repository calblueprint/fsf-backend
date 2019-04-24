module Api
    module V1
        class MessagesController < ApplicationController
            def get_latest_messages


                # Get the stringified parameter passed in from the end point
                # un-stringify deadtime
                two_hours_ago = DateTime.now - (5.0/24) # dummy value
                last_sent = two_hours_ago

                # filter Message.where('updated_at' > something)

                # messages = Message.order('updated_at DESC').last(1);

                messages = Message.where(updated_at: last_sent..DateTime.now).order('updated_at DESC')
        

                # @users = User.where(created_at: from_date.beginning_of_day..to_date.end_of_day)


                # upon deletion, also do something (not found here)
                
                render json: {status: 'SUCCESS', message:'Loaded FSF Messages', data: messages}, status: :ok
            end
        end
    end
end
