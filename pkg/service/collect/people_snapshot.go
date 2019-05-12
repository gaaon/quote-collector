package collect

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type peopleSnapshotService struct {
	snapshotLocation string
}

func NewPeopleSnapshotService() *peopleSnapshotService {
	return &peopleSnapshotService{
		snapshotLocation: "data/snapshot.txt",
	}
}

func (service *peopleSnapshotService) writePeopleListIntoWriter(
	writer io.Writer, peopleList []model.Person) (err error) {

	for i, name := range peopleList {
		_, _ = writer.Write([]byte(
			name.FullName + "\t" + name.ReversedName + "\t" + name.Link))

		if i != len(peopleList)-1 {
			_, _ = writer.Write([]byte("\n"))
		}
	}

	return
}

func (service *peopleSnapshotService) SavePeopleList(peopleList []model.Person) (err error) {
	var f *os.File
	if f, err = os.Create(service.snapshotLocation); err != nil {
		return
	}
	defer f.Close()

	return service.writePeopleListIntoWriter(f, peopleList)
}

func (service *peopleSnapshotService) IsPeopleListExist() bool {
	_, err := os.Stat(service.snapshotLocation)

	return !os.IsNotExist(err)
}

func (service *peopleSnapshotService) readPeopleListFromReader(reader io.Reader) (
	peopleList []model.Person, err error) {

	var content []byte
	if content, err = ioutil.ReadAll(reader); err != nil {
		return
	}

	splits := strings.Split(string(content), "\n")
	for _, split := range splits {
		values := strings.Split(split, "\t")
		peopleList = append(peopleList, model.Person{FullName: values[0], ReversedName: values[1], Link: values[2]})
	}

	return
}

func (service *peopleSnapshotService) FindPeopleList() (
	peopleList []model.Person, err error) {

	var f *os.File
	if f, err = os.Open(service.snapshotLocation); err != nil {
		return
	}
	defer f.Close()

	return service.readPeopleListFromReader(f)
}
