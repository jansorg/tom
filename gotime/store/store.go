package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/satori/uuid"
)

func nextID() string {
	return uuid.NewV4().String()
}

type Project struct {
	Id        string `json:"id"`
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
}

type Frame struct {
	Id        string     `json:"id"`
	ProjectId string     `json:"project"`
	Start     *time.Time `json:"start,omitempty"`
	End       *time.Time `json:"end,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

func (f *Frame) IsStopped() bool {
	return f.End != nil && !f.End.IsZero()
}

func (f *Frame) IsActive() bool {
	return f.End == nil || f.End.IsZero()
}

func (f *Frame) Stop() {
	now := time.Now()
	f.End = &now
	f.Updated = &now
}

func (f *Frame) Duration() time.Duration {
	if f.IsStopped() {
		return f.End.Sub(*f.Start)
	}
	return time.Duration(0)
}

func (f *Frame) IsBefore(other *Frame) bool {
	return f.Start != nil && other.Start != nil && f.Start.Before(*other.Start)
}

func (f *Frame) IsAfter(other *Frame) bool {
	return !f.IsBefore(other) && f.Start != nil && other.Start != nil && f.Start.After(*other.Start)
}

func NewStartedFrame(project Project) Frame {
	now := time.Now()
	return Frame{
		Id:        nextID(),
		ProjectId: project.Id,
		Start:     &now,
		Updated:   &now,
	}
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Store interface {
	Reset(projects, tags, frames bool) error

	Projects() []Project
	AddProject(project Project) (Project, error)
	UpdateProject(project Project) (Project, error)
	RemoveProject(id string) error
	FindProject(id string) (Project, error)
	FindProjects(func(Project) bool) []Project

	Tags() []Tag
	AddTag(tag Tag) (Tag, error)
	UpdateTag(tag Tag) (Tag, error)
	RemoveTag(id string) error
	FindTag(id string) (Tag, error)
	FindTags(func(Tag) bool) []Tag

	Frames() []Frame
	AddFrame(frame Frame) (Frame, error)
	UpdateFrame(frame Frame) (Frame, error)
	RemoveFrame(id string) error
	FindFrame(id string) (Frame, error)
	FindFrames(func(Frame) bool) []Frame
}

func NewStore(dir string) (Store, error) {
	store := &DataStore{
		path:        dir,
		projectFile: filepath.Join(dir, "projects.json"),
		tagFile:     filepath.Join(dir, "tags.json"),
		frameFile:   filepath.Join(dir, "frames.json"),
	}

	if err := store.load(); err != nil {
		// 	return nil, err
	}
	return store, nil
}

type DataStore struct {
	path string

	tagFile     string
	frameFile   string
	projectFile string

	mu       sync.Mutex
	projects []Project
	tags     []Tag
	frames   []Frame
}

func (d *DataStore) load() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.loadLocked()
}

func (d *DataStore) loadLocked() error {
	var data []byte
	var err error

	if data, err = ioutil.ReadFile(d.projectFile); err != nil {
		return err
	}
	if err = json.Unmarshal(data, &d.projects); err != nil {
		return err
	}

	if data, err = ioutil.ReadFile(d.tagFile); err != nil {
		return err
	}
	if err = json.Unmarshal(data, &d.tags); err != nil {
		return err
	}

	if data, err = ioutil.ReadFile(d.frameFile); err != nil {
		return err
	}
	if err = json.Unmarshal(data, &d.frames); err != nil {
		return err
	}

	return nil
}

func (d *DataStore) save() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.saveLocked()
}

func (d *DataStore) saveLocked() error {
	var data []byte
	var err error

	if data, err = json.Marshal(d.projects); err != nil {
		return err
	}
	if err := ioutil.WriteFile(d.projectFile, data, 0600); err != nil {
		return err
	}

	if data, err = json.Marshal(d.tags); err != nil {
		return err
	}
	if err := ioutil.WriteFile(d.tagFile, data, 0600); err != nil {
		return err
	}

	if data, err = json.Marshal(d.frames); err != nil {
		return err
	}
	if err := ioutil.WriteFile(d.frameFile, data, 0600); err != nil {
		return err
	}

	return nil
}

func (d *DataStore) Reset(projects, tags, frames bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if projects {
		d.projects = []Project{}
	}
	if tags {
		d.tags = []Tag{}
	}
	if frames {
		d.frames = []Frame{}
	}

	return d.saveLocked()
}

func (d *DataStore) Projects() []Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.projects
}

func (d *DataStore) AddProject(project Project) (Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	project.Id = nextID()
	d.projects = append(d.projects, project)
	return project, d.saveLocked()
}

func (d *DataStore) UpdateProject(project Project) (Project, error) {
	panic("implement me")
}

func (d *DataStore) RemoveProject(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, p := range d.projects {
		if p.Id == id {
			d.projects = append(d.projects[:i], d.projects[i+1:]...)
			return d.saveLocked()
		}
	}

	return fmt.Errorf("project %s not found", id)
}

func (d *DataStore) FindProject(id string) (Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, p := range d.projects {
		if p.Id == id {
			return p, nil
		}
	}
	return Project{}, fmt.Errorf("project %s not found", id)
}

func (d *DataStore) FindProjects(filter func(Project) bool) []Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []Project
	for _, p := range d.projects {
		if filter(p) {
			result = append(result, p)
		}
	}
	return result
}

func (d *DataStore) Tags() []Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.tags
}

func (d *DataStore) AddTag(tag Tag) (Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tag.Id = nextID()
	d.tags = append(d.tags, tag)
	return Tag{}, d.saveLocked()
}

func (d *DataStore) UpdateTag(tag Tag) (Tag, error) {
	panic("implement me")
}

func (d *DataStore) RemoveTag(id string) error {
	panic("implement me")
}

func (d *DataStore) FindTag(id string) (Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, tag := range d.tags {
		if tag.Id == id {
			return tag, nil
		}
	}
	return Tag{}, fmt.Errorf("tag %s not found", id)
}

func (d *DataStore) FindTags(filter func(Tag) bool) []Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []Tag
	for _, tag := range d.tags {
		if filter(tag) {
			result = append(result, tag)
		}
	}
	return result
}

func (d *DataStore) Frames() []Frame {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.frames
}

func (d *DataStore) AddFrame(frame Frame) (Frame, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	frame.Id = nextID()
	d.frames = append(d.frames, frame)
	return frame, d.saveLocked()
}

func (d *DataStore) UpdateFrame(frame Frame) (Frame, error) {
	if frame.Id == "" {
		return Frame{}, fmt.Errorf("id of frame undefined")
	}

	if err := d.RemoveFrame(frame.Id); err != nil {
		return Frame{}, err
	}
	return d.AddFrame(frame)
}

func (d *DataStore) RemoveFrame(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, frame := range d.frames {
		if frame.Id == id {
			d.frames = append(d.frames[:i], d.frames[i+1:]...)
			return d.saveLocked()
		}
	}
	return fmt.Errorf("frame %s not found", id)
}

func (d *DataStore) FindFrame(id string) (Frame, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, frame := range d.frames {
		if frame.Id == id {
			return frame, nil
		}
	}
	return Frame{}, fmt.Errorf("frame %s not found", id)
}

func (d *DataStore) FindFrames(filter func(Frame) bool) []Frame {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []Frame
	for _, frame := range d.frames {
		if filter(frame) {
			result = append(result, frame)
		}
	}
	return result
}
