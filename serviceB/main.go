package main

import (
	json2 "encoding/json"
	"fmt"
	"log"
	UserModel "myservice/Entity/User"
	MySQL "myservice/Storage/MySQL"
	Redis "myservice/Storage/Redis"
	"myservice/myservice_pb/myservice_pb"
	"net"

	"github.com/go-redis/redis"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	cache *redis.Client
	db    *sqlx.DB
}

func (s *server) SignUp(ctx context.Context, req *myservice_pb.SignUpRequest) (*myservice_pb.SignUpResponse, error) {

	//marshal request to json
	json, err := json2.Marshal(UserModel.User{Name: req.GetName(), Age: req.GetAge(), Email: req.GetEmail()})
	log.Println(json)
	if err != nil {
		log.Fatalf("Marshall failed %v", err)
		log.Fatalf("Marshall failed %v", json)
	}

	//create ramdom new id
	id := uuid.New()

	//set value to cache
	//err = s.cache.Set(id.String(), json, 0).Err()
	//if err != nil {
	//	log.Fatalf("Set value to cache failed %v", err)
	//}

	//set value to database
	query := fmt.Sprintf("INSERT INTO Users(ID, Name, Age, Email) VALUES(\"%s\", \"%s\", %d,\"%s\")", id, req.GetName(), req.GetAge(), req.GetEmail())
	if _, err := s.db.Exec(query); err != nil {
		log.Fatalf("Set value to database failed %v", err)
	}

	//define response
	resp := &myservice_pb.SignUpResponse{
		Response: fmt.Sprintf("%v signed up successful", req.GetName()),
		Id:       id.String(),
	}
	log.Printf("%v signed up successful", req.GetName())
	return resp, nil

}
func (s *server) GetUserByID(ctx context.Context, req *myservice_pb.GetUserRequest) (*myservice_pb.GetUserResponse, error) {

	//fmt.Println(ctx)
	md, _ := metadata.FromIncomingContext(ctx)
	//time.Sleep(time.Second * 2)
	fmt.Println(md)
	fmt.Println(md["name"])
	fmt.Println(md["age"])
	fmt.Println(ctx.Value("key1"))

	var user UserModel.User
	//get value in cache
	val, err := s.cache.Get(req.GetId()).Result()
	fmt.Printf("%v, %T", val, val)
	if err != nil || val == "{}" {
		log.Printf("Don't exist this user in case: %v", err)

		//get value in db
		err = s.db.Get(&user, "select * from Users where id=?", req.GetId())
		if err != nil {
			log.Println("Id not exist in database")
		}
		//set value to cache
		json, err := json2.Marshal(UserModel.User{Name: user.Name, Age: user.Age, Email: user.Email})
		if err != nil {
			log.Fatalf("Marshal failed %v", err)
		}
		err = s.cache.Set(req.GetId(), json, 0).Err()
		if err != nil {
			log.Fatalf("Set value to cache failed %v", err)
		}

		//return response
		resp := &myservice_pb.GetUserResponse{
			Name:  user.Name,
			Age:   user.Age,
			Email: user.Email,
		}

		log.Printf("%v get information", user.Name)
		return resp, nil

	}

	//if user exist in redis, unmarshal json to struct user
	json.Unmarshal([]byte(val), &user)

	// return response
	resp := &myservice_pb.GetUserResponse{
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	}
	log.Printf("%v get information", user.Name)
	return resp, nil
}

func (s *server) DeleteUser(ctx context.Context, req *myservice_pb.DeleteUserRequest) (*myservice_pb.DeleteUserResponse, error) {
	//delete user from cache
	s.cache.Del(req.GetId())

	//delete user from database
	_, err := s.db.Exec("DELETE FROM Users where id=?", req.GetId())
	if err != nil {
		log.Fatalf("Error while delete user in database %v", err)
	}

	// return response
	resp := &myservice_pb.DeleteUserResponse{
		Response: fmt.Sprintf("%v delete ok", req.GetId()),
	}
	log.Printf("%v delete user", req.GetId())
	return resp, nil
}

func (s *server) UpdateUser(ctx context.Context, req *myservice_pb.UpdateUserRequest) (*myservice_pb.UpdateUserResponse, error) {
	log.Println("request ", req)
	log.Printf("%v update user", req.GetId())

	//marshal request to json
	json, err := json2.Marshal(UserModel.User{ID: req.GetId(), Name: req.GetName(), Age: req.GetAge(), Email: req.GetEmail()})
	if err != nil {
		log.Fatalf("Marshall failed %v", err)
		log.Fatalf("Marshall failed %v", json)
	}

	//delete user from cache
	s.cache.Del(req.GetId())

	//update to database
	fmt.Println(json)
	fmt.Println(req)
	query := fmt.Sprintf("UPDATE Users SET Name = \"%s\", Age = %d, Email = \"%s\" WHERE ID = \"%s\"", req.GetName(), req.GetAge(), req.GetEmail(), req.GetId())
	_, err = s.db.Exec(query)
	if err != nil {
		log.Fatalf("Error while update user in database %v", err)
	}
	// return response
	resp := &myservice_pb.UpdateUserResponse{
		Response: fmt.Sprintf("%v update ok", req.GetId()),
	}
	return resp, nil
}

func main() {
	//create listen
	lis, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Error while create listen %v", err)
	}

	s := grpc.NewServer()

	cache := Redis.NewRedisClient(context.Background())
	db := MySQL.NewMySQLConnect(context.Background())

	defer cache.Close()
	defer db.Close()

	myservice_pb.RegisterTestRPCServiceServer(s, &server{cache, db})

	fmt.Println("service is running")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Error while serve %v", err)
	}
}
