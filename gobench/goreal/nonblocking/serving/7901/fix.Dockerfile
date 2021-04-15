FROM golang:1.13
# Clone the project to local
RUN git clone https://github.com/knative/serving.git /go/src/github.com/knative/serving

# Install package dependencies
RUN apt-get update && \
	apt-get install -y vim python3

# Clone git porject dependencies


# Get go package dependencies


# Checkout the fixed version of this bug
WORKDIR /go/src/github.com/knative/serving
RUN git reset --hard 5d91bee9539f9051a3e25cf99dc73d4e0bc9829b




RUN go test ./pkg/reconciler/autoscaling/hpa -race -c -o /go/gobench.test