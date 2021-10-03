package service

import (
	"strconv"
	"time"
	"to-do-list/model"
	"to-do-list/pkg/e"
	"to-do-list/pkg/logging"
	"to-do-list/serializer"
)

//展示任务详情的服务
type ShowTaskService struct {
}

//删除任务的服务
type DeleteTaskService struct {
}

//更新任务的服务
type UpdateTaskService struct {
	ID            uint   `form:"id" json:"id"`
	Title         string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Content          string `form:"content" json:"content" binding:"max=1000"`
	Status 		  int   `form:"status" json:"status"`   //0 待办   1已完成
}

//创建任务的服务
type CreateTaskService struct {
	Title         string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Content          string `form:"content" json:"content" binding:"max=1000"`
	Status 		  int    `form:"status" json:"status"`   //0 待办   1已完成
	StartTime time.Time `json:"start_time" form:"start_time"`
	EndTime time.Time `json:"end_time" form:"end_time"`
}

//搜索任务的服务
type SearchTaskService struct {
	Info string `form:"info" json:"info"`
}

type ListTasksService struct {
	Limit      int  `form:"limit" json:"limit"`
	Start      int  `form:"start" json:"start"`
	CategoryID uint `form:"category_id" json:"category_id"`
}

func (service *CreateTaskService) Create(id string) serializer.Response {
	idInt,_ := strconv.Atoi(id)
	task := model.Task{
		UserId: uint(idInt),
		Title:         service.Title,
		Content:          service.Content,
		Status:          0,
		StartTime:service.StartTime,
		EndTime:service.EndTime,
	}
	code := e.SUCCESS
	err := model.DB.Create(&task).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    e.GetMsg(code),
	}
}

func (service *ListTasksService) List(id string) serializer.Response {
	var tasks []model.Task
	total := 0
	if service.Limit == 0 {
		service.Limit = 15
	}
	model.DB.Model(model.Task{}).Where("uid = ?",id).Count(&total).
		Limit(service.Limit).Offset((service.Start-1)*service.Limit).
		Find(&tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(total))
}

func (service *ShowTaskService) Show(id string) serializer.Response {
	var task model.Task
	code := e.SUCCESS
	err := model.DB.First(&task, id).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	task.AddView()  //增加点击数
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    e.GetMsg(code),
	}
}

func (service *DeleteTaskService) Delete(id string) serializer.Response {
	var task model.Task
	code := e.SUCCESS
	err := model.DB.First(&task, id).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	err = model.DB.Delete(&task).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *UpdateTaskService) Update(id string) serializer.Response {
	var task model.Task
	model.DB.Model(model.Task{}).Where("id = ?",task).First(&task)
	task.Content = service.Content
	task.Status = service.Status
	task.Title = service.Title
	code := e.SUCCESS
	err := model.DB.Save(&task).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:  "修改成功",
	}
}

func (service *SearchTaskService) Search() serializer.Response {
	var tasks []model.Task
	code := e.SUCCESS
	err := model.DB.Where("title LIKE ? OR content LIKE ?",
		"%"+service.Info+"%","%"+service.Info+"%").Find(&tasks).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildTasks(tasks),
	}
}
