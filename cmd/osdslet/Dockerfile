# Dock build usage:	docker build . -t leonwanghui/opensds:v1alpha1
# Docker run usage:	docker run -it --net=host -v /var/log/opensds:/var/log/opensds leonwanghui/opensds:v1alpha1

FROM ubuntu
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdslet /usr/bin

# Define default command.
CMD ["/usr/bin/osdslet", "--endpoint=YOUR_HOST_IP:50040"]