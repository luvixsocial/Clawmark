package constants

const (
	ResourceNotFound    = `{"success":false,"message":"Slow down, bucko! We couldn't find this resource *anywhere*!"}`
	FileNotFound        = `{"success":false,"message":"Slow down, bucko! We couldn't find the requested file *anywhere*!"}`
	EndpointNotFound    = `{"success":false,"message":"Slow down, bucko! You got the path wrong or something but this endpoint doesn't exist!"}`
	BadRequest          = `{"success":false,"message":"Slow down, bucko! You're doing something illegal!!!"}`
	Forbidden           = `{"success":false,"message":"Slow down, bucko! You're not allowed to do this!"}`
	Unauthorized        = `{"success":false,"message":"Slow down, bucko! You're not authorized to do this or did you forget a API token somewhere?"}`
	InternalServerError = `{"success":false,"message":"Slow down, bucko! Something went wrong on our end!"}`
	MethodNotAllowed    = `{"success":false,"message":"Slow down, bucko! That method is not allowed for this endpoint!!!"}`
	BodyRequired        = `{"success":false,"message":"Slow down, bucko! A body is required for this endpoint!!!"}`
	BackTick            = "`"
	DoubleBackTick      = "``"
)
