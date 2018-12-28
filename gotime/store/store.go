package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/satori/uuid"
)

func nextID() string {
	return uuid.NewV4().String()
}

type Store interface {
	Reset(projects, tags, frames bool) error

	Projects() []*Project
	AddProject(project Project) (*Project, error)
	UpdateProject(project Project) (*Project, error)
	RemoveProject(id string) error
	FindFirstProject(func(*Project) bool) (*Project, error)
	FindProjects(func(*Project) bool) []*Project

	Tags() []*Tag
	AddTag(tag Tag) (*Tag, error)
	UpdateTag(tag Tag) (*Tag, error)
	RemoveTag(id string) error
	FindFirstTag(func(*Tag) bool) (*Tag, error)
	FindTags(func(*Tag) bool) []*Tag

	Frames() []*Frame
	AddFrame(frame Frame) (*Frame, error)
	UpdateFrame(frame Frame) (*Frame, error)
	RemoveFrame(id string) error
	FindFirstFrame(func(*Frame) bool) (*Frame, error)
	FindFrames(func(*Frame) bool) []*Frame
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
	projects []*Project
	tags     []*Tag
	frames   []*Frame
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
		d.projects = []*Project{}
	}
	if tags {
		d.tags = []*Tag{}
	}
	if frames {
		d.frames = []*Frame{}
	}

	return d.saveLocked()
}

func (d *DataStore) Projects() []*Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.projects
}

func (d *DataStore) AddProject(project Project) (*Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	project.ID = nextID()
	d.projects = append(d.projects, &project)
	return &project, d.saveLocked()
}

func (d *DataStore) UpdateProject(project Project) (*Project, error) {
	if err := d.RemoveProject(project.ID); err != nil {
		return nil, err
	}

	return d.AddProject(project)
}

func (d *DataStore) RemoveProject(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, p := range d.projects {
		if p.ID == id {
			d.projects = append(d.projects[:i], d.projects[i+1:]...)
			return d.saveLocked()
		}
	}

	return fmt.Errorf("project %s not found", id)
}

func (d *DataStore) FindFirstProject(filter func(*Project) bool) (*Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, p := range d.projects {
		if filter(p) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no matching project found")
}

func (d *DataStore) FindProjects(filter func(*Project) bool) []*Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []*Project
	for _, p := range d.projects {
		if filter(p) {
			result = append(result, p)
		}
	}
	return result
}

func (d *DataStore) Tags() []*Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.tags
}

func (d *DataStore) AddTag(tag Tag) (*Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tag.ID = nextID()
	d.tags = append(d.tags, &tag)
	return nil, d.saveLocked()
}

func (d *DataStore) UpdateTag(tag Tag) (*Tag, error) {
	if err := d.RemoveTag(tag.ID); err != nil {
		return nil, err
	}

	return d.AddTag(tag)
}

func (d *DataStore) RemoveTag(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, t := range d.tags {
		if t.ID == id {
			d.tags = append(d.tags[:i], d.tags[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("tag %s not found", id)
}

func (d *DataStore) FindTag(id string) (*Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, tag := range d.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return nil, fmt.Errorf("tag %s not found", id)
}

func (d *DataStore) FindFirstTag(filter func(*Tag) bool) (*Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, tag := range d.tags {
		if filter(tag) {
			return tag, nil
		}
	}
	return nil, fmt.Errorf("no matching tag found")
}

func (d *DataStore) FindTags(filter func(*Tag) bool) []*Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []*Tag
	for _, tag := range d.tags {
		if filter(tag) {
			result = append(result, tag)
		}
	}
	return result
}

func (d *DataStore) Frames() []*Frame {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.frames
}

func (d *DataStore) AddFrame(frame Frame) (*Frame, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	frame.Id = nextID()
	d.frames = append(d.frames, &frame)
	return &frame, d.saveLocked()
}

func (d *DataStore) UpdateFrame(frame Frame) (*Frame, error) {
	if frame.Id == "" {
		return nil, fmt.Errorf("id of frame undefined")
	}

	if err := d.RemoveFrame(frame.Id); err != nil {
		return nil, err
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

func (d *DataStore) FindFirstFrame(filter func(*Frame) bool) (*Frame, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, frame := range d.frames {
		if filter(frame) {
			return frame, nil
		}
	}
	return nil, fmt.Errorf("no matching frame found")
}

func (d *DataStore) FindFrames(filter func(*Frame) bool) []*Frame {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []*Frame
	for _, frame := range d.frames {
		if filter(frame) {
			result = append(result, frame)
		}
	}
	return result
}
