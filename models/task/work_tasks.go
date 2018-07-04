package task

import (
	"time"
	"github.com/astaxie/beego/orm"
	"encoding/json"
)

type WorkTasks struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Describe string `json:"describe"`
	Deadline time.Time `json:"deadline"`
	SaveFolderId int `json:"save_folder_id"`
	SaveFolderName string `json:"save_folder_name"`
	Publisher int `json:"publisher"`
	PublisherName string `json:"publisher_name"`
	TaskNum int `json:"task_num"`
	FinishNum int `json:"finish_num"`
	Status int `json:"status"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"auto_now"`
}

type WorkTasksOperator struct {
	Id int `json:"id"`
	Operator int `json:"operator"`
	OperatorName string `json:"operator_name"`
	CoursewareId int `json:"courseware_id"`
	CoursewareName string `json:"courseware_name"`
	UploadTime time.Time `json:"upload_time"`
	WorkTasksId int `json:"work_tasks_id"`
	Status int `json:"status"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"auto_now"`
}

type WorkTasksCc struct {
	Id int `json:"id"`
	Cc int `json:"cc"`
	CcName string `json:"cc_name"`
	WorkTasksId int `json:"work_tasks_id"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"auto_now"`
}

func (wt *WorkTasks) TableName() string {
	return "work_tasks"
}

func (wt *WorkTasksOperator) TableName() string {
	return "work_tasks_operator"
}

func (wt *WorkTasksCc) TableName() string {
	return "work_tasks_cc"
}

func init()  {
	orm.RegisterModel(new(WorkTasks), new(WorkTasksOperator), new(WorkTasksCc))
}

func (wt *WorkTasks) Save(operator, cc []map[string]interface{}) error {
	o := orm.NewOrm()
	o.Begin()

	if _, err := o.Insert(wt); err != nil {
		o.Rollback()

		return err
	}

	var wtos []WorkTasksOperator
	for _, value := range operator {
		recipientId := int(value["id"].(float64))
		wtr := WorkTasksOperator{Operator:recipientId, OperatorName:value["name"].(string), WorkTasksId:wt.Id}
		wtos = append(wtos, wtr)
	}
	if _, err := o.InsertMulti(len(wtos), wtos); err != nil {
		o.Rollback()

		return err
	}

	if len(cc) != 0 {
		var wtcs []WorkTasksCc
		for _, value := range cc {
			ccId := int(value["id"].(float64))
			wtc := WorkTasksCc{Cc:ccId, CcName:value["name"].(string), WorkTasksId:wt.Id}
			wtcs = append(wtcs, wtc)
		}

		if _, err := o.InsertMulti(len(wtcs), wtcs); err != nil {
			o.Rollback()

			return err
		}
	}
	o.Commit()

	return nil
}

func (wt *WorkTasks) Get() ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	var tasks []*WorkTasks

	o := orm.NewOrm()

	if num, err := o.QueryTable(wt).All(&tasks); err != nil {
		return res, err
	} else if num <= 0 {
		return res, err
	}

	var taskIds []int
	for _, value := range tasks {
		taskIds = append(taskIds, value.Id)
	}

	var operators []WorkTasksOperator
	o.QueryTable(new(WorkTasksOperator)).Filter("work_tasks_id__in", taskIds).All(&operators)

	var ccs []WorkTasksCc
	o.QueryTable(new(WorkTasksCc)).Filter("work_tasks_id__in", taskIds).All(&ccs)

	for _, value := range tasks {
		var op []WorkTasksOperator
		for _, ov := range operators {
			if ov.WorkTasksId == value.Id {
				op = append(op, ov)
			}
		}
		var oc []WorkTasksCc
		for _, cv := range ccs {
			if cv.WorkTasksId == value.Id {
				oc = append(oc, cv)
			}
		}
		var maps map[string]interface{}
		jsons, _ := json.Marshal(value)
		json.Unmarshal(jsons, &maps)
		maps["operator"] = op
		maps["cc"] = oc
		res = append(res, maps)
	}

	return res, nil
}

func (wt *WorkTasks) GetInfoById() (map[string]interface{}, error) {
	var res map[string]interface{}
	o := orm.NewOrm()

	if err := o.Read(wt); err != nil {
		return res, err
	}
	jsons, _ := json.Marshal(wt)
	json.Unmarshal(jsons, &res)

	var wto []WorkTasksOperator
	if _, err := o.QueryTable(new(WorkTasksOperator)).Filter("work_tasks_id", wt.Id).All(&wto); err != nil {
		return res, err
	}
	res["operator"] = wto

	var wtc []WorkTasksCc
	if _, err := o.QueryTable(new(WorkTasksCc)).Filter("work_tasks_id", wt.Id).All(&wtc); err != nil {
		return res, err
	}
	res["cc"] = wtc

	return res, nil
}