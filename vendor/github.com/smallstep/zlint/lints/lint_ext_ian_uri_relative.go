package lints

/*
 * ZLint Copyright 2018 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

/*************************************************************************
When the issuerAltName extension contains a URI, the name MUST be
stored in the uniformResourceIdentifier (an IA5String).  The name
MUST NOT be a relative URI, and it MUST follow the URI syntax and
encoding rules specified in [RFC3986].  The name MUST include both a
scheme (e.g., "http" or "ftp") and a scheme-specific-part.  URIs that
include an authority ([RFC3986], Section 3.2) MUST include a fully
qualified domain name or IP address as the host.  Rules for encoding
Internationalized Resource Identifiers (IRIs) are specified in
Section 7.4.
*************************************************************************/

import (
	"net/url"

	"github.com/smallstep/zcrypto/x509"
	"github.com/smallstep/zlint/util"
)

type uriRelative struct{}

func (l *uriRelative) Initialize() error {
	return nil
}

func (l *uriRelative) CheckApplies(c *x509.Certificate) bool {
	return util.IsExtInCert(c, util.IssuerAlternateNameOID)
}

func (l *uriRelative) Execute(c *x509.Certificate) *LintResult {
	for _, uri := range c.IANURIs {
		parsed_uri, err := url.Parse(uri)

		if err != nil {
			return &LintResult{Status: Error}
		}

		if !parsed_uri.IsAbs() {
			return &LintResult{Status: Error}
		}
	}
	return &LintResult{Status: Pass}
}

func init() {
	RegisterLint(&Lint{
		Name:          "e_ext_ian_uri_relative",
		Description:   "When issuerAltName extension is present and the URI is used, the name MUST NOT be a relative URI",
		Citation:      "RFC 5280: 4.2.1.7",
		Source:        RFC5280,
		EffectiveDate: util.RFC5280Date,
		Lint:          &uriRelative{},
	})
}
