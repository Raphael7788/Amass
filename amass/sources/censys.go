// Copyright 2017 Jeff Foley. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package sources

import (
	"fmt"

	"github.com/OWASP/Amass/amass/core"
	"github.com/OWASP/Amass/amass/utils"
)

type Censys struct {
	BaseDataSource
}

func NewCensys(srv core.AmassService) DataSource {
	c := new(Censys)

	c.BaseDataSource = *NewBaseDataSource(srv, core.CERT, "Censys")
	return c
}

func (c *Censys) Query(domain, sub string) []string {
	var unique []string

	if domain != sub {
		return unique
	}

	url := c.getURL(domain)
	page, err := utils.GetWebPage(url, nil)
	if err != nil {
		c.Service.Config().Log.Printf("%s: %v", url, err)
		return unique
	}

	c.Service.SetActive()
	re := utils.SubdomainRegex(domain)
	for _, sd := range re.FindAllString(page, -1) {
		if u := utils.NewUniqueElements(unique, sd); len(u) > 0 {
			unique = append(unique, u...)
		}
	}
	return unique
}

func (c *Censys) getURL(domain string) string {
	format := "https://www.censys.io/domain/%s/table"

	return fmt.Sprintf(format, domain)
}
