FROM golang:1.13
# Clone the project to local
RUN git clone https://github.com/moby/moby.git /go/src/github.com/moby/moby

# Install package dependencies
RUN apt-get update && \
	apt-get install -y vim python3

# Clone git porject dependencies


# Get go package dependencies


# Checkout the fixed version of this bug
WORKDIR /go/src/github.com/moby/moby
RUN git reset --hard 48ed4f0639d2f290603a04ec146beb3f9569280f




RUN sed -i '68s/--rm//' Makefile && \
	sed -i '68s/MOUNT)/MOUNT) -v \/go\/test:\/go\/test --name moby_29733_cntr/' Makefile && \
	sed -i '52s/-cover//' hack/make/test-unit && \
	sed -i '52s/go test.*/&\n\t&/' hack/make/test-unit && \
	sed -i '52s/$/ -i/' hack/make/test-unit && \
	sed -i '53s/\$pkg_list//' hack/make/test-unit && \
	sed -i '53s/$/ -c -o \/go\/gobench.test/' hack/make/test-unit