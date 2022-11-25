FROM ubuntu:20.04

COPY ./change_hpa_container /usr/local/bin/change_hpa_container

ENTRYPOINT [ "/usr/local/bin/change_hpa_container" ]