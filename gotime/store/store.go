package store

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/satori/uuid"
)

type Project struct {
	Id        string `json:"id"`
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
}

type Frame struct {
	Id      string    `json:"id"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Updated time.Time `json:"updated,omitempty"`
	Notes   string    `json:"notes"`
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Store interface {
	Reset() error

	Projects() []Project
	AddProject(project Project) error
	UpdateProject(project Project) error
	RemoveProject(id string) error

	Tags() []Tag
	AddTag(tag Tag) error
	UpdateTag(tag Tag) error
	RemoveTag(id string) error

	Frames() []Frame
	AddFrame(frame Frame) error
	UpdateFrame(frame Frame) error
	RemoveFrame(id string) error
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

func (d *DataStore) nextID() string {
	return uuid.NewV4().String()
}

func (d *DataStore) Reset() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.projects = []Project{}
	d.tags = []Tag{}
	d.frames = []Frame{}

	return d.saveLocked()
}

func (d *DataStore) Projects() []Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.projects
}

func (d *DataStore) AddProject(project Project) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	project.Id = d.nextID()
	d.projects = append(d.projects, project)
	return d.saveLocked()
}

func (d *DataStore) UpdateProject(project Project) error {
	panic("implement me")
}

func (d *DataStore) RemoveProject(id string) error {
	panic("implement me")
}

func (d *DataStore) Tags() []Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.tags
}

func (d *DataStore) AddTag(tag Tag) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tag.Id = d.nextID()
	d.tags = append(d.tags, tag)
	return d.saveLocked()
}

func (d *DataStore) UpdateTag(tag Tag) error {
	panic("implement me")
}

func (d *DataStore) RemoveTag(id string) error {
	panic("implement me")
}

func (d *DataStore) Frames() []Frame {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.frames
}

func (d *DataStore) AddFrame(frame Frame) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	frame.Id = d.nextID()
	d.frames = append(d.frames, frame)
	return d.saveLocked()
}

func (d *DataStore) UpdateFrame(frame Frame) error {
	panic("implement me")
}

func (d *DataStore) RemoveFrame(id string) error {
	panic("implement me")
}
