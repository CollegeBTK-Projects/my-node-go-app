package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(coll *mongo.Collection) *UserRepository {
	return &UserRepository{coll: coll}
}

func (u *UserRepository) InsertOne(ctx context.Context, um models.User) (primitive.ObjectID, error) {
	res, err := u.coll.InsertOne(ctx, um)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (u *UserRepository) InsertMany(ctx context.Context, um []models.User) error {
	docs := make([]interface{}, len(um))

	for i, u := range um {
		docs[i] = u
	}

	_, err := u.coll.InsertMany(ctx, docs)
	return err
}

func (r *UserRepository) FindByAll(ctx context.Context) ([]models.User, error) {
	curs, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNilCursor) {
			return nil, err
		}
		return nil, err
	}
	defer curs.Close(ctx)

	var u []models.User
	if err := curs.All(ctx, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var u models.User
	err := r.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) FindByGroup(ctx context.Context, group string) ([]models.User, error) {
	filter := bson.D{{Key: "group_name", Value: group}}
	curs, err := r.coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNilCursor) {
			return nil, err
		}
		return nil, err
	}
	defer curs.Close(ctx)

	var u []models.User
	if err := curs.All(ctx, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) SortByGrades(ctx context.Context) ([]models.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "grade", Value: -1}})

	curs, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNilCursor) {
			return nil, err
		}
		return nil, err
	}
	defer curs.Close(ctx)

	var u []models.User
	if err := curs.All(ctx, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindPaginated(ctx context.Context, page int64, pageSize int64) ([]models.User, error) {
	opts := options.Find().SetSkip((page - 1) * pageSize).SetLimit(pageSize)

	cursor, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateGrade(ctx context.Context, id primitive.ObjectID, grade int) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "grade", Value: grade}}}}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	fmt.Println("Найдено:", res.MatchedCount, "Изменено:", res.ModifiedCount)
	return nil
}

func (r *UserRepository) IncGradesByCourse(ctx context.Context, course int) (int64, int64, error) {
	filter := bson.D{{Key: "course", Value: course}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "grade", Value: 1}}}}

	res, err := r.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, 0, err
	}

	return res.MatchedCount, res.ModifiedCount, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}

	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *UserRepository) DeleteByMinGrade(ctx context.Context, maxGrade int) (int64, error) {
	filter := bson.D{{Key: "grade", Value: bson.D{{Key: "$lt", Value: maxGrade}}}}

	res, err := r.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (r *UserRepository) FindNamesAndGrades(ctx context.Context) ([]models.User, error) {
	opts := options.Find().SetProjection(bson.D{
		{Key: "name", Value: 1},
		{Key: "grade", Value: 1},
		{Key: "_id", Value: 0},
	})

	cursor, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) ProcessUsers(ctx context.Context) error {
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User

		if err := cursor.Decode(&user); err != nil {
			return err
		}

		fmt.Println(user.Name, user.Grade)
	}

	return cursor.Err()
}

func (r *UserRepository) ProcessUsersInBatches(ctx context.Context) error {
	opts := options.Find().SetBatchSize(1000)

	cursor, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User

		if err := cursor.Decode(&user); err != nil {
			return err
		}

		fmt.Println(user.Name)
	}

	return cursor.Err()
}

func (r *UserRepository) AverageGradeByGroup(ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$group_name"},
			{Key: "average_grade", Value: bson.D{
				{Key: "$avg", Value: "$grade"},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "average_grade", Value: -1},
		}}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *UserRepository) CountByCourse(ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$course"},
			{Key: "count", Value: bson.D{
				{Key: "$sum", Value: 1},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id", Value: 1},
		}}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *UserRepository) Top10Students(ctx context.Context) ([]models.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "grade", Value: -1}}).SetLimit(10)

	cursor, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GenerateMillion(ctx context.Context) error {
	const total = 1000000
	const batchSize = 1000

	for start := 0; start < total; start += batchSize {
		end := start + batchSize
		if end > total {
			end = total
		}

		docs := make([]any, 0, end-start)

		for i := start; i < end; i++ {
			docs = append(docs, models.User{
				Name:      fmt.Sprintf("Student%d", i),
				GroupName: fmt.Sprintf("Group%d", i%100),
				Course:    i%4 + 1,
				Grade:     i % 101,
				CreatedAt: time.Now(),
			})
		}

		if _, err := r.coll.InsertMany(ctx, docs); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) ReadStreaming(ctx context.Context) error {
	opts := options.Find().SetBatchSize(1000)

	cursor, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return err
		}
		count++
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	fmt.Println("Обработано:", count)
	return nil
}

func (r *UserRepository) CreateGroupIndex(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "group_name", Value: 1}},
	})
	return err
}

func (r *UserRepository) CreateGroupGradeIndex(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "group_name", Value: 1},
			{Key: "grade", Value: -1},
		},
	})
	return err
}

func (r *UserRepository) ExplainGroupQuery(ctx context.Context) (bson.M, error) {
	filter := bson.D{{Key: "group_name", Value: "477"}}

	command := bson.D{
		{Key: "explain", Value: bson.D{
			{Key: "find", Value: r.coll.Name()},
			{Key: "filter", Value: filter},
		}},
		{Key: "verbosity", Value: "executionStats"},
	}

	var result bson.M
	err := r.coll.Database().RunCommand(ctx, command).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserRepository) CreateTextIndex(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}},
	})

	return err
}

func (r *UserRepository) SearchByName(ctx context.Context, query string) ([]bson.M, error) {
	filter := bson.D{
		{Key: "$text", Value: bson.D{
			{Key: "$search", Value: query},
		}},
	}

	opts := options.Find().
		SetProjection(bson.D{
			{Key: "name", Value: 1},
			{Key: "grade", Value: 1},
			{Key: "score", Value: bson.D{
				{Key: "$meta", Value: "textScore"},
			}},
		}).
		SetSort(bson.D{
			{Key: "score", Value: bson.D{
				{Key: "$meta", Value: "textScore"},
			}},
		})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
