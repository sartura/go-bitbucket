package bitbucket

import (
	"encoding/json"
	"os"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type Project struct {
	Key  string
	Name string
}

type Repository struct {
	c *Client

	Project     Project
	Slug        string
	Full_name   string
	Description string
	Fork_policy string
	Type        string
	Owner       map[string]interface{}
	Links       map[string]interface{}
}

func (r *Repository) Create(ro *RepositoryOptions) (Repository, error) {
	data := r.buildRepositoryBody(ro)
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.Repo_slug)
	response := r.c.execute("POST", url, data)

	return decodeRepository(response)
}

func (r *Repository) Get(ro *RepositoryOptions) (Repository, error) {
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.Repo_slug)
	response := r.c.execute("GET", url, "")

	return decodeRepository(response)
}

func (r *Repository) Delete(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.Repo_slug)
	return r.c.execute("DELETE", url, "")
}

func (r *Repository) ListWatchers(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/watchers", ro.Owner, ro.Repo_slug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) ListForks(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/forks", ro.Owner, ro.Repo_slug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) buildRepositoryBody(ro *RepositoryOptions) string {

	body := map[string]interface{}{}

	if ro.Scm != "" {
		body["scm"] = ro.Scm
	}
	//if ro.Scm != "" {
	//		body["name"] = ro.Name
	//}
	if ro.Is_private != "" {
		body["is_private"] = ro.Is_private
	}
	if ro.Description != "" {
		body["description"] = ro.Description
	}
	if ro.Fork_policy != "" {
		body["fork_policy"] = ro.Fork_policy
	}
	if ro.Language != "" {
		body["language"] = ro.Language
	}
	if ro.Has_issues != "" {
		body["has_issues"] = ro.Has_issues
	}
	if ro.Has_wiki != "" {
		body["has_wiki"] = ro.Has_wiki
	}
	if ro.Project != "" {
		body["project"] = map[string]string{
			"key": ro.Project,
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}

func decodeRepository(json interface{}) (Repository, error) {
	jsonMap := json.(map[string]interface{})

	if jsonMap["type"] == "error" {
		return Repository{}, DecodeError(jsonMap)
	}

	var repository Repository
	err := mapstructure.Decode(jsonMap, &repository)
	if err != nil {
		return Repository{}, err
	}

	return repository, nil
}
