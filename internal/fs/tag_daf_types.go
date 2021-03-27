package fs

type TagGetAllDafT = func() ([]string, error)

type TagAddDafT = func(newTags []string) error
