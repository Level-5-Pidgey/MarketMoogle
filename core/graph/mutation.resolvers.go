package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"fmt"
)

// CreateUser is the resolver for the CreateUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, lodestoneID *int) (*schema.User, error) {
	/*//Create base object
	user := schema.User{
		LodestoneID: lodestoneID,
	}
	user.IsPremium = util.BoolPointer(false)

	//Make API request to the lodestone to fill the rest of the data
	xivApiProvider := providers.XivApiProvider{}
	lodestoneUser, err := xivApiProvider.GetLodestoneInfoById(lodestoneID)
	characterData := &lodestoneUser.Character
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(characterData.ClassJobs); i++ {
		classJob := characterData.ClassJobs[i]
		newJob := schema.Job{
			JobID: classJob.JobID,
			Level: classJob.Level,
		}

		user.Jobs = append(user.Jobs, newJob)
	}

	user.DataCenter = &characterData.DC
	user.Server = &characterData.Server
	user.PortraitAddress = &characterData.Avatar

	dbErr := r.DbProvider.
	if dbErr != nil {
		return nil, err
	}
	return &user, nil*/

	panic(fmt.Errorf("not implemented"))
}

// UpdateUser is the resolver for the UpdateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, playerID *int, input schema.UserInput) (*schema.User, error) {
	/*updatedUser := schema.User{
		ID:              *playerID,
		LodestoneID:     input.LodestoneID,
		DataCenter:      input.DataCenter,
		Server:          input.Server,
		PortraitAddress: input.PortraitAddress,
		IsPremium:       input.IsPremium,
	}

	//Cast "JobInput" back to "Job"
	for _, job := range input.Jobs {
		newJob := schema.Job{
			ID:    job.ID,
			JobID: job.JobID,
			Level: job.Level,
		}

		updatedUser.Jobs = append(updatedUser.Jobs, newJob)
	}

	dbErr := r.DbProvider.SaveRecipe(updatedUser).Error
	if dbErr != nil {
		return nil, dbErr
	}
	return &updatedUser, nil*/

	panic(fmt.Errorf("not implemented"))
}

// DeleteUser is the resolver for the DeleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, playerID *int) (bool, error) {
	/*internalErrorString := "internal error"
	deleteUser := schema.User{}
	r.DbProvider.Find(&deleteUser, *playerID)

	if tx := r.DbProvider.First(deleteUser, "id = ?", playerID); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return false, tx.Error
		}

		return false, errors.New(internalErrorString)
	}

	if tx := r.DbProvider.Delete(deleteUser, "id = ?", playerID); tx.Error != nil {
		return false, errors.New(internalErrorString)
	}

	return true, nil*/

	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
