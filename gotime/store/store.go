package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/satori/uuid"
)

func nextID() string {
	return uuid.NewV4().String()
}

type Store interface {
	DirPath() string
	StartBatch()
	StopBatch()

	Reset(projects, tags, frames bool) error

	Projects() []*Project
	ProjectByID(id string) (*Project, error)
	ProjectIsChild(parentID, id string) bool
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
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", dir)
	}

	store := &DataStore{
		path:        dir,
		projectFile: filepath.Join(dir, "projects.json"),
		tagFile:     filepath.Join(dir, "tags.json"),
		frameFile:   filepath.Join(dir, "frames.json"),
	}

	if err := store.loadLocked(); err != nil {
		return nil, err
	}
	return store, nil
}

type DataStore struct {
	path      string
	batchMode int32

	tagFile     string
	frameFile   string
	projectFile string

	mu          sync.RWMutex
	projectsMap map[string]*Project
	projects    []*Project
	tags        []*Tag
	frames      []*Frame
}

func (d *DataStore) DirPath() string {
	return filepath.Dir(d.projectFile)
}

func (d *DataStore) StartBatch() {
	swapped := atomic.CompareAndSwapInt32(&d.batchMode, 0, 1)
	if !swapped {
		log.Fatal(fmt.Errorf("StartBatch() in batch mode"))
	}
}

func (d *DataStore) StopBatch() {
	swapped := atomic.CompareAndSwapInt32(&d.batchMode, 1, 0)
	if !swapped {
		log.Fatal(fmt.Errorf("StopBatch() called without prior StartBatch()"))
	}

	_ = d.save()
}

func (d *DataStore) sortProjects() {
	sort.SliceStable(d.projects, func(i, j int) bool {
		return strings.Compare(d.projects[i].FullName, d.projects[j].FullName) < 0
	})
}

func (d *DataStore) sortTags() {
	sort.SliceStable(d.tags, func(i, j int) bool {
		return strings.Compare(d.tags[i].Name, d.tags[j].Name) < 0
	})
}

func (d *DataStore) sortFrames() {
	sort.SliceStable(d.frames, func(i, j int) bool {
		return d.frames[i].IsBefore(d.frames[j])
	})
}

func (d *DataStore) load() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.loadLocked()
}

func (d *DataStore) loadLocked() error {
	var data []byte
	var err error

	if fileExists(d.projectFile) {
		if data, err = ioutil.ReadFile(d.projectFile); err != nil {
			return err
		}
		if err = json.Unmarshal(data, &d.projects); err != nil {
			return err
		}
	}

	if fileExists(d.tagFile) {
		if data, err = ioutil.ReadFile(d.tagFile); err != nil {
			return err
		}
		if err = json.Unmarshal(data, &d.tags); err != nil {
			return err
		}
	}

	if fileExists(d.frameFile) {
		if data, err = ioutil.ReadFile(d.frameFile); err != nil {
			return err
		}
		if err = json.Unmarshal(data, &d.frames); err != nil {
			return err
		}
	}

	// update internal data
	d.updateProjectsMapping()
	for _, p := range d.projects {
		d.updateProjectInternals(p)
	}
	d.sortProjects()
	d.sortTags()
	d.sortFrames()

	return nil
}

func (d *DataStore) save() error {
	if atomic.LoadInt32(&d.batchMode) == 1 {
		return nil;
	}

	fmt.Println("saving...")
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.saveLocked()
}

func (d *DataStore) saveLocked() error {
	d.sortProjects()
	d.sortTags()
	d.sortFrames()

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

func (d *DataStore) ProjectByID(id string) (*Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	p, ok := d.projectsMap[id]
	if !ok {
		return nil, fmt.Errorf("no project found for %s", id)
	}
	return p, nil
}

func (d *DataStore) AddProject(project Project) (*Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	project.ID = nextID()
	d.updateProjectInternals(&project)
	d.projects = append(d.projects, &project)
	d.updateProjectsMapping()

	return &project, d.saveLocked()
}

func (d *DataStore) updateProjectInternals(p *Project) {
	p.FullName = p.Name
	if p.ParentID == "" {
		return
	}

	parents := []string{p.Name}

	id := p.ParentID
	for id != "" {
		parent, err := d.findFirstProjectLocked(func(current *Project) bool {
			return current.ID == id
		})

		if err != nil {
			log.Fatal(fmt.Errorf("unable to find project %s", id))
		}

		id = parent.ParentID
		parents = append([]string{parent.Name}, parents...)
	}

	p.FullName = strings.Join(parents, "/")
}

func (d *DataStore) UpdateProject(project Project) (*Project, error) {
	if err := d.RemoveProject(project.ID); err != nil {
		return nil, err
	}
	// fixme keep id!
	return d.AddProject(project)
}

func (d *DataStore) RemoveProject(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, p := range d.projects {
		if p.ID == id {
			d.projects = append(d.projects[:i], d.projects[i+1:]...)
			d.updateProjectsMapping()
			return d.saveLocked()
		}
	}

	return fmt.Errorf("project %s not found", id)
}

func (d *DataStore) FindFirstProject(filter func(*Project) bool) (*Project, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.findFirstProjectLocked(filter)
}

func (d *DataStore) findFirstProjectLocked(filter func(*Project) bool) (*Project, error) {
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

func (d *DataStore) ProjectIsChild(parentID, id string) bool {
	if parentID == id {
		return true
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	for id != "" {
		if id == parentID {
			return true
		}

		project, ok := d.projectsMap[id]
		if !ok {
			return false
		}
		id = project.ParentID
	}
	return false
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
	return &tag, d.saveLocked()
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
	d.mu.RLock()
	defer d.mu.RUnlock()

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

	frame.ID = nextID()
	d.frames = append(d.frames, &frame)
	return &frame, d.saveLocked()
}

func (d *DataStore) UpdateFrame(frame Frame) (*Frame, error) {
	if frame.ID == "" {
		return nil, fmt.Errorf("id of frame undefined")
	}

	if err := d.RemoveFrame(frame.ID); err != nil {
		return nil, err
	}
	return d.AddFrame(frame)
}

func (d *DataStore) RemoveFrame(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, frame := range d.frames {
		if frame.ID == id {
			d.frames = append(d.frames[:i], d.frames[i+1:]...)
			return d.saveLocked()
		}
	}
	return fmt.Errorf("frame %s not found", id)
}

func (d *DataStore) FindFirstFrame(filter func(*Frame) bool) (*Frame, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, frame := range d.frames {
		if filter(frame) {
			return frame, nil
		}
	}
	return nil, fmt.Errorf("no matching frame found")
}

func (d *DataStore) FindFrames(filter func(*Frame) bool) []*Frame {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var result []*Frame
	for _, frame := range d.frames {
		if filter(frame) {
			result = append(result, frame)
		}
	}
	return result
}

func (d *DataStore) updateProjectsMapping() {
	d.projectsMap = map[string]*Project{}
	for _, p := range d.projects {
		d.projectsMap[p.ID] = p
	}
}
