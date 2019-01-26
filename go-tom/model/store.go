package model

type PropertyHolder interface {
	GetProperties() map[string]string
}

type Store interface {
	DirPath() string
	StartBatch()
	StopBatch()

	Reset(projects, tags, frames bool) (int, int, int, error)

	Projects() ProjectList
	ProjectByID(id string) (*Project, error)
	ProjectIsSameOrChild(parentID, id string) bool
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

	Frames() FrameList
	AddFrame(frame Frame) (*Frame, error)
	UpdateFrame(frame Frame) (*Frame, error)
	RemoveFrame(id string) error
	FindFirstFrame(func(*Frame) bool) (*Frame, error)
	FindFrames(func(*Frame) (bool, error)) ([]*Frame, error)
}
