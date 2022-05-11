package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type MergeRequestService struct {
	*BaseService
}

type MergeRequest struct {
	SourceBranch *string `json:"source_branch,omitempty" yaml:"source_branch,omitempty"`
	TargetBranch *string `json:"target_branch,omitempty" yaml:"target_branch,omitempty"`
	Title        *string `json:"title,omitempty" yaml:"title,omitempty"`
}

func NewMergeRequestService(p *Provisioner) *MergeRequestService {
	return &MergeRequestService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

func (m *MergeRequestService) Create(pc *ProjectCtx) {
	for _, mergeRequest := range pc.Project.MergeRequests {
		_, _, err := m.provisioner.Client.MergeRequests.CreateMergeRequest(pc.ProjectID, &gitlab.CreateMergeRequestOptions{
			SourceBranch: mergeRequest.SourceBranch,
			TargetBranch: mergeRequest.TargetBranch,
			Title:        mergeRequest.Title,
		})
		if err != nil {
			fmt.Printf("[ERROR] Merge request `%s` could not be created for project `%s` -- %s\n", *mergeRequest.Title, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Created merge request `%s` for project `%s`.\n", *mergeRequest.Title, pc.ProjectID)
		}
	}
}
