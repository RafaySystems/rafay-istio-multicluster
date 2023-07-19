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

/***********************************************************************
4.1.2.5.1.  UTCTime
   The universal time type, UTCTime, is a standard ASN.1 type intended
   for representation of dates and time.  UTCTime specifies the year
   through the two low-order digits and time is specified to the
   precision of one minute or one second.  UTCTime includes either Z
   (for Zulu, or Greenwich Mean Time) or a time differential.

   For the purposes of this profile, UTCTime values MUST be expressed in
   Greenwich Mean Time (Zulu) and MUST include seconds (i.e., times are
   YYMMDDHHMMSSZ), even where the number of seconds is zero.  Conforming
   systems MUST interpret the year field (YY) as follows:

      Where YY is greater than or equal to 50, the year SHALL be
      interpreted as 19YY; and

      Where YY is less than 50, the year SHALL be interpreted as 20YY.
***********************************************************************/

import (
	"time"

	"github.com/smallstep/zcrypto/x509"
	"github.com/smallstep/zlint/util"
)

type utcTimeGMT struct {
}

func (l *utcTimeGMT) Initialize() error {
	return nil
}

func (l *utcTimeGMT) CheckApplies(c *x509.Certificate) bool {
	firstDate, secondDate := util.GetTimes(c)
	beforeTag, afterTag := util.FindTimeType(firstDate, secondDate)
	date1Utc := beforeTag == 23
	date2Utc := afterTag == 23
	return date1Utc || date2Utc
}

func (l *utcTimeGMT) Execute(c *x509.Certificate) *LintResult {
	var r LintStatus
	firstDate, secondDate := util.GetTimes(c)
	beforeTag, afterTag := util.FindTimeType(firstDate, secondDate)
	date1Utc := beforeTag == 23
	date2Utc := afterTag == 23
	if date1Utc {
		// UTC Tests on notBefore
		utcNotGmt(c.NotBefore, &r)
	}
	if date2Utc {
		// UTC Tests on NotAfter
		utcNotGmt(c.NotAfter, &r)
	}
	return &LintResult{Status: r}
}

func utcNotGmt(t time.Time, r *LintStatus) {
	// If we already ran this test and it resulted in error, don't want to discard that
	// And now we use the afterBool to make sure we test the right time
	if *r == Error {
		return
	}
	if t.Location() != time.UTC {
		*r = Error
	} else {
		*r = Pass
	}
}

func init() {
	RegisterLint(&Lint{
		Name:          "e_utc_time_not_in_zulu",
		Description:   "UTCTime values MUST be expressed in Greenwich Mean Time (Zulu)",
		Citation:      "RFC 5280: 4.1.2.5.1",
		Source:        RFC5280,
		EffectiveDate: util.RFC2459Date,
		Lint:          &utcTimeGMT{},
	})
}