package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/service"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"github.com/carinfinin/shortener-url/internal/app/storage/storefile"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

// Server заускает сервер и содержит ссылку на хранилище.
type Server struct {
	http.Server
	Store  service.Repository
	config *config.Config
}

// New конструктор для Server принимает кофиг.
func New(config *config.Config) (*Server, error) {

	var server Server

	switch {
	case config.DBPath != "":
		s, err := storepg.New(config)
		if err != nil {
			return nil, err
		}
		s.CreateTableForDB(context.Background())
		server.Store = s

	case config.FilePath != "":
		s, err := storefile.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s

	default:
		s, err := store.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s
	}

	server.Addr = config.Addr
	s := service.New(server.Store, config)

	server.Handler = router.ConfigureRouter(s, config).Handle
	server.config = config

	fmt.Println(config)

	return &server, nil
}

// Stop останавливает server
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

// Start запускает сервер.
func (s *Server) Start() error {
	if s.config.TLS {

		m := &autocert.Manager{
			Cache:      autocert.DirCache("secret-dir"),
			Prompt:     autocert.AcceptTOS,
			Email:      "example@example.org",
			HostPolicy: autocert.HostWhitelist(s.config.Addr),
		}
		s.TLSConfig = m.TLSConfig()

		return s.ListenAndServeTLS("", "")
	}
	return s.ListenAndServe()
}

// generateTLS
func generateTLS() {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Сохраняем сертификат в файл
	certFile, _ := os.Create("cert.pem")
	certFile.Write(certPEM.Bytes())
	certFile.Close()

	// Сохраняем приватный ключ в файл
	keyFile, _ := os.Create("key.pem")
	keyFile.Write(privateKeyPEM.Bytes())
	keyFile.Close()

}
