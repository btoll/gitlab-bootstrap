package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/xanzy/go-gitlab"
)

type APIToken string

type BaseService struct {
	provisioner *Provisioner
}

type Provisioner struct {
	Client *gitlab.Client
	Token  APIToken

	Configs []Group

	Branches      *BranchService
	Groups        *GroupService
	Invites       *InviteService
	Issues        *IssueService
	Labels        *LabelService
	MergeRequests *MergeRequestService
	Projects      *ProjectService
	Releases      *ReleaseService
	//	Users    *UserService
}

func NewProvisioner(configs []Group) *Provisioner {
	apiToken, isSet := os.LookupEnv("GITLAB_API_PRIVATE_TOKEN")
	if apiToken == "" || !isSet {
		panic("[ERROR] Must set $GITLAB_API_PRIVATE_TOKEN!")
	}

	gitlabClient, err := gitlab.NewClient(apiToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	p := &Provisioner{
		Client:  gitlabClient,
		Token:   APIToken(apiToken),
		Configs: configs,
	}

	p.Branches = NewBranchService(p)
	p.Groups = NewGroupService(p)
	p.Invites = NewInviteService(p)
	p.Issues = NewIssueService(p)
	p.Labels = NewLabelService(p)
	p.MergeRequests = NewMergeRequestService(p)
	p.Projects = NewProjectService(p)
	p.Releases = NewReleaseService(p)
	//	p.Users = NewUserService(p)

	return p
}

func (p *Provisioner) ProcessConfig(config Group, destroy bool) {
	var wg sync.WaitGroup
	wg.Add(len(config.Projects))

	group, err := p.Groups.Get(config.Group)
	//	LookupGroup(group)
	if err != nil {
		panic(err)
	}

	if group == nil && config.Parent != nil {
		group = p.Groups.CreateSubgroup(config)
	}

	p.Groups.Mapper[config.Group] = group.ID

	for _, project := range config.Projects {
		pc := ProjectCtx{
			Group:     group,
			Project:   project,
			ProjectID: fmt.Sprintf("%s/%s", group.FullPath, project.Name),
		}
		go func(pc ProjectCtx) {
			if !destroy {
				//				p.Projects.Replace()
				p.Projects.Create(&pc)
			} else {
				p.Projects.Delete(&pc)
			}
			wg.Done()
		}(pc)
	}

	wg.Wait()
}

func (p *Provisioner) ProcessConfigs(destroy bool) {
	for _, config := range p.Configs {
		p.ProcessConfig(config, destroy)
	}
}
