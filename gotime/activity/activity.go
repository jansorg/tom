package activity

import (
	"fmt"
	"sort"
	"time"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

var ProjectNotFoundErr = fmt.Errorf("project not found")

type ActivityControl struct {
	ctx                   *context.GoTimeContext
	createMissingProjects bool
	createMissingTags     bool
	allowMultipleActives  bool
	startStopTime         time.Time
}

func NewActivityControl(ctx *context.GoTimeContext, createMissing bool, allowMultipleActives bool, startStopTime time.Time) *ActivityControl {
	return &ActivityControl{
		ctx:                   ctx,
		createMissingProjects: createMissing,
		createMissingTags:     createMissing,
		allowMultipleActives:  allowMultipleActives,
		startStopTime:         startStopTime,
	}
}

func (a *ActivityControl) Start(projectNameOrID string, notes string, tags []*store.Tag) (*store.Frame, error) {
	projects := a.ctx.Query.ProjectsByShortNameOrID(projectNameOrID)

	var err error
	var project *store.Project
	if len(projects) == 0 && a.createMissingProjects == false {
		return nil, ProjectNotFoundErr
	} else if len(projects) == 0 {
		if project, err = a.ctx.Store.AddProject(store.Project{Name: projectNameOrID}); err != nil {
			return nil, err
		}
	} else if len(projects) == 1 {
		project = projects[0]
	} else {
		return nil, fmt.Errorf("more than one project found for %s", projectNameOrID)
	}

	frame := store.NewStartedFrame(project)
	frame.Notes = notes
	frame.Start = &a.startStopTime
	frame.AddTags(tags...)
	return a.ctx.Store.AddFrame(frame)
}

func (a *ActivityControl) StopNewest(notes string, tags []*store.Tag) (*store.Frame, error) {
	var frames []*store.Frame
	var err error
	if frames, err = a.stopActivities(false, notes, tags); err != nil {
		return nil, err
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no running activity found")
	}
	return frames[0], nil
}

func (a *ActivityControl) StopAll(notes string, tags []*store.Tag) ([]*store.Frame, error) {
	return a.stopActivities(true, notes, tags)
}

func (a *ActivityControl) stopActivities(all bool, notes string, tags []*store.Tag) ([]*store.Frame, error) {
	actives := a.ctx.Query.ActiveFrames()

	if !all && len(actives) > 0 {
		sort.SliceStable(actives, func(i, j int) bool {
			return actives[i].Start.After(*actives[j].Start)
		})
		actives = actives[:1]
	}

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
