package Models

import(
    "golang.org/x/crypto/bcrypt"
    "fmt"
)

// Create default admin account when system stated
func CreateDefaultAdmin(login, password string) {

    // Check if admin is exists
    var admin Admin
    if err := GetOneAdminByLogin(&admin, login); err == nil {

        // Admin already exists
        return
    }

    fmt.Println("Create account for system administrator [login:", login, "]")

    // Create new admin record
    admin.Login = login

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    admin.Password = string(hashedPassword)

    // Save new admin record
    if err := PutOneAdmin(&admin); err != nil {
        fmt.Println("WARNING! New admin account not created")
    }
}
