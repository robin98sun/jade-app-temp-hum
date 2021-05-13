FROM debian
COPY ./jade-app /jade-app
EXPOSE 8080
ENTRYPOINT /jade-app
