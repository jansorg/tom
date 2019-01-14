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

	"github.com/jansorg/tom/go-tom/model"
)

var ErrTagNotFound = fmt.Errorf("tag not found")

func NewStore(dir string) (model.Store, error) {
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", dir)
	}

	store := &DataStore{
		path:        dir,
		ProjectFile: filepath.Join(dir, "projects.json"),
		TagFile:     filepath.Join(dir, "tags.json"),
		FrameFile:   filepath.Join(dir, "frames.json"),
	}

	if err := store.loadLocked(); err != nil {
		return nil, err
	}
	return store, nil
}

type DataStore struct {
	path      string
	batchMode int32

	ProjectFile string
	TagFile     string
	FrameFile   string

	mu          sync.RWMutex
	projectsMap map[string]*model.Project
	projects    []*model.Project
	tags        []*model.Tag
	frames      []*model.Frame
}

func (d *DataStore) DirPath() string {
	return filepath.Dir(d.ProjectFile)
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
		return strings.Compare(d.projects[i].GetFullName("/"), d.projects[j].GetFullName("/")) < 0
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

	if fileExists(d.ProjectFile) {
		if data, err = ioutil.ReadFile(d.ProjectFile); err != nil {
			return err
		}
		if err = json.Unmarshal(data, &d.projects); err != nil {
			return err
		}
	}

	if fileExists(d.TagFile) {
		if data, err = ioutil.ReadFile(d.TagFile); err != nil {
			return err
		}
		if err = json.Unmarshal(data, &d.tags); err != nil {
			return err
		}
	}

	if fileExists(d.FrameFile) {
		if data, err = ioutil.ReadFile(d.FrameFile); err != nil {
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
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.saveLocked()
}

func (d *DataStore) saveLocked() error {
	if atomic.LoadInt32(&d.batchMode) == 1 {
		return nil
	}

	d.sortProjects()
	d.sortTags()
	d.sortFrames()

	var data []byte
	var err error

	if data, err = json.Marshal(d.projects); err != nil {
		return err
	}
	if err := ioutil.WriteFile(d.ProjectFile, data, 0600); err != nil {
		return err
	}

	if data, err = json.Marshal(d.tags); err != nil {
		return err
	}
	if err := ioutil.WriteFile(d.TagFile, data, 0600); err != nil {
		return err
	}

	if data, err = json.Marshal(d.frames); err != nil {
		return err
	}
	if err := ioutil.WriteFile(d.FrameFile, data, 0600); err != nil {
		return err
	}

	return nil
}

func (d *DataStore) Reset(projects, tags, frames bool) (int, int, int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var projectCount, tagCount, frameCount int

	if projects {
		projectCount = len(d.projects)
		d.projects = []*model.Project{}
	}
	if tags {
		tagCount = len(d.tags)
		d.tags = []*model.Tag{}
	}
	if frames {
		frameCount = len(d.frames)
		d.frames = []*model.Frame{}
	}

	return projectCount, tagCount, frameCount, d.saveLocked()
}

func (d *DataStore) Projects() []*model.Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.projects
}

func (d *DataStore) ProjectByID(id string) (*model.Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	p, ok := d.projectsMap[id]
	if !ok {
		return nil, fmt.Errorf("no project found for %s", id)
	}
	return p, nil
}

func (d *DataStore) AddProject(project model.Project) (*model.Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	project.ID = model.NextID()
	d.updateProjectInternals(&project)
	d.projects = append(d.projects, &project)
	d.updateProjectsMapping()

	return &project, d.saveLocked()
}

func (d *DataStore) updateProjectInternals(p *model.Project) {
	p.Store = d

	if p.Properties == nil {
		p.Properties = make(map[string]string)
	}

	p.FullName = []string{p.Name}
	if p.ParentID == "" {
		return
	}

	parents := []string{p.Name}

	id := p.ParentID
	for id != "" {
		parent, err := d.findFirstProjectLocked(func(current *model.Project) bool {
			return current.ID == id
		})

		if err != nil {
			log.Fatal(fmt.Errorf("unable to find project %s", id))
		}

		id = parent.ParentID
		parents = append([]string{parent.Name}, parents...)
	}

	p.FullName = parents
}

func (d *DataStore) UpdateProject(project model.Project) (*model.Project, error) {
	if err := project.Validate(); err != nil {
		return nil, err
	}

	existing, err := d.ProjectByID(project.ID)
	if err != nil {
		return nil, err
	}
	*existing = project
	d.updateProjectInternals(existing)
	return existing, d.save()
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

func (d *DataStore) FindFirstProject(filter func(*model.Project) bool) (*model.Project, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.findFirstProjectLocked(filter)
}

func (d *DataStore) findFirstProjectLocked(filter func(*model.Project) bool) (*model.Project, error) {
	for _, p := range d.projects {
		if filter(p) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no matching project found")
}

func (d *DataStore) FindProjects(filter func(*model.Project) bool) []*model.Project {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []*model.Project
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

func (d *DataStore) Tags() []*model.Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.tags
}

func (d *DataStore) AddTag(tag model.Tag) (*model.Tag, error) {
	if err := tag.Validate(); err != nil {
		return nil, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	tag.ID = model.NextID()
	d.tags = append(d.tags, &tag)
	return &tag, d.saveLocked()
}

func (d *DataStore) UpdateTag(tag model.Tag) (*model.Tag, error) {
	if err := tag.Validate(); err != nil {
		return nil, err
	}

	existing, err := d.FindTag(tag.ID)
	if err != nil {
		return nil, err
	}

	*existing = tag
	return existing, d.saveLocked()
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

func (d *DataStore) FindTag(id string) (*model.Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, tag := range d.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return nil, fmt.Errorf("tag %s not found", id)
}

func (d *DataStore) FindFirstTag(filter func(*model.Tag) bool) (*model.Tag, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, tag := range d.tags {
		if filter(tag) {
			return tag, nil
		}
	}
	return nil, ErrTagNotFound
}

func (d *DataStore) FindTags(filter func(*model.Tag) bool) []*model.Tag {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []*model.Tag
	for _, tag := range d.tags {
		if filter(tag) {
			result = append(result, tag)
		}
	}
	return result
}

func (d *DataStore) Frames() []*model.Frame {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.frames
}

func (d *DataStore) AddFrame(frame model.Frame) (*model.Frame, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	frame.ID = model.NextID()
	d.frames = append(d.frames, &frame)
	return &frame, d.saveLocked()
}

func (d *DataStore) UpdateFrame(frame model.Frame) (*model.Frame, error) {
	if err := frame.Validate(); err != nil {
		return nil, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, f := range d.frames {
		if f.ID == frame.ID {
			*f = frame
			return f, d.saveLocked()
		}
	}
	return nil, fmt.Errorf("no frame with ID %s found", frame.ID)
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

func (d *DataStore) FindFirstFrame(filter func(*model.Frame) bool) (*model.Frame, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, frame := range d.frames {
		if filter(frame) {
			return frame, nil
		}
	}
	return nil, fmt.Errorf("no matching frame found")
}

func (d *DataStore) FindFrames(filter func(*model.Frame) bool) []*model.Frame {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var result []*model.Frame
	for _, frame := range d.frames {
		if filter(frame) {
			result = append(result, frame)
		}
	}
	return result
}

func (d *DataStore) updateProjectsMapping() {
	d.projectsMap = map[string]*model.Project{}
	for _, p := range d.projects {
		d.projectsMap[p.ID] = p
	}
}
