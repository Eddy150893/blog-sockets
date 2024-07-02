#Construir la imagen con docker build . -t platzi-ws-rest-db
#Correr el contenedor docker run -p 54321:5432 platzi-ws-rest-db
FROM postgres:10.3

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]