#!/bin/bash
./node_modules/.bin/prettier --write lib/tasks/db.rake
./node_modules/.bin/prettier --write app/controllers/api/v1/notices_controller.rb
./node_modules/.bin/prettier --write app/models/petition.rb
./node_modules/.bin/prettier --write app/models/article.rb