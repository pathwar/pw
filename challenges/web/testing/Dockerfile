FROM pathwar/supervisord:onbuild
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)
EXPOSE 22 80

# OpenSSH
RUN apt-get update \
 && apt-get install -y -q openssh-server \
 && apt-get clean
RUN mkdir -p /var/run/sshd /var/log/supervisor \
 && sed -i '/^PermitRootLogin/d;$s/$/\nPermitRootLogin yes/' /etc/ssh/sshd_config \
 && echo "root:root" | chpasswd \
 && chsh -s /bin/bash www-data \
 && echo "www-data:www-data" | chpasswd

# TTY.js
RUN apt-get -qq update      \
 && apt-get -y -qq upgrade  \
 && apt-get -y -qq install  \
    nodejs                  \
    npm                     \
    make                    \
    g++                     \
 && apt-get clean
RUN ln -s /usr/bin/nodejs /usr/bin/node
RUN npm install -g tty.js
RUN for pts in 0 1 2 3 4 5 6 7 8; do echo /dev/pts/$pts >> /etc/securetty; done \
 && echo /dev/pty >> /etc/securetty
ADD ./patches/srv/ /srv/
