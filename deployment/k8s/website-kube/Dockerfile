FROM centos:7

EXPOSE 80

COPY website /root/website

RUN \
   ls /root/website; \
   chmod 755 /root/website/website;

CMD ["/root/website/website", "-port", "80"]
