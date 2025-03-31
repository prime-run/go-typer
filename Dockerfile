FROM golang:1.24.1-alpine

ENV COLORTERM=truecolor

RUN apk add --no-cache zsh

# Set the working directory
WORKDIR /app

COPY . .

RUN go build -o go-typer .

ENTRYPOINT ["./go-typer"]
CMD ["start"]

