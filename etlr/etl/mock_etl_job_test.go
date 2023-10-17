// Copyright © 2023 OSINTAMI. This is not yours.
package etl

type MockETLJob struct {
	tools  *Toolbox
	source *Source
	info   *ETLJobInfo
}

func NewMockETLJob(tools *Toolbox, source *Source, info *ETLJobInfo) IETLJob {
	return &MockETLJob{
		tools:  tools,
		source: source,
		info:   info}
}
func (x *MockETLJob) Refresh() error {
	return nil
}
func (x *MockETLJob) Extract() error {
	return nil
}
func (x *MockETLJob) Transform() error {
	return nil
}
func (x *MockETLJob) Load() error {
	return nil
}
func (x *MockETLJob) Publish() error {
	return nil
}
func (x *MockETLJob) Info() *ETLJobInfo {
	return x.info
}
func (x *MockETLJob) Source() *Source {
	return x.source
}
func (x *MockETLJob) Tools() *Toolbox {
	return x.tools
}
