package ca

import (
	"context"
	"crypto"
	"crypto/tls"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/smallstep/certificates/api"
	"go.step.sm/crypto/jose"
)

type tokenClaims struct {
	SHA string `json:"sha"`
	jose.Claims
}

// Bootstrap is a helper function that initializes a client with the
// configuration in the bootstrap token.
func Bootstrap(token string) (*Client, error) {
	tok, err := jose.ParseSigned(token)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing token")
	}
	var claims tokenClaims
	if err := tok.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, errors.Wrap(err, "error parsing token")
	}

	// Validate bootstrap token
	switch {
	case claims.SHA == "":
		return nil, errors.New("invalid bootstrap token: sha claim is not present")
	case !strings.HasPrefix(strings.ToLower(claims.Audience[0]), "http"):
		return nil, errors.New("invalid bootstrap token: aud claim is not a url")
	}

	return NewClient(claims.Audience[0], WithRootSHA256(claims.SHA))
}

// BootstrapClient is a helper function that using the given bootstrap token
// return an http.Client configured with a Transport prepared to do TLS
// connections using the client certificate returned by the certificate
// authority. By default the server will kick off a routine that will renew the
// certificate after 2/3rd of the certificate's lifetime has expired.
//
// Usage:
//
//	// Default example with certificate rotation.
//	client, err := ca.BootstrapClient(ctx.Background(), token)
//
//	// Example canceling automatic certificate rotation.
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	client, err := ca.BootstrapClient(ctx, token)
//	if err != nil {
//	  return err
//	}
//	resp, err := client.Get("https://internal.smallstep.com")
func BootstrapClient(ctx context.Context, token string, options ...TLSOption) (*http.Client, error) {
	b, err := createBootstrap(token) //nolint:contextcheck // deeply nested context; temporary
	if err != nil {
		return nil, err
	}

	// Make sure the tlsConfig has all supported roots on RootCAs.
	//
	// The roots request is only supported if identity certificates are not
	// required. In all cases the current root is also added after applying all
	// options too.
	if !b.RequireClientAuth {
		options = append(options, AddRootsToRootCAs())
	}

	transport, err := b.Client.Transport(ctx, b.SignResponse, b.PrivateKey, options...)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

// BootstrapServer is a helper function that using the given token returns the
// given http.Server configured with a TLS certificate signed by the Certificate
// Authority. By default the server will kick off a routine that will renew the
// certificate after 2/3rd of the certificate's lifetime has expired.
//
// Without any extra option the server will be configured for mTLS, it will
// require and verify clients certificates, but options can be used to drop this
// requirement, the most common will be only verify the certs if given with
// ca.VerifyClientCertIfGiven(), or add extra CAs with
// ca.AddClientCA(*x509.Certificate).
//
// Usage:
//
//	// Default example with certificate rotation.
//	srv, err := ca.BootstrapServer(context.Background(), token, &http.Server{
//	    Addr: ":443",
//	    Handler: handler,
//	})
//
//	// Example canceling automatic certificate rotation.
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	srv, err := ca.BootstrapServer(ctx, token, &http.Server{
//	    Addr: ":443",
//	    Handler: handler,
//	})
//	if err != nil {
//	    return err
//	}
//	srv.ListenAndServeTLS("", "")
func BootstrapServer(ctx context.Context, token string, base *http.Server, options ...TLSOption) (*http.Server, error) {
	if base.TLSConfig != nil {
		return nil, errors.New("server TLSConfig is already set")
	}

	b, err := createBootstrap(token) //nolint:contextcheck // deeply nested context; temporary
	if err != nil {
		return nil, err
	}

	// Make sure the tlsConfig has all supported roots on RootCAs.
	//
	// The roots request is only supported if identity certificates are not
	// required. In all cases the current root is also added after applying all
	// options too.
	if !b.RequireClientAuth {
		options = append(options, AddRootsToCAs())
	}

	tlsConfig, err := b.Client.GetServerTLSConfig(ctx, b.SignResponse, b.PrivateKey, options...)
	if err != nil {
		return nil, err
	}

	base.TLSConfig = tlsConfig
	return base, nil
}

// BootstrapListener is a helper function that using the given token returns a
// TLS listener which accepts connections from an inner listener and wraps each
// connection with Server.
//
// Without any extra option the server will be configured for mTLS, it will
// require and verify clients certificates, but options can be used to drop this
// requirement, the most common will be only verify the certs if given with
// ca.VerifyClientCertIfGiven(), or add extra CAs with
// ca.AddClientCA(*x509.Certificate).
//
// Usage:
//
//	inner, err := net.Listen("tcp", ":443")
//	if err != nil {
//	  return nil
//	}
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	lis, err := ca.BootstrapListener(ctx, token, inner)
//	if err != nil {
//	    return err
//	}
//	srv := grpc.NewServer()
//	... // register services
//	srv.Serve(lis)
func BootstrapListener(ctx context.Context, token string, inner net.Listener, options ...TLSOption) (net.Listener, error) {
	b, err := createBootstrap(token) //nolint:contextcheck // deeply nested context; temporary
	if err != nil {
		return nil, err
	}

	// Make sure the tlsConfig has all supported roots on RootCAs.
	//
	// The roots request is only supported if identity certificates are not
	// required. In all cases the current root is also added after applying all
	// options too.
	if !b.RequireClientAuth {
		options = append(options, AddRootsToCAs())
	}

	tlsConfig, err := b.Client.GetServerTLSConfig(ctx, b.SignResponse, b.PrivateKey, options...)
	if err != nil {
		return nil, err
	}

	return tls.NewListener(inner, tlsConfig), nil
}

type bootstrap struct {
	Client            *Client
	RequireClientAuth bool
	SignResponse      *api.SignResponse
	PrivateKey        crypto.PrivateKey
}

func createBootstrap(token string) (*bootstrap, error) {
	client, err := Bootstrap(token)
	if err != nil {
		return nil, err
	}

	version, err := client.Version()
	if err != nil {
		return nil, err
	}

	req, pk, err := CreateSignRequest(token)
	if err != nil {
		return nil, err
	}

	sign, err := client.Sign(req)
	if err != nil {
		return nil, err
	}

	return &bootstrap{
		Client:            client,
		RequireClientAuth: version.RequireClientAuthentication,
		SignResponse:      sign,
		PrivateKey:        pk,
	}, nil
}
