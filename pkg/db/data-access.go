package db

import "fmt"

type ApplicationResponse struct {
	Name      string `json:"name"`
	APIToken  string `json:"api_token"`
	CreatedAt string `json:"created_at"`
	APISecret string `json:"api_secret"`
}

func CreateApplicationAndApiTokenAndSecret(name string, apiToken string, apiSecret string) error {
	app := Application{
		Name:      name,
		APIToken:  apiToken,
		APISecret: apiSecret,
	}
	something := GetSQLiteDB()
	if err := something.Create(&app).Error; err != nil {
		return err
	}
	return nil
}

func GetAllApplications() ([]ApplicationResponse, error) {
	var unFilteredApplications []Application
	something := GetSQLiteDB()
	if err := something.Find(&unFilteredApplications).Error; err != nil {
		return nil, err
	}
	applications := make([]ApplicationResponse, len(unFilteredApplications))
	for i, app := range unFilteredApplications {
		applications[i] = ApplicationResponse{
			Name:      app.Name,
			APIToken:  app.APIToken,
			CreatedAt: app.CreatedAt.String(),
			APISecret: app.APISecret,
		}
	}
	fmt.Println(applications)
	return applications, nil
}

func UpdateApplicationTokenAndSecret(name string, newToken string, newSecret string) error {
	something := GetSQLiteDB()
	if err := something.Model(&Application{}).Where("name = ?", name).Updates(Application{
		APIToken:  newToken,
		APISecret: newSecret,
	}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteApplication(name string) error {
	something := GetSQLiteDB()
	if err := something.Where("name = ?", name).Delete(&Application{}).Error; err != nil {
		return err
	}
	fmt.Println("Application deleted successfully")
	return nil
}

func GetApplicationByTokenAndSecret(token string, secret string) (*Application, error) {
	var app Application
	something := GetSQLiteDB()
	if err := something.Where("api_token = ? AND api_secret = ?", token, secret).First(&app).Error; err != nil {
		return nil, err
	}

	return &app, nil
}
