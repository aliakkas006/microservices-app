package models

// type Todo struct {
// 	ID          int64  `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	Completed   bool   `json:"completed"`
// }

type Todo struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
    UserEmail   string `json:"userEmail"` // Must match producer's field name
}
