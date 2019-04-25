module Api
    module V1
        class MessagesController < ApplicationController
            def get_latest_messages
                if(params.has_key?(:last_sent))
                    # Get the stringified parameter passed in from the end point
                    last_sent_param = params['last_sent']
                    
                    # parse time
                    last_sent = Time.parse(last_sent_param)
                    
                    # filter Message.where('updated_at' > something)
                    messages = Message.where(updated_at: last_sent..DateTime.now).order('updated_at DESC')
                    render json: {status: 'SUCCESS', message:'Loaded FSF Messages', data: messages}, status: :ok
                else
                    # messages = Message.order('updated_at DESC').last(1);
                    render json: {status: 'FAILURE', message:'Please pass in a last_sent parameter', data: {}}, status: :bad_request
                end
            end
        end
    end
end
