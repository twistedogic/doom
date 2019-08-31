package job

import (
	"log"

	"github.com/fatih/structs"
	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/target"
)

type Job struct {
	Src tap.Tap
	Dst target.Target
}

func New() *Job {
	return &Job{}
}

func (j *Job) Set(src tap.Tap, dst target.Target) {
	j.SetSrc(src)
	j.SetDst(dst)
}

func (j *Job) SetSrc(src tap.Tap) {
	j.Src = src
}

func (j *Job) SetDst(dst target.Target) {
	j.Dst = dst
}

func (j *Job) Execute() error {
	return j.Src.Update(j.Dst)
}

func (j *Job) Run() {
	if err := j.Execute(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Done tap: %s, target: %s", structs.Name(j.Src), structs.Name(j.Dst))
}
