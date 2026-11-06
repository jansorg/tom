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
	"time"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

var ErrTagNotFound = fmt.Errorf("tag not found")

func NewStore(dir string, backupDir string, maxBackups int) (model.Store, error) {
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", dir)
	}

	store := &DataStore{
		path:         dir,
		backupPath:   backupDir,
		maxBackups:   maxBackups,
		ProjectFile:  filepath.Join(dir, "projects.json"),
		TagFile:      filepath.Join(dir, "tags.json"),
		FrameFile:    filepath.Join(dir, "frames.json"),
		PropertyFile: filepath.Join(dir, "properties.json"),
	}

	if err := store.loadLocked(); err != nil {
		return nil, err
	}
	return store, nil
}

type DataStore struct {
	path       string
	backupPath string
	batchMode  int32

	maxBackups int

	ProjectFile  string
	TagFile      string
	FrameFile    string
	PropertyFile string

	mu          sync.RWMutex
	projectsMap map[string]*model.Project
	projects    []*model.Project
	tags        []*model.Tag
	frames      []*model.Frame
}

func (d *DataStore) DirPath() string {
	return d.path
}

func (d *DataStore) BackupDirPath() string {
	return d.backupPath
}

func (d *DataStore) MaxBackups() int {
	return d.maxBackups
}

func (d *DataStore) StartBatch() {
	swapped := atomic.CompareAndSwapInt32(&d.batchMode, 0, 1)
	if !swapped {
		log.Fatal(fmt.Errorf("StartBatch() in batch mode"))
	}
}

func (d *DataStore) StopBatch() error {
	swapped := atomic.CompareAndSwapInt32(&d.batchMode, 1, 0)
	if !swapped {
		log.Fatal(fmt.Errorf("StopBatch() called without prior StartBatch()"))
	}

	return d.save()
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
		if err := d.updateProjectInternals(p); err != nil {
			// don't make all commands fail because of inconsistent data, the user has to be able to fix it
			_, _ = fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		}
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

	if err := d.backupLocked(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: unable to create backup: %v\n", err)
	}

	d.sortProjects()
	d.sortTags()
	d.sortFrames()

	var data []byte
	var err error

	// projects
	if data, err = json.Marshal(d.projects); err != nil {
		return err
	}
	if err := writeFileAtomic(d.ProjectFile, data, 0600); err != nil {
		return err
	}

	// tags
	if data, err = json.Marshal(d.tags); err != nil {
		return err
	}
	if err := writeFileAtomic(d.TagFile, data, 0600); err != nil {
		return err
	}

	// frames
	if data, err = json.Marshal(d.frames); err != nil {
		return err
	}
	if err := writeFileAtomic(d.FrameFile, data, 0600); err != nil {
		return err
	}

	return nil
}

// writeFileAtomic writes data into a temporary file next to path and renames it to path,
// so that a crash or a full disk never leaves a truncated data file behind.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp, err := ioutil.TempFile(filepath.Dir(path), filepath.Base(path)+".tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	_, err = tmp.Write(data)
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Chmod(tmpName, perm)
	}
	if err == nil {
		err = os.Rename(tmpName, path)
	}
	if err != nil {
		_ = os.Remove(tmpName)
	}
	return err
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

func (d *DataStore) Projects() model.ProjectList {
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

	if err := d.validateParentLocked("", project.ParentID); err != nil {
		return nil, err
	}

	project.ID = model.NextID()
	d.projects = append(d.projects, &project)
	d.updateProjectsMapping()
	_ = d.updateProjectInternals(&project)

	return &project, d.saveLocked()
}

// updateProjectInternals updates the cached full name of p.
// It returns an error if a parent is missing or if the parents form a cycle.
// In both cases the full name contains the parents which could be resolved.
func (d *DataStore) updateProjectInternals(p *model.Project) error {
	p.Store = d

	names := []string{p.Name}
	visited := map[string]bool{p.ID: true}

	id := p.ParentID
	for id != "" {
		if visited[id] {
			p.FullName = names
			return fmt.Errorf("project %s is its own ancestor", p.ID)
		}
		visited[id] = true

		parent, ok := d.projectsMap[id]
		if !ok {
			p.FullName = names
			return fmt.Errorf("parent project %s of project %s not found", id, p.ID)
		}

		id = parent.ParentID
		names = append([]string{parent.Name}, names...)
	}

	p.FullName = names
	return nil
}

