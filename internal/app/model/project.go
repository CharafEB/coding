package model

type Project struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	StidentIDB        int    `json:"student_id"`
	ProjectPartnerIDs string `json:"project_partner_id"`
	Status            string `json:"status"`
	Date              string `json:"date"`
	DateEdit          string `json:"date_edin"`
	SupportStucter    string `json:"role"`
}
