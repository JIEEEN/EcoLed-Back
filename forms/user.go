package forms

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type UserForm struct {}

// Register service's input value (from body)
type RegisterForm struct {
    Email       string `form:"email" json:"email" binding:"required,email"`
    Password    string `form:"password" json:"password" binding:"required,min=4,max=30"`
    Nickname    string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
    Accountname string `form:"accountname" json:"accountname" binding:"required,min=2,max=30"`
}

// Login service's input value (from body)
type LoginForm struct {
    Email    string `form:"email" json:"email" binding:"required,email"`
    Password string `form:"password" json:"password" binding:"required,min=4,max=30"`
}


// Login service's return value
type UserReturnForm struct {
    Email     string
    Nickname  string
    CreatedAt string 
}

type FindUserInfoForm struct {
    Email string `form:"email" json:"email" binding:"required,email"`
}

type VerifyCodeForm struct {
    Email string `form:"email" json:"email" binding:"required,email"`
    Code string `form:"code" json:"code" binding:"required,min=6,max=6"`
}


// Custom validation error messages for RegisterForm
func (f UserForm) Register(form RegisterForm) string {
	validate := validator.New()
	err := validate.Struct(form)

	if err == nil {
		return ""
	}

    switch err.(type) {
    case validator.ValidationErrors:

        if _, ok := err.(*json.UnmarshalTypeError); ok {
            return "Something went wrong, please try again later"
        }

        for _, err := range err.(validator.ValidationErrors) {
            switch err.Field() {
            case "Email":
                return f.Email(err.Tag())
            case "Password":
                return f.Password(err.Tag())
            case "Nickname":
                return f.Nickname(err.Tag())
            case "Accountname":
                return f.Accountname(err.Tag())
            }

        }
    default:
        return "Invalid request"
    }

    return "Something went wrong, please try again later"
}

// Custom validation error messages for LoginForm
func (f UserForm) Login(form LoginForm) string {
	validate := validator.New()
	err := validate.Struct(form)
	
	if err == nil {
		return ""
	}

    switch err.(type) {
    case validator.ValidationErrors:

        if _, ok := err.(*json.UnmarshalTypeError); ok {
            return "Something went wrong, please try again later"
        }

        for _, err := range err.(validator.ValidationErrors) {
            switch err.Field() {
            case "Email":
                return f.Email(err.Tag())
            case "Password":
                return f.Password(err.Tag())
            }
        }
    default:
        return "Invalid request"
    }

    return "Something went wrong, please try again later"
}

func (f UserForm) FindUserInfo(form FindUserInfoForm) string {
    validate := validator.New()
    err := validate.Struct(form)
    
    if err == nil {
        return ""
    }

    switch err.(type) {
    case validator.ValidationErrors:

        if _, ok := err.(*json.UnmarshalTypeError); ok {
            return "Something went wrong, please try again later"
        }

        for _, err := range err.(validator.ValidationErrors) {
            switch err.Field() {
            case "Email":
                return f.Email(err.Tag())
            }
        }
    default:
        return "Invalid request"
    }

    return "Something went wrong, please try again later"
}

func (f UserForm) VerifyCode(form VerifyCodeForm) string {
    validate := validator.New()
    err := validate.Struct(form)
    if err == nil {
        return ""
    }

    switch err.(type) {
    case validator.ValidationErrors:

        if _, ok := err.(*json.UnmarshalTypeError); ok {
            return "Something went wrong, please try again later"
        }

        for _, err := range err.(validator.ValidationErrors) {
            switch err.Field() {
            case "Email":
                return f.Email(err.Tag())
            case "Code":
                return "Invalid code"
            }
        }
    default:
        return "Invalid request"
    }

    return "Something went wrong, please try again later"
}

// Custom validation error messages for each field
func (f UserForm) Email(tag string) string {
    switch tag {
    case "required":
        return "Email is required"
    case "email":
        return "Invalid email format"
    default:
        return "Invalid email"
    }
}

func (f UserForm) Password(tag string) string {
    switch tag {
    case "required":
        return "Password is required"
    case "min":
        return "Password is too short"
    case "max":
        return "Password is too long"
    default:
        return "Invalid password"
    }
}

func (f UserForm) Nickname(tag string) string {
    switch tag {
    case "required":
        return "Nickname is required"
    case "min":
        return "Nickname is too short"
    case "max":
        return "Nickname is too long"
    default:
        return "Invalid nickname"
    }
}

func (f UserForm) Accountname(tag string) string {
    switch tag {
    case "required":
        return "Account name is required"
    case "min":
        return "Account name is too short"
    case "max":
        return "Account name is too long"
    default:
        return "Invalid account name"
    }
}

func (f UserForm) Code(tag string) string {
    switch tag {
    case "required":
        return "Code is required"
    case "min":
        return "Code must be 6 digits"
    case "max":
        return "Code must be 6 digits"
    default:
        return "Invalid code"
    }
}