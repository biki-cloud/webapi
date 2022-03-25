FROM centos:7

EXPOSE 80

COPY apigw /root/apigw

RUN \
   ls /root/apigw; \
   chmod 755 /root/apigw/apigw;

CMD ["/root/apigw/apigw", "-port", "80"]