// validateParentLocked returns an error if parentID does not exist
// or if setting it as parent of the project with ID projectID would create a cycle.
func (d *DataStore) validateParentLocked(projectID, parentID string) error {
	visited := map[string]bool{}
	for id := parentID; id != ""; {
		if projectID != "" && id == projectID {
			return fmt.Errorf("a project can't be moved into itself or into one of its subprojects")
		}
		if visited[id] {
			return fmt.Errorf("the parents of project %s form a cycle", parentID)
		}
		visited[id] = true

		parent, ok := d.projectsMap[id]
		if !ok {
			return fmt.Errorf("parent project %s not found", id)
		}
		id = parent.ParentID
	}
	return nil
}

func (d *DataStore) UpdateProject(project model.Project) (*model.Project, error) {
	if err := project.Validate(); err != nil {
		return nil, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	existing, ok := d.projectsMap[project.ID]
	if !ok {
		return nil, fmt.Errorf("no project found for %s", project.ID)
	}
	if err := d.validateParentLocked(project.ID, project.ParentID); err != nil {
		return nil, err
	}

	*existing = project

	// the full names of all subprojects depend on the name and parent of this project
	for _, p := range d.projects {
		if err := d.updateProjectInternals(p); err != nil {
			return nil, err
		}
	}
	return existing, d.saveLocked()
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

func (d *DataStore) ProjectIsSameOrChild(parentID, id string) bool {
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

	d.mu.Lock()
	defer d.mu.Unlock()

	*existing = tag
	return existing, d.saveLocked()
}

func (d *DataStore) RemoveTag(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, t := range d.tags {
		if t.ID == id {
			d.tags = append(d.tags[:i], d.tags[i+1:]...)
			return d.saveLocked()
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

func (d *DataStore) Frames() model.FrameList {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.frames
}

func (d *DataStore) AddFrame(frame model.Frame) (*model.Frame, error) {
	if err := frame.Validate(false); err != nil {
		return nil, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	frame.ID = model.NextID()
	d.frames = append(d.frames, &frame)
	return &frame, d.saveLocked()
}

func (d *DataStore) UpdateFrame(frame model.Frame) (*model.Frame, error) {
	if err := frame.Validate(true); err != nil {
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

func (d *DataStore) FindFrames(filter func(*model.Frame) (bool, error)) ([]*model.Frame, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var result []*model.Frame
	for _, frame := range d.frames {
		if ok, err := filter(frame); err != nil {
			return nil, err
		} else if ok {
			result = append(result, frame)
		}
	}
	return result, nil
}

func (d *DataStore) updateProjectsMapping() {
	d.projectsMap = map[string]*model.Project{}
	for _, p := range d.projects {
		d.projectsMap[p.ID] = p
	}
}

func (d *DataStore) backupLocked() error {
	if d.backupPath == "" || d.Empty() {
		return nil
	}

	var sourceFiles []string
	for _, sourceFile := range []string{d.ProjectFile, d.FrameFile, d.TagFile} {
		if _, err := os.Stat(sourceFile); err == nil {
			sourceFiles = append(sourceFiles, sourceFile)
		}
	}
	if len(sourceFiles) == 0 {
		// nothing saved yet
		return nil
	}

	// UTC and no colons: valid on Windows and sorts chronologically
	targetDir := filepath.Join(d.backupPath, fmt.Sprintf("tom-%s", time.Now().UTC().Format("20060102T150405.000000000Z")))
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("backup directory already exists: %s", targetDir)
	}
	err := os.MkdirAll(targetDir, 0700)
	if err != nil {
		return err
	}

	for _, sourceFile := range sourceFiles {
		if err = util.CopyFile(sourceFile, filepath.Join(targetDir, filepath.Base(sourceFile)), false); err != nil {
			// fixme cleanup?
			return err
		}
	}

	// remove backups > max backups
	if d.maxBackups > 0 {
		infos, err := ioutil.ReadDir(d.backupPath)
		if err == nil && len(infos) > d.maxBackups {
			// filter and validate
			var matching []string
			for _, info := range infos {
				baseName := info.Name()
				if info.IsDir() && strings.HasPrefix(baseName, "tom-") {
					matching = append(matching, filepath.Join(d.backupPath, baseName))
				}
			}

			// remove oldest backups first
			if len(matching) > d.maxBackups {
				// sort newest to oldest
				sort.Slice(matching, func(i, j int) bool {
					return strings.Compare(matching[i], matching[j]) >= 0
				})

				for _, dirName := range matching[d.maxBackups:] {
					_ = os.RemoveAll(dirName)
				}
			}
		}
	}

	return nil
}

func (d *DataStore) Empty() bool {
	return len(d.projects) == 0 && len(d.tags) == 0 && len(d.frames) == 0
}
