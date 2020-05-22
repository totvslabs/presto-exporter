FROM scratch
LABEL maintainer="devops@totvslabs.com"
COPY presto-exporter /bin/presto-exporter
ENTRYPOINT ["/bin/presto-exporter"]
CMD [ "-h" ]
