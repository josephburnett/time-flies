package file

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/josephburnett/time-flies/pkg/types"
	"github.com/niklasfasching/go-org/org"
)

func (c *FileConfig) GetOrgFiles() []string {
	if c == nil {
		return nil
	}
	return c.OrgFiles
}

func (c *FileConfig) ReadOrg() (types.Log, error) {
	var allLogs types.Log
	for _, f := range c.GetOrgFiles() {
		bytes, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		log, err := c.ParseOrg(string(bytes))
		if err != nil {
			return nil, err
		}
		allLogs = mergeLogs(allLogs, log)
	}
	return allLogs, nil
}

func (c *FileConfig) ParseOrg(doc string) (types.Log, error) {
	d := org.New().Parse(bytes.NewReader([]byte(doc)), "")
	fmt.Printf("%+v\n", d)
	return nil, nil
}
