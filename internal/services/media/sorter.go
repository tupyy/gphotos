package media

import "github.com/tupyy/gophoto/internal/entity"

type mediaSorter struct {
	medias []entity.Media
}

func newSorter(m []entity.Media) *mediaSorter {
	return &mediaSorter{m}
}

func (ms *mediaSorter) Len() int {
	return len(ms.medias)
}

func (ms *mediaSorter) Swap(i, j int) {
	ms.medias[i], ms.medias[j] = ms.medias[j], ms.medias[i]
}

func (ms *mediaSorter) Less(i, j int) bool {
	return ms.medias[i].CreateDate.Before(ms.medias[j].CreateDate)
}
