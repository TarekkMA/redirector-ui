package datasource

import (
	"github.com/TarekkMA/redirector-ui/data"
	"github.com/jinzhu/gorm"
)

var subscribers []func(data []*data.Redirect)

type RedirectsDataSource struct {
	DB *gorm.DB
}

func (r *RedirectsDataSource) Subscribe(callback func(data []*data.Redirect)) {
	subscribers = append(subscribers, callback)
}

func (r *RedirectsDataSource) AddItem(redirect *data.Redirect) error {
	defer r.notify()
	return r.DB.Create(redirect).Error
}

func (r *RedirectsDataSource) RemoveItem(redirect *data.Redirect) error {
	defer r.notify()
	return r.DB.Delete(redirect).Error
}

func (r *RedirectsDataSource) GetAll() ([]*data.Redirect, error) {
	results := []*data.Redirect{}
	return results, r.DB.Find(&results).Error
}

func (r *RedirectsDataSource) notify() {
	data, _ := r.GetAll()
	for i := 0; i < len(subscribers); i++ {
		subscribers[i](data)
	}
}
