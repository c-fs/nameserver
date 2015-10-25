FROM scratch

ADD nameserver /
ADD server/nameserver.conf /

EXPOSE 15525
ENTRYPOINT ["/nameserver"]
