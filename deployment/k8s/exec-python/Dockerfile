FROM python:latest

EXPOSE 80

COPY exec /root/exec

RUN \
   ls /root/exec; \
   chmod 755 /root/exec/exec;

CMD ["/root/exec/exec", "-port", "80"]
