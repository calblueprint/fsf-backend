#This module is the API for news Petitions
module Api
    module V1
        class PetitionsController < ApplicationController
            def index
                petitions = Petition.order('created_at DESC');
                render json: {status: 'SUCCESS', message:'Loaded FSF Petitions', data: petitions}, status: :ok
            end

            def show
                petition = Petition.find(params[:id])
                render json: {status: 'SUCCESS', message:'Loaded FSF Petition', data: petitions}, status: :ok
            end

            def create
                petition = Petition.new(petition_params)
                if petition.save
                    render json: {status: 'SUCCESS', message:'Created new FSF Petition', data: petition}, status: :ok
                else
                    render json: {status: 'ERROR', message:'Petition not saved', data: petition.errors}, status: :unprocessable_entity
                end
            end

            private
            def petition_params
                params.permit(
                    :title,
                    :link,
                    :description
                )
            end
        end
    end
end