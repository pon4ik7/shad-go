//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
outer:
	for {
		getResp, getErr := c.Get(&kvapi.GetRequest{Key: key})

		var oldValue *string
		var oldVersion uuid.UUID

		switch {
		case getErr == nil:
			oldValue = &getResp.Value
			oldVersion = getResp.Version

		case errors.Is(getErr, kvapi.ErrKeyNotFound):
			oldValue = nil
			oldVersion = uuid.Nil

		case isAuthError(getErr):
			return getErr

		default:
			continue outer
		}

		newValue, updateErr := updateFn(oldValue)
		if updateErr != nil {
			return updateErr
		}

		var attemptedVersions []uuid.UUID
	setLoop:
		for {
			newVersion := uuid.Must(uuid.NewV4())
			attemptedVersions = append(attemptedVersions, newVersion)

			_, setErr := c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: oldVersion,
				NewVersion: newVersion,
			})

			switch {
			case setErr == nil:
				return nil

			case isConflictError(setErr):
				conflictErr := extractConflictError(setErr)
				for _, attempted := range attemptedVersions {
					if conflictErr.ExpectedVersion == attempted {
						return nil
					}
				}
				continue outer

			case errors.Is(setErr, kvapi.ErrKeyNotFound):
				createValue, err := updateFn(nil)
				if err != nil {
					return err
				}

				var createAttemptedVersions []uuid.UUID
			createLoop:
				for {
					createVersion := uuid.Must(uuid.NewV4())
					createAttemptedVersions = append(createAttemptedVersions, createVersion)

					_, createErr := c.Set(&kvapi.SetRequest{
						Key:        key,
						Value:      createValue,
						OldVersion: uuid.Nil,
						NewVersion: createVersion,
					})

					switch {
					case createErr == nil:
						return nil

					case isConflictError(createErr):
						createConflictErr := extractConflictError(createErr)
						for _, attempted := range createAttemptedVersions {
							if createConflictErr.ExpectedVersion == attempted {
								return nil
							}
						}
						continue outer

					case isAuthError(createErr):
						return createErr

					default:
						continue createLoop
					}
				}

			case isAuthError(setErr):
				return setErr

			default:
				continue setLoop
			}
		}
	}
}

func isConflictError(err error) bool {
	var conflictErr *kvapi.ConflictError
	return errors.As(err, &conflictErr)
}

func extractConflictError(err error) *kvapi.ConflictError {
	var conflictErr *kvapi.ConflictError
	errors.As(err, &conflictErr)
	return conflictErr
}

func isAuthError(err error) bool {
	var authErr *kvapi.AuthError
	return errors.As(err, &authErr)
}
