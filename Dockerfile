FROM alpine:3.5

RUN mkdir -pv /home/ottemo/media
RUN mkdir -pv /home/ottemo/foundation
RUN mkdir -pv /home/ottemo/foundation/var/log
RUN mkdir -pv /home/ottemo/foundation/var/session

COPY foundation /home/ottemo/foundation/

# create links for proper logging
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/cron.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/events.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/paypal.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/rest.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/subscription.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/mongo.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/events.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/models.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/usps.log
RUN ln -sf /dev/stdout /home/ottemo/foundation/var/log/product
RUN ln -sf /dev/stderr /home/ottemo/foundation/var/log/errors.log

COPY bin/docker-entrypoint.sh /home/ottemo/foundation/

EXPOSE 3000
WORKDIR /home/ottemo/foundation
CMD ./docker-entrypoint.sh
