docker run -d \
--name user_service \
-p 5431:5432 \
-e POSTGRES_USER=user \
-e POSTGRES_PASSWORD=password \
-e POSTGRES_DB=user_service_db postgres;