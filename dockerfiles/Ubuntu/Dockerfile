FROM golang:1.15

WORKDIR /root
RUN \
  git clone https://github.com/Chathuru/kubernetes-cluster-autoscaler.git && \
  cd kubernetes-cluster-autoscaler && \
  go build -o autoscaler  cmd/main.go


FROM golang:1.15
WORKDIR /go/bin
COPY --from=0 /root/kubernetes-cluster-autoscaler/autoscaler .
CMD ["./autoscaler"]
