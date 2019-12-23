package activity

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

var ProjectNotFoundErr = fmt.Errorf("project not found")
var NoteRequiredErr = fmt.Errorf("note required")

type Control struct {
	ctx                   *context.TomContext
	createMissingProjects bool
	createMissingTags     bool
	allowMultipleActives  bool
	startStopTime         time.Time
}

func NewActivityControl(ctx *context.TomContext, createMissing bool, allowMultipleActives bool, startStopTime time.Time) *Control {
	return &Control{
		ctx:                   ctx,
		createMissingProjects: createMissing,
		createMissingTags:     createMissing,
		allowMultipleActives:  allowMultipleActives,
		startStopTime:         startStopTime,
	}
}

func (a *Control) Start(projectNameOrID string, notes string, tags []*model.Tag) (*model.Frame, error) {
	project, err := a.ctx.Query.ProjectByID(projectNameOrID)
	if err != nil {
		project, err = a.ctx.Query.ProjectByFullName(strings.Split(projectNameOrID, "/"))
	}

	if project == nil && a.createMissingProjects == false {
		return nil, ProjectNotFoundErr
	} else if project == nil {
		if project, err = a.ctx.Store.AddProject(model.Project{Name: projectNameOrID}); err != nil {
			return nil, err
		}
	}

	frame := model.NewStartedFrame(project)
	frame.Notes = notes
	frame.Start = &a.startStopTime
	frame.AddTags(tags...)
	return a.ctx.Store.AddFrame(frame)
}

func (a *Control) StopNewest(notes string, tags []*model.Tag) (*model.Frame, error) {
	var frames []*model.Frame
	var err error
	if frames, err = a.stopActivities(false, notes, tags); err != nil {
		return nil, err
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no running activity found")
	}
	return frames[0], nil
}

func (a *Control) StopAll(notes string, tags []*model.Tag) ([]*model.Frame, error) {
	return a.stopActivities(true, notes, tags)
}

func (a *Control) stopActivities(all bool, notes string, tags []*model.Tag) ([]*model.Frame, error) {
	actives := a.ctx.Query.ActiveFrames()

	if !all && len(actives) > 0 {
		sort.SliceStable(actives, func(i, j int) bool {
			return actives[i].Start.After(*actives[j].Start)
		})
		actives = actives[:1]
	}

	// validate, that a note is defined when required
	if notes == "" {
		for _, frame := range actives {
			noteRequired, err := a.ctx.Query.IsNoteRequired(frame.ProjectId)
			if err == nil && noteRequired != nil && *noteRequired == true {
				return nil, NoteRequiredErr
			}
		}
	}

	// update data
	for _, frame := range actives {
		frame.StopAt(a.startStopTime)

		if notes != "" {
			frame.Notes = notes
		}

		frame.AddTags(tags...)

		if _, err := a.ctx.Store.UpdateFrame(*frame); err != nil {
			return nil, err
		}
	}

	return actives, nil
}
