FROM busybox:glibc
EXPOSE 12345
ADD main /bin/shopping-list
CMD ["shopping-list"]
