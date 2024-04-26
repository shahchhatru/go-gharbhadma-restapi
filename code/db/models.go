package db

import(
	
	"time"
)

type Credentials struct {
    ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Email string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
    Pass  string `json:"pass"`
}


type UserProfile struct {
    ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Username     string `json:"username"`
    Credentials  Credentials `gorm:"foreignKey:CredentialsID"`
    CredentialsID uint // Foreign key
}

// Task represents a task object
type Task struct {
    ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}

// TaskReminder represents a reminder associated with a task
type TaskReminder struct {
    ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    
    ReminderType string         `json:"reminder_type"` // e.g., daily, weekly
    Date         *time.Time     `json:"date,omitempty"` // Date for one-time reminder
    Day          string         `json:"day,omitempty"`  // Day of the week for weekly reminder
    Time         string         `json:"time,omitempty"` // Time of the day for notification
    TaskID       uint           `json:"task_id"`
	Task         Task           `gorm:"foreignKey:TaskID"`
	UserProfileID uint           `json:"user_profile_id"` // Foreign key for UserProfile
    UserProfile   UserProfile   `gorm:"foreignKey:UserProfileID"`
}



