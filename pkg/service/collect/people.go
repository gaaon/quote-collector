package collect

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
)

func FindPeopleListFromSnapshot() (peopleList []model.Person, err error) {
	if repository.IsExistPeopleListSnapshot() {
		if peopleList, err = FindPeopleListFromAToZ(); err != nil {
			return
		}

		if err = repository.SavePeopleListIntoSnapshot(peopleList); err != nil {
			return
		}

		println("Read from server")

		return
	} else {
		if peopleList, err = repository.FindPeopleListFromSnapshot(); err != nil {
			return
		}

		println("Read from local")

		return
	}
}