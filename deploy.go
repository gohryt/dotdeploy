package dotdeploy

import "errors"

type (
	Deploy struct {
		Folder string
		Keep   bool

		Remote Remote
		Do     Do
	}

	Path struct {
		Connection string
		Path       string
	}
)

var (
	ErrDeployFolderEmpty = errors.New(`deploy.Folder == ""`)
)

func (deploy *Deploy) Prepare() error {
	err := deploy.Validate()
	if err != nil {
		return err
	}

	base := Do(nil)

	for i := range deploy.Do {
		action := deploy.Do[i]

		switch action.Data.(type) {
		case *Copy:
			copy := action.Data.(*Copy)
			copyMeta := &CopyMeta{
				Path: deploy.Folder,
			}

			if copy.From.Connection != "" {
				copyMeta.From, err = deploy.Remote.Find(copy.From.Connection)
				if err != nil {
					return err
				}
			}

			if copy.To.Connection != "" {
				copyMeta.To, err = deploy.Remote.Find(copy.To.Connection)
				if err != nil {
					return err
				}
			}

			action.Meta = copyMeta
		case *Move:
			move := action.Data.(*Move)
			moveMeta := &MoveMeta{
				Path: deploy.Folder,
			}

			if move.From.Connection != "" {
				moveMeta.From, err = deploy.Remote.Find(move.From.Connection)
				if err != nil {
					return err
				}
			}

			action.Meta = moveMeta
		case *Execute:
			execute := action.Data.(*Execute)
			executeMeta := new(ExecuteMeta)

			if execute.Path.Connection != "" {
				executeMeta.Path, err = deploy.Remote.Find(action.Data.(*Execute).Path.Connection)
				if err != nil {
					return err
				}
			}

			action.Meta = executeMeta
		}

		if action.Follow != "" {
			follow, err := deploy.Do.Find(action.Follow)
			if err != nil {
				return err
			}

			follow.Next = append(follow.Next, action)
		} else {
			base = append(base, action)
		}
	}

	deploy.Do = base
	return nil
}

func (deploy *Deploy) Validate() error {
	join := []error(nil)

	if deploy.Folder == "" {
		join = append(join, ErrDeployFolderEmpty)
	}

	for i := range deploy.Remote {
		err := deploy.Remote[i].Data.(Validable).Validate()
		if err != nil {
			join = append(join, err)
		}
	}

	for i := range deploy.Do {
		err := deploy.Do[i].Data.(Validable).Validate()
		if err != nil {
			join = append(join, err)
		}
	}

	return errors.Join(join...)
}
