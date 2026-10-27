package repository

import (
	"context"
	"math"
	"strings"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(coll *mongo.Collection) *UserRepository {
	return &UserRepository{coll: coll}
}

func (r *UserRepository) CreateIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "email", Value: "text"},
			},
			Options: options.Index().
				SetName("users_text").
				SetWeights(bson.D{{Key: "name", Value: 5}, {Key: "email", Value: 1}}).
				SetDefaultLanguage("none"),
		},
	})
	return err
}

func (u *UserRepository) CreateUser(ctx context.Context, um models.User) (bson.ObjectID, error) {
	res, err := u.coll.InsertOne(ctx, um)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.NilObjectID, err
		}
		return bson.NilObjectID, err
	}
	return res.InsertedID.(bson.ObjectID), nil
}

func (u *UserRepository) FindById(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	cur, err := u.coll.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var um models.User
	for cur.Next(ctx) {
		if err := cur.Decode(&um); err != nil {
			return nil, err
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &um, nil
}

func (u *UserRepository) FindByAll(ctx context.Context) ([]models.User, error) {
	cur, err := u.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	um := []models.User{}
	if err := cur.All(ctx, &um); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return um, nil
}

func (u *UserRepository) UpdateById(ctx context.Context, id bson.ObjectID, um models.User) (int64, error) {
	update := bson.M{
		"$set": bson.M{
			"name":       um.Name,
			"email":      um.Email,
			"group_name": um.GroupName,
			"age":        um.Age,
		},
	}

	res, err := u.coll.UpdateByID(ctx, id, update)
	if err != nil {
		return 0, err
	}

	if res.ModifiedCount == 0 {
		return 0, nil
	}

	return res.MatchedCount, nil
}

func (u *UserRepository) DeleteById(ctx context.Context, id bson.ObjectID) error {
	del, err := u.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if del.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func buildFilter(q models.UserQuery) bson.M {
	f := bson.M{}
	if q.GroupName != "" {
		f["group_name"] = q.GroupName
	}
	if q.Course != nil {
		f["course"] = *q.Course
	}
	if q.AgeMin != nil || q.AgeMax != nil {
		age := bson.M{}
		if q.AgeMin != nil {
			age["$gte"] = *q.AgeMin
		}
		if q.AgeMax != nil {
			age["$lte"] = *q.AgeMax
		}
		f["age"] = age
	}
	if q.Search != "" {
		f["$text"] = bson.M{"$search": q.Search}
	}
	return f
}

var sortFields = map[string]string{
	"name":       "name",
	"email":      "email",
	"age":        "age",
	"created_at": "createdAt",
}

func buildSort(q models.UserQuery) bson.D {
	if q.Sort != "" {
		dir := 1
		field := q.Sort
		if strings.HasPrefix(field, "-") {
			dir = -1
			field = field[1:]
		}
		return bson.D{{Key: sortFields[field], Value: dir}, {Key: "_id", Value: 1}}
	}
	if q.Search != "" {
		return bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}}
	}
	return bson.D{{Key: "createdAt", Value: -1}, {Key: "_id", Value: 1}}
}

func baseFindOptions(q models.UserQuery) *options.FindOptionsBuilder {
	opts := options.Find().SetSort(buildSort(q))
	if q.Search != "" {
		opts.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
	}
	return opts
}

func (r *UserRepository) List(ctx context.Context, q models.UserQuery) ([]models.User, int64, error) {
	filter := buildFilter(q)
	page, limit := q.PageLimit()

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := baseFindOptions(q).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	users := []models.User{}
	if err := cur.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) ListAll(ctx context.Context, q models.UserQuery) ([]models.User, error) {
	cur, err := r.coll.Find(ctx, buildFilter(q), baseFindOptions(q))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	users := []models.User{}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Stats(ctx context.Context) (*models.Stats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.D{
			{Key: "total", Value: bson.A{
				bson.D{{Key: "$count", Value: "n"}},
			}},
			{Key: "avg", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.M{"age": bson.M{"$gt": 0}}}},
				bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": "$age"}}}},
			}},
			{Key: "by_group", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.M{"_id": "$group_name", "count": bson.M{"$sum": 1}}}},
				bson.D{{Key: "$sort", Value: bson.M{"_id": 1}}},
			}},
			{Key: "by_course", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.M{
					"_id":   bson.M{"$ifNull": bson.A{"$course", 0}}, // 0 = курс не указан
					"count": bson.M{"$sum": 1},
				}}},
				bson.D{{Key: "$sort", Value: bson.M{"_id": 1}}},
			}},
		}}},
	}

	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var raw []struct {
		Total []struct {
			N int64 `bson:"n"`
		} `bson:"total"`
		Avg []struct {
			Avg float64 `bson:"avg"`
		} `bson:"avg"`
		ByGroup  []models.Bucket `bson:"by_group"`
		ByCourse []models.Bucket `bson:"by_course"`
	}
	if err := cur.All(ctx, &raw); err != nil {
		return nil, err
	}

	st := &models.Stats{ByGroup: []models.Bucket{}, ByCourse: []models.Bucket{}}
	if len(raw) == 0 {
		return st, nil
	}
	res := raw[0]
	if len(res.Total) > 0 {
		st.Total = res.Total[0].N
	}
	if len(res.Avg) > 0 {
		st.AvgAge = math.Round(res.Avg[0].Avg*100) / 100
	}
	st.ByGroup, st.ByCourse = res.ByGroup, res.ByCourse
	return st, nil
}
