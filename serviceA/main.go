package main

import (
	"errors"
	"fmt"
	"log"
	UserModel "myservice/Entity/User"
	"myservice/myservice_pb/myservice_pb"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	cc, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("err while dial %v", err)
	}
	defer cc.Close()

	client := myservice_pb.NewTestRPCServiceClient(cc)
	log.Println("serice client ", client)

	r := gin.Default()
	r.POST("/users", createUser(client))
	r.GET("/users/:id", getUser(client))
	r.DELETE("/users/:id", deleteUser(client))
	r.PUT("/users/:id", updateUser(client))
	r.Run()

}

// define funcion call sign in rpc
func callSignUp(c myservice_pb.TestRPCServiceClient, user *UserModel.User) {
	log.Println("callSignUp calling")

	SignUpRequest := myservice_pb.SignUpRequest{
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	}
	resp, err := c.SignUp(context.Background(), &SignUpRequest)
	if err != nil {
		log.Fatalf("call Sign Up err %v", err)
	}
	user.ID = resp.GetId()
	log.Printf("call Sign Up response %v %v", resp.GetResponse(), resp.GetId())

}

// define POST handler function create user
func createUser(client myservice_pb.TestRPCServiceClient) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		user := UserModel.User{}
		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		callSignUp(client, &user)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"name":    user.Name,
			"id":      user.ID,
		})

		fmt.Printf("create user %v", user.Name)

	}
}

// define rpc function call get user
func callGetUser(c myservice_pb.TestRPCServiceClient, id string) UserModel.User {
	log.Println("callGetUserByID calling")

	//init context and set time out to context
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	//add metadata to context
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("age", "25666", "name", "phu van hop"),
	)

	//add value to context
	ctx = context.WithValue(ctx, "key1", "value1")

	getUserRequest := myservice_pb.GetUserRequest{
		Id: id,
	}

	fmt.Println(ctx)

	resp, err := c.GetUserByID(ctx, &getUserRequest)
	fmt.Println(err, "=========================")

	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("ContextDeadlineExceeded: true")
	}
	if os.IsTimeout(err) {
		log.Println("IsTimeoutError: true")

	}

	//if err != nil {
	//	log.Fatalf("call Sign Up err %v", err)
	//}

	log.Printf(" service A is running")
	user := UserModel.User{
		Name:  resp.GetName(),
		Age:   resp.GetAge(),
		Email: resp.GetEmail(),
	}
	return user
}

// define GET gin handler func to get user by id
func getUser(client myservice_pb.TestRPCServiceClient) gin.HandlerFunc {

	return func(c *gin.Context) {
		var user UserModel.User
		//extract id from http param
		id := c.Param("id")
		// call rpc funtion
		user = callGetUser(client, id)
		if user.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
			return
		}
		//return json
		c.JSON(http.StatusOK, gin.H{"data": user})

	}
}

// define rpc function call delete user
func callDeleteUser(client myservice_pb.TestRPCServiceClient, id string) string {
	log.Println("callDeleteUser calling")

	deleteUserRequest := myservice_pb.DeleteUserRequest{
		Id: id,
	}
	resp, err := client.DeleteUser(context.Background(), &deleteUserRequest)
	if err != nil {
		log.Fatalf("call Delete err %v", err)
	}

	return resp.GetResponse()
}

// define DELETE gin handler func to delete user by id
func deleteUser(client myservice_pb.TestRPCServiceClient) gin.HandlerFunc {

	return func(c *gin.Context) {
		//extract id from http param
		id := c.Param("id")
		//call rpc function
		resp := callDeleteUser(client, id)
		//return json
		c.JSON(http.StatusOK, gin.H{"message": resp})
	}
}

// define rpc function to call update user
func callUpdateUser(client myservice_pb.TestRPCServiceClient, user *UserModel.User) string {
	log.Println("callUpdate calling")

	UpdateUserRequest := myservice_pb.UpdateUserRequest{
		Id:    user.ID,
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	}

	fmt.Println("UpdateUserRequest ", UpdateUserRequest)
	resp, err := client.UpdateUser(context.Background(), &UpdateUserRequest)
	if err != nil {
		log.Fatalf("call Sign Up err %v", err)
	}
	return resp.GetResponse()
}

// define PUT gin handler func to update user by id
func updateUser(client myservice_pb.TestRPCServiceClient) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//extract id from http param
		id := ctx.Param("id")

		//parse data to user struct
		user := UserModel.User{}

		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		user.ID = id
		fmt.Println(user)
		//call rpc function
		resp := callUpdateUser(client, &user)
		//return json
		ctx.JSON(http.StatusOK, gin.H{"message": resp})
	}
}
