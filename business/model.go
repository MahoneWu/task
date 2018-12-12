package business

//define user obj struct
type User struct {
	Name string
	Password string
}

//define response struct
type Response struct {
	code string
	message string
}

//define quota struct
type UserQuota struct {
	WriteSpeed int
	ReadSpeed int
	UserId int64
}



const (
	WRITE = "WRITE"
	READ = "READ"
	LOGIN = "LOGIN"
)
