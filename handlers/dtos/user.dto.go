package dtos

import "jxb-eprocurement/models"

type UserDTO struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Role  RoleDTO `json:"role"`
}

func ToUserDTO(user models.User) UserDTO {
	return UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  ToRoleDTO(user.Role),
	}
}

func ToUserNoRoleDTO(user models.User) UserDTO {
	return UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func ToLoanDTOs(loans []models.Loan) []LoanDTO {
	loanDTOs := make([]LoanDTO, len(loans))
	for i, loan := range loans {
		loanDTOs[i] = ToLoanDTO(loan)
	}
	return loanDTOs
}
