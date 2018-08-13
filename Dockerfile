FROM busybox:glibc
EXPOSE 12345
ADD shopping-list /bin/shopping-list
CMD ["shopping-list"]
