package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type BranchService struct {
	*BaseService
}

type Branch struct {
	Branch    string `json:"branch,omitempty" yaml:"branch,omitempty"`
	Ref       string `json:"ref,omitempty" yaml:"ref,omitempty"`
	Protected bool   `json:"protected,omitempty" yaml:"protected,omitempty"`
}

func NewBranchService(p *Provisioner) *BranchService {
	return &BranchService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

func (b *BranchService) Create(pc *ProjectCtx) {
	for _, branch := range pc.Project.Branches {
		_, _, err := b.provisioner.Client.Branches.CreateBranch(pc.ProjectID, &gitlab.CreateBranchOptions{
			Branch: &branch.Branch,
			Ref:    &branch.Ref,
		})
		if err != nil {
			fmt.Printf("[ERROR] Could not create branch `%s` for project `%s` -- %v\n", branch.Branch, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Created branch `%s` for project `%s`.\n", branch.Branch, pc.ProjectID)

			if branch.Protected {
				_, _, err := b.Protect(pc.ProjectID, branch.Branch)
				if err != nil {
					fmt.Printf("[ERROR] Could not protect branch `%s` for project %s` -- %s\n", branch.Branch, pc.ProjectID, err)
				} else {
					fmt.Printf("[INFO] Created branch protections for branch `%s` for project `%s`.\n", branch.Branch, pc.ProjectID)
				}
			}
		}
	}
}

func (b *BranchService) Protect(pid, branch string) (*gitlab.Branch, *gitlab.Response, error) {
	t := true
	return b.provisioner.Client.Branches.ProtectBranch(pid, branch, &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  &t,
		DevelopersCanMerge: &t,
	})
}
