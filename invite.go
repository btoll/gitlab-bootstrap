package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type Invite struct {
	AccessLevel string `json:"access_level,omitempty" yaml:"access_level,omitempty"`
	Email       string `json:"email,omitempty" yaml:"email,omitempty"`
}

// https://docs.gitlab.com/ee/development/permissions.html#members
const (
	None       gitlab.AccessLevelValue = 0
	Minimal    gitlab.AccessLevelValue = 5
	Guest      gitlab.AccessLevelValue = 10
	Reporter   gitlab.AccessLevelValue = 20
	Developer  gitlab.AccessLevelValue = 30
	Maintainer gitlab.AccessLevelValue = 40
	Owner      gitlab.AccessLevelValue = 50
)

func getAccessLevel(accessLevel string) *gitlab.AccessLevelValue {
	var v gitlab.AccessLevelValue
	switch accessLevel {
	case "None":
		v = None
	case "Minimal":
		v = Minimal
	case "Guest":
		v = Guest
	case "Reporter":
		v = Reporter
	case "Maintainer":
		v = Maintainer
	case "Owner":
		v = Owner
	default:
		v = Developer
	}
	return gitlab.AccessLevel(v)
}

func (pc ProjectCtx) createInvites() {
	for _, invite := range pc.Project.Invites {
		_, _, err := pc.Client.Invites.ProjectInvites(pc.ProjectID, &gitlab.InvitesOptions{
			Email:       &invite.Email,
			AccessLevel: getAccessLevel(invite.AccessLevel),
		})
		if err != nil {
			fmt.Printf("[ERROR] Invite for `%s` could not be sent for project `%s` -- %s\n", invite.Email, pc.ProjectID, err)
		} else {
			fmt.Printf("[INFO] Invite for `%s` sent for project `%s`.\n", invite.Email, pc.ProjectID)
		}
	}
}

func replaceInvites(pc *ProjectCtx, api API) {
}
