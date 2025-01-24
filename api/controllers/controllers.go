git merge origin/mainpackage controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daiki-trnsk/go-debatemap/api/database"
	"github.com/daiki-trnsk/go-debatemap/api/models"
	generate "github.com/daiki-trnsk/go-debatemap/api/tokens"
	"github.com/labstack/echo/v4"

	// "github.com/labstack/echo/v4/middleware"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var DebateMapCollection *mongo.Collection = database.DebateMapData(database.Client, "debatemaps")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func SignUp(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	validationErr := Validate.Struct(user)
	if validationErr != nil {
		fmt.Println("Validation Error:", validationErr)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": validationErr.Error()})
	}

	count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User already exists"})
	}
	password := HashPassword(*user.Password)
	user.Password = &password

	user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_ID = user.ID.Hex()
	token, refreshtoken, err := generate.TokenGenerator(*user.Email, *user.Nickname, user.User_ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating tokens"})
	}
	user.Token = &token
	user.Refresh_Token = &refreshtoken
	_, inserterr := UserCollection.InsertOne(ctx, user)
	if inserterr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "the user did not created"})
	}
	return c.JSON(http.StatusCreated, "Succesfully signed in!")
}

func Login(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	var founduser models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
	defer cancel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
	}
	PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
	defer cancel()
	if !PasswordIsValid {
		fmt.Println(msg)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": msg})
	}
	token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.Nickname, founduser.User_ID)
	defer cancel()
	generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)
	fmt.Println("founduser", http.StatusFound, founduser)
	return c.JSON(http.StatusOK, founduser)
}

func AddDebateMap(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var debatemap models.DebateMap
	defer cancel()
	if err := c.Bind(&debatemap); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Printf("Received DebateMap: %+v\n", debatemap)
	debatemap.ID = primitive.NewObjectID()
	debatemap.RegistrationDate = time.Now()
	debatemap.NodesJSON = map[string]interface{}{
		"id":       "root",
		"topic":    debatemap.RootNodeTopic,
		"expanded": true,
		"children": []interface{}{},
	}
	if debatemap.NodesJSON["topic"] == "" {
		return c.JSON(http.StatusBadRequest, "Enter Topic")
	}
	_, anyerr := DebateMapCollection.InsertOne(ctx, debatemap)
	if anyerr != nil {
		return c.JSON(http.StatusInternalServerError, "Not Created")
	}
	defer cancel()
	return c.JSON(http.StatusOK, "Map Created!")
}

func GetDebateMaps(c echo.Context) error {
	var debatemaplist []models.DebateMap
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := DebateMapCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
	}
	err = cursor.All(ctx, &debatemaplist)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, "Error occurred while fetching debatemaps")
	}
	defer cursor.Close(ctx)
	if err := cursor.Err(); err != nil {
		// Don't forget to log errors. I log them really simple here just
		// to get the point across.
		log.Println(err)
		return c.JSON(400, "invalid")
	}
	defer cancel()
	return c.JSON(200, debatemaplist)
}

func GetDebateMap(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var debatemap models.DebateMap
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	err := DebateMapCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&debatemap)
	if err != nil {
		return c.JSON(404, "Invalid DebateMap ID")
	}
	defer cancel()
	return c.JSON(200, debatemap)
}

func UpdateDebateMap(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var debatemap models.DebateMap
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := c.Bind(&debatemap); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	update := bson.M{
		"$set": bson.M{
			"root_node_topic": debatemap.NodesJSON["topic"],
			"nodes_json": debatemap.NodesJSON,
		},
	}
	_, err := DebateMapCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error occurred while updating debatemap"})
	}
	defer cancel()
	return c.JSON(http.StatusOK, "DebateMap updated successfully")
}

func SearchDebateMapByQuery(c echo.Context) error {
	var searchdebatemaps []models.DebateMap
	queryParam := c.QueryParam("name")
	if queryParam == "" {
		log.Println("query is empty")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid Search Index"})
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	searchquerydb, err := DebateMapCollection.Find(ctx, bson.M{"root_node_topic": bson.M{"$regex": queryParam, "$options": "i"}})
	if err != nil {
		return c.JSON(404, "something went wrong in fetching the dbquery")
	}
	err = searchquerydb.All(ctx, &searchdebatemaps)
	if err != nil {
		log.Println(err)
		return c.JSON(400, "invalid")
	}
	defer searchquerydb.Close(ctx)
	if err := searchquerydb.Err(); err != nil {
		log.Println(err)
		return c.JSON(400, "invalid request")
	}
	defer cancel()
	return c.JSON(200, searchdebatemaps)
}
