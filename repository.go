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
	ForkPolicy  string
	Type        string
	Owner       map[string]interface{}
	Links       map[string]interface{}
}

type Pipeline struct {
	Type       string
	Enabled    bool
	Repository Repository
}

type PipelineVariable struct {
	Type    string
	Uuid    string
	Key     string
	Value   string
	Secured bool
}

type PipelineKeyPair struct {
	Type       string
	Uuid       string
	PublicKey  string
	PrivateKey string
}

func (r *Repository) Create(ro *RepositoryOptions) (*Repository, error) {
	data := r.buildRepositoryBody(ro)
	urlStr := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.RepoSlug)
	response, err := r.c.execute("POST", urlStr, data)
	if err != nil {
		return nil, err
	}

	return decodeRepository(response)
}

func (r *Repository) Get(ro *RepositoryOptions) (*Repository, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.RepoSlug)
	response, err := r.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeRepository(response)
}

func (r *Repository) Delete(ro *RepositoryOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.RepoSlug)
	return r.c.execute("DELETE", urlStr, "")
}

func (r *Repository) GetCommit(gco *GetCommitOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s/commit/%s", gco.Owner, gco.RepoSlug, gco.Node)
	return r.c.execute("GET", urlStr, "")
}

func (r *Repository) ListWatchers(ro *RepositoryOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s/watchers", ro.Owner, ro.RepoSlug)
	return r.c.execute("GET", urlStr, "")
}

func (r *Repository) ListForks(ro *RepositoryOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s/forks", ro.Owner, ro.RepoSlug)
	return r.c.execute("GET", urlStr, "")
}

func (r *Repository) UpdatePipelineConfig(rpo *RepositoryPipelineOptions) (*Pipeline, error) {
	data := r.buildPipelineBody(rpo)
	urlStr := r.c.requestUrl("/repositories/%s/%s/pipelines_config", rpo.Owner, rpo.RepoSlug)
	response, err := r.c.execute("PUT", urlStr, data)
	if err != nil {
		return nil, err
	}
	return decodePipelineRepository(response)
}

func (r *Repository) AddPipelineVariable(rpvo *RepositoryPipelineVariableOptions) (*PipelineVariable, error) {
	data := r.buildPipelineVariableBody(rpvo)
	urlStr := r.c.requestUrl("/repositories/%s/%s/pipelines_config/variables/", rpvo.Owner, rpvo.RepoSlug)

	response, err := r.c.execute("POST", urlStr, data)
	if err != nil {
		return nil, err
	}
	return decodePipelineVariableRepository(response)
}

func (r *Repository) AddPipelineKeyPair(rpkpo *RepositoryPipelineKeyPairOptions) (*PipelineKeyPair, error) {
	data := r.buildPipelineKeyPairBody(rpkpo)
	urlStr := r.c.requestUrl("/repositories/%s/%s/pipelines_config/ssh/key_pair", rpkpo.Owner, rpkpo.RepoSlug)

	response, err := r.c.execute("PUT", urlStr, data)
	if err != nil {
		return nil, err
	}

	return decodePipelineKeyPairRepository(response)
}

//Sourcecontent used to retrieve the contents of a single file,
// or the contents of a directory at a specified revision
func (r *Repository) SourceContent(sco *SourceContentOption) (interface{}, error) {
	args := []interface{}{sco.Owner, sco.RepoSlug, sco.Revision, sco.Path}
	urlStr := r.c.requestUrl("/repositories/%s/%s/src/%s/%s", args...)
	return r.c.execute("GET", urlStr, "")
}

// CreateBuildStatus create build status against a repository on bitbucket
func (r *Repository) CreateBuildStatus(so *CreateBuildStatusOptions) (interface{}, error) {
	args := []interface{}{so.Owner, so.RepoSlug, so.Node}
	urlStr := r.c.requestUrl("/repositories/%s/%s/commit/%s/statuses/build", args...)
	body := r.buildJsonBody(so.Body)
	return r.c.execute("POST", urlStr, body)
}

func (r *Repository) buildJsonBody(body map[string]interface{}) string {

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}

func (r *Repository) buildRepositoryBody(ro *RepositoryOptions) string {

	body := map[string]interface{}{}

	if ro.Scm != "" {
		body["scm"] = ro.Scm
	}
	//if ro.Scm != "" {
	//		body["name"] = ro.Name
	//}
	if ro.IsPrivate != "" {
		body["is_private"] = ro.IsPrivate
	}
	if ro.Description != "" {
		body["description"] = ro.Description
	}
	if ro.ForkPolicy != "" {
		body["fork_policy"] = ro.ForkPolicy
	}
	if ro.Language != "" {
		body["language"] = ro.Language
	}
	if ro.HasIssues != "" {
		body["has_issues"] = ro.HasIssues
	}
	if ro.HasWiki != "" {
		body["has_wiki"] = ro.HasWiki
	}
	if ro.Project != "" {
		body["project"] = map[string]string{
			"key": ro.Project,
		}
	}

	return r.buildJsonBody(body)
}

func (r *Repository) buildPipelineBody(rpo *RepositoryPipelineOptions) string {

	body := map[string]interface{}{}

	body["enabled"] = rpo.Enabled

	return r.buildJsonBody(body)
}

func (r *Repository) buildPipelineVariableBody(rpvo *RepositoryPipelineVariableOptions) string {

	body := map[string]interface{}{}

	if rpvo.Uuid != "" {
		body["uuid"] = rpvo.Uuid
	}
	body["key"] = rpvo.Key
	body["value"] = rpvo.Value
	body["secured"] = rpvo.Secured

	return r.buildJsonBody(body)
}

func (r *Repository) buildPipelineKeyPairBody(rpkpo *RepositoryPipelineKeyPairOptions) string {

	body := map[string]interface{}{}

	if rpkpo.PrivateKey != "" {
		body["private_key"] = rpkpo.PrivateKey
	}
	if rpkpo.PublicKey != "" {
		body["public_key"] = rpkpo.PublicKey
	}

	return r.buildJsonBody(body)
}

func decodeRepository(repoResponse interface{}) (*Repository, error) {
	repoMap := repoResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var repository = new(Repository)
	err := mapstructure.Decode(repoMap, repository)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func decodePipelineRepository(repoResponse interface{}) (*Pipeline, error) {
	repoMap := repoResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var pipeline = new(Pipeline)
	err := mapstructure.Decode(repoMap, pipeline)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func decodePipelineVariableRepository(repoResponse interface{}) (*PipelineVariable, error) {
	repoMap := repoResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var pipelineVariable = new(PipelineVariable)
	err := mapstructure.Decode(repoMap, pipelineVariable)
	if err != nil {
		return nil, err
	}

	return pipelineVariable, nil
}

func decodePipelineKeyPairRepository(repoResponse interface{}) (*PipelineKeyPair, error) {
	repoMap := repoResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var pipelineKeyPair = new(PipelineKeyPair)
	err := mapstructure.Decode(repoMap, pipelineKeyPair)
	if err != nil {
		return nil, err
	}

	return pipelineKeyPair, nil
}
