package film_models

import (
	"fmt"
)

type LoadInfo struct {
	ID       int
	FilmInfo FilmInfo
	Err      error
}

func (loadInfo LoadInfo) String() string {
	if loadInfo.Err == nil {
		return loadInfo.FilmInfo.String()
	} else {
		return fmt.Sprintf("%d - error: %s", loadInfo.ID, loadInfo.Err)
	}
}
