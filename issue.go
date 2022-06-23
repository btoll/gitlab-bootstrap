package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type IssueService struct {
	*BaseService
}

type IssueType struct {
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
}

func NewIssueService(p *Provisioner) *IssueService {
	return &IssueService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

const (
	Incident string = "incident"
	Issue    string = "issue"
	TestCase string = "test_case"
)

func getIssueType(issueType string) *string {
	var s string
	switch issueType {
	case Incident:
		s = "incident"
	case Issue:
		s = "issue"
	case TestCase:
		s = "test_case"
	default:
		s = "issue"
	}
	return gitlab.String(s)
}

func (i *IssueService) Create(pc *ProjectCtx) {
	for _, issue := range pc.Project.Issues {
		_, _, err := i.provisioner.Client.Issues.CreateIssue(pc.ProjectID, &gitlab.CreateIssueOptions{
			Title:     &issue.Title,
			IssueType: getIssueType(issue.Type),
		})
		if err != nil {
			fmt.Printf("[ERROR] Issue `%s` could not be created for project `%s` -- %s\n", issue.Title, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Created issue `%s` for project `%s`.\n", issue.Title, pc.ProjectID)
		}
	}
}

func replaceIssues(pc *ProjectCtx, api API) {
}
