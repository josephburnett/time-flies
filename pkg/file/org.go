package file

import (
	"bytes"
	"io/ioutil"
	"strings"

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
	done, todo, err := sectionEntries(d.Outline.Section)
	if err != nil {
		return nil, err
	}
	return types.Log{&types.Week{
		// TODO: date from COMPLETED or CREATED
		// TODO: header from PROPERTIES
		Done: done,
		Todo: todo,
	}}, nil
}

func sectionEntries(s *org.Section) (done, todo []*types.Entry, err error) {
	if s == nil {
		return
	}
	if s.Headline != nil {
		if s.Headline.Status == "TODO" {
			lines := []string{}
			for _, n := range s.Headline.Title {
				lines = append(lines, n.String())
			}
			todo = append(todo, &types.Entry{
				Line:   strings.Join(lines, " "),
				Labels: tagLabels(s.Headline.Tags),
			})
		}
		if s.Headline.Status == "DONE" {
			lines := []string{}
			for _, n := range s.Headline.Title {
				lines = append(lines, n.String())
			}
			done = append(done, &types.Entry{
				Line:   strings.Join(lines, " "),
				Labels: tagLabels(s.Headline.Tags),
			})
		}
	}
	for _, c := range s.Children {
		cDone, cTodo, err := sectionEntries(c)
		if err != nil {
			return nil, nil, err
		}
		done = append(done, cDone...)
		todo = append(todo, cTodo...)
	}
	return
}

func tagLabels(tags []string) map[string]string {
	labels := map[string]string{}
	for _, t := range tags {
		t = strings.ReplaceAll(t, "_", "-")
		parts := strings.Split(t, "@")
		if len(parts) != 2 {
			continue
		}
		k, v := parts[0], parts[1]
		if k == "" || v == "" {
			continue
		}
		labels[k] = v
	}
	return labels
}
