package domain

type Employees struct {
	EmployeeId uint64 `gorm:"primaryKey;column:employee_id" json:"employee_id"`
	UserId     uint64 `gorm:"column:user_id" json:"user_id"`

	Nrp int `gorm:"column:nrp;unique" json:"nrp"`

	DepartmentId uint64 `gorm:"column:department_id" json:"department_id"`
	PositionId   uint64 `gorm:"column:position_id" json:"position_id"`

	BaseModel
}

func (Employees) TableName() string {
	return "employees"
}

type EmployeeDepartments struct {
	DepartmentId   uint64 `gorm:"primaryKey;column:department_id" json:"department_id"`
	DepartmentName string `gorm:"column:department_name;size:255" json:"department_name"`
	BaseModel
}

func (EmployeeDepartments) TableName() string {
	return "employee_departments"
}

type EmployeePositions struct {
	PositionId   uint64 `gorm:"primaryKey;column:position_id" json:"position_id"`
	PositionName string `gorm:"column:position_name;size:255" json:"position_name"`
	BaseModel
}

func (EmployeePositions) TableName() string {
	return "employee_positions"
}
