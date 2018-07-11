FROM debian:jessie

ADD install.sh install.sh

RUN chmod +x install.sh && ./install.sh && rm install.sh

EXPOSE 8900
EXPOSE 8901

VOLUME /data

ENTRYPOINT ["gopusher-comet"]
