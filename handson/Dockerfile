FROM ubuntu:20.04

RUN apt update
RUN apt install -y vim

COPY .vimrc /root/.vimrc

# デフォルト命令
CMD date +"%Y/%m/%d %H:%M:%S (UTC)"
