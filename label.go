package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type LabelService struct {
	*BaseService
}

type Label struct {
	Name        string  `json:"name,omitempty" yaml:"name,omitempty"`
	Color       string  `json:"color,omitempty" yaml:"color,omitempty"`
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty" yaml:"priority,omitempty"`
}

func NewLabelService(p *Provisioner) *LabelService {
	return &LabelService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

func (l *LabelService) Create(pc *ProjectCtx) {
	for _, label := range pc.Project.Labels {
		_, _, err := l.provisioner.Client.Labels.CreateLabel(pc.ProjectID, &gitlab.CreateLabelOptions{
			Name:        &label.Name,
			Color:       &label.Color,
			Description: label.Description,
			Priority:    label.Priority,
		})
		if err != nil {
			fmt.Printf("[ERROR] Label `%s` could not be created for project `%s` -- %s\n", label.Name, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Created label `%s` for project `%s`.\n", label.Name, pc.ProjectID)
		}
	}
}
