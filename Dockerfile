# Gunakan image dasar golang versi 1.21
FROM golang:1.21-alpine

# Setel direktori kerja
WORKDIR /app

# Salin go.mod dan go.sum
COPY go.mod go.sum ./

# Unduh dependensi
RUN go mod download

# Salin seluruh source code ke direktori kerja
COPY . .

# Kompilasi aplikasi
RUN go build -o main ./cmd/main.go

# Eksekusi aplikasi
CMD ["./main"]
