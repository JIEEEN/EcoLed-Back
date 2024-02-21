package forms

import "github.com/go-playground/validator/v10"

// Paylog Service's input value (from body)
type PaylogForm struct {
	Date     string  `form:"date" json:"date" binding:"required,datetime=2006-01-02"`
	Time     string  `form:"time" json:"time" binding:"required,datetime=15:04"`
	Content  string  `form:"content" json:"content" binding:"required,min=1,max=100"`
	Cost     int64   `form:"cost" json:"cost" binding:"required"`
	Name     string  `form:"name" json:"name" binding:"required,min=1,max=30"`
	Place    string  `form:"place" json:"place" binding:"required,min=1,max=30"`
	Material string  `form:"material" json:"material" binding:"required,oneof=plastic can metal paper food glass vinyl styrofoam weee trash n_a"`
	Ecoscore float64 `form:"ecoscore" json:"ecoscore" binding:"required,min=-1000,max=1000"`
}

func (f PaylogForm) Validate() string {
	validate := validator.New()
	err := validate.Struct(f)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Date":
				return f.DateError(err.Tag())
			case "Time":
				return f.TimeError(err.Tag())
			case "Content":
				return f.ContentError(err.Tag())
			case "Cost":
				return f.CostError(err.Tag())
			case "Name":
				return f.NameError(err.Tag())
			case "Place":
				return f.PlaceError(err.Tag())
			case "Material":
				return f.MaterialError(err.Tag())
			case "Ecoscore":
				return f.EcoscoreError(err.Tag())
			}
		}
	}
	return ""
}

// Custom validation error messages for Date field
func (f PaylogForm) DateError(tag string) string {
	switch tag {
	case "required":
		return "Date is required"
	case "datetime":
		return "Date must be in the format of YYYYMMDD"
	default:
		return "Invalid date"
	}
}

// Custom validation error messages for Time field
func (f PaylogForm) TimeError(tag string) string {
	switch tag {
	case "required":
		return "Time is required"
	case "datetime":
		return "Time must be in the format of HHMM"
	default:
		return "Invalid time"
	}
}

// Custom validation error messages for Content field
func (f PaylogForm) ContentError(tag string) string {
	switch tag {
	case "required":
		return "Content is required"
	case "min":
		return "Content is too short"
	case "max":
		return "Content is too long"
	default:
		return "Invalid content"
	}
}

// Custom validation error messages for Cost field
func (f PaylogForm) CostError(tag string) string {
	switch tag {
	case "required":
		return "Cost is required"
	default:
		return "Invalid cost"
	}
}

// Custom validation error messages for Name field
func (f PaylogForm) NameError(tag string) string {
	switch tag {
	case "required":
		return "Name is required"
	case "min":
		return "Name is too short"
	case "max":
		return "Name is too long"
	default:
		return "Invalid name"
	}
}

// Custom validation error messages for Place field
func (f PaylogForm) PlaceError(tag string) string {
	switch tag {
	case "required":
		return "Place is required"
	case "min":
		return "Place is too short"
	case "max":
		return "Place is too long"
	default:
		return "Invalid place"
	}
}

// Custom validation error messages for Material field
func (f PaylogForm) MaterialError(tag string) string {
	switch tag {
	case "required":
		return "Material is required"
	case "oneof":
		return "Material must be one of plastic, can, metal, paper, food, glass, vinyl, styrofoam, weee, trash, n_a"
	default:
		return "Invalid material"
	}
}

// Custom validation error messages for Ecoscore field
func (f PaylogForm) EcoscoreError(tag string) string {
	switch tag {
	case "required":
		return "Ecoscore is required"
	case "min":
		return "Ecoscore can't be lower than -10"
	case "max":
		return "Ecoscore can't be higher than +10"
	default:
		return "Invalid ecoscore"
	}
}
