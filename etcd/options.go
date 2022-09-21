package etcd

// ClientOption etcd client proxy option
type ClientOption func(*clientConfig)

// WithEndpoints set endpoints
func WithEndpoints(endpoints []string) ClientOption {
	return func(cc *clientConfig) {
		cc.Endpoints = endpoints
	}
}

// WithTimeout set timeout
//
// default 1000 millisecond
func WithTimeout(timeout int) ClientOption {
	return func(cc *clientConfig) {
		cc.Timeout = timeout
	}
}

// WithAuth set username and password for authentication
func WithAuth(username, password string) ClientOption {
	return func(cc *clientConfig) {
		cc.Username = username
		cc.Password = password
	}
}

// WithTLSKeyPath set tls key file path
func WithTLSKeyPath(tlsKeyPath string) ClientOption {
	return func(cc *clientConfig) {
		cc.TLSKeyPath = tlsKeyPath
	}
}

// WithTLSCertPath set tls cert file path
func WithTLSCertPath(tlsCertPath string) ClientOption {
	return func(cc *clientConfig) {
		cc.TLSCertPath = tlsCertPath
	}
}

// WithCACertPath set ca cert file path
func WithCACertPath(caCertPath string) ClientOption {
	return func(cc *clientConfig) {
		cc.CACertPath = caCertPath
	}
}
