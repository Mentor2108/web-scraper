package data

import (
	"backend-service/defn"
	"backend-service/util"
	"context"
)

func ActionByMap(ctx context.Context, action defn.DatabaseAction, data map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	switch action {
	case defn.DatabaseActionCreate:
		return CreateEntityByMap(ctx, data)
	case defn.DatabaseActionRead:
		return ReadEntityByMap(ctx, data)
	case defn.DatabaseActionUpdate:
		return UpdateEntityByMap(ctx, data)
	case defn.DatabaseActionDelete:
		return DeleteEntityByMap(ctx, data)
	default:
		cerr := util.NewCustomError(ctx, defn.ErrCodeInvalidDatabaseAction, defn.ErrInvalidDatabaseAction)
		log.Println(cerr)
		return nil, cerr
	}
}

func CreateEntityByMap(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	// db := GetDatabaseConnection()
	// // db.Pool.CopyFrom()
	// a := pgx.Conn{}
	// a.Cop
	// args := pgx.StrictNamedArgs(data)
	return nil, nil
}

func ReadEntityByMap(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	return nil, nil
}

func UpdateEntityByMap(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	return nil, nil
}

func DeleteEntityByMap(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	return nil, nil
}
