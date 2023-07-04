package validation

type PlaybookStruct struct {
	Name        string            `yaml:"name"`
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	Owner       string            `yaml:"owner"`
	Metadata    map[string]string `yaml:"metadata"`
	Tasks       []TasksStruct     `yaml:"tasks"`
}
type TasksStruct struct {
	Title       string        `yaml:"title"`
	Description string        `yaml:"description"`
	Type1       string        `yaml:"type"`
	Process     ProcessStruct `yaml:"process"`
}
type ProcessStruct struct {
	Org             string            `yaml:"org"`
	Project         string            `yaml:"project"`
	Repo            string            `yaml:"repo"`
	Entrypoint      string            `yaml:"entrypoint"`
	RepoBranchOrTag string            `yaml:"repoBranchOrTag"`
	Arguments       map[string]string `yaml:"arguments"`
}
