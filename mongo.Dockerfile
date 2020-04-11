FROM mongo:4.2

# init.js will be executed when the mongo container runs
COPY ./init-mongo.js ./docker-entrypoint-initdb.d
