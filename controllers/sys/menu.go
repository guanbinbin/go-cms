package sys

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/syyongx/php2go"
	"go-cms/controllers"
	"go-cms/models"
	"go-cms/pkg/e"
	"go-cms/pkg/vo"
	"log"
	"strings"
)

type MenuController struct {
	controllers.BaseController
}

func (c *MenuController) Prepare() {

}

/**
获取列表数据
 */
func (c *MenuController) Index() {
	if c.Ctx.Input.IsPost() {
		model := models.NewMenu()
		
		data := c.Ctx.Input.RequestBody
		//json数据封装到user对象中
		err := json.Unmarshal(data, &model)
		if err != nil {
			c.JsonResult(e.ERROR, err.Error())
		}
		
		dataMap := make(map[string]interface{}, 0)
		
		if !php2go.Empty(model.Visible) {
			dataMap["visible"] = model.Visible
		}
		
		if !php2go.Empty(model.MenuName) {
			dataMap["menu_name"] = model.MenuName
		}

		//开始时间
		if !php2go.Empty(model.StartTime) {
			dataMap["start_time"] = model.StartTime
		}

		//结束时间
		if !php2go.Empty(model.EndTime) {
			dataMap["end_time"] = model.EndTime
		}

		if php2go.Empty(model.Page) {
			model.Page = 1
		}else{
			if model.Page <= 0 {
				model.Page = 1
			}
		}

		if php2go.Empty(model.PageSize) {
			model.PageSize = 10
		}else {
			if model.Page <= 0 {
				model.Page = 10
			}
		}

		var orderBy string
		if !php2go.Empty(model.OrderColumnName) && !php2go.Empty(model.OrderType){
			orderBy = strings.Join([]string{model.OrderColumnName,model.OrderType}," ")
		}else {
			orderBy = "created_at DESC"
		}
		
		result, count,err := models.NewMenu().FindByMap((model.Page-1)*model.PageSize, model.PageSize, dataMap,orderBy)
		if err != nil{
			c.JsonResult(e.ERROR, "获取数据失败")
		}
		c.JsonResult(e.SUCCESS, "ok", result, count, model.Page, model.PageSize)
	}
}

/**
创建数据
*/
func (c *MenuController) Create() {
	if c.Ctx.Input.IsPost() {
		model := models.NewMenu()
        data := c.Ctx.Input.RequestBody
		//1.压入数据 json数据封装到对象中
		
		err := json.Unmarshal(data, model)
		if err != nil {
			c.JsonResult(e.ERROR, err.Error())
		}
		
		if err := c.ParseForm(model); err != nil {
			c.JsonResult(e.ERROR, "赋值失败")
		}

		//2.验证
		valid := validation.Validation{}
		if b, _ := valid.Valid(model); !b {
			for _, err := range valid.Errors {
				log.Println(err.Key, err.Message)
			}
			c.JsonResult(e.ERROR, "验证失败")
		}

		//3.插入数据
		if _, err := model.Create(); err != nil {
			c.JsonResult(e.ERROR, "创建失败")
		}
		c.JsonResult(e.SUCCESS, "添加成功")
	}
}

/**
更新数据
*/
func (c *MenuController) Update() {
	model := models.NewMenu()
	data := c.Ctx.Input.RequestBody
	//json数据封装到对象中

	err := json.Unmarshal(data, model)

	if err != nil {
		c.JsonResult(e.ERROR, err.Error())
	}

	//save
	if c.Ctx.Input.IsPut() {
		post, err := models.NewMenu().FindById(model.Id)
		if err != nil||php2go.Empty(post) {
			c.JsonResult(e.ERROR, "没找到数据")
		}
		
		valid := validation.Validation{}
		if b, _ := valid.Valid(model); !b {
			for _, err := range valid.Errors {
				log.Println(err.Key, err.Message)
			}
			c.JsonResult(e.ERROR, "验证失败")
		}
		
		if _, err := model.Update(); err != nil {
			c.JsonResult(e.ERROR, "修改失败")
		}
		c.JsonResult(e.SUCCESS, "修改成功")
	}

	//get
	if c.Ctx.Input.IsPost() {
		res,err := model.FindById(model.Id)
		if err != nil{
			c.JsonResult(e.ERROR, "获取失败")
		}
		c.JsonResult(e.SUCCESS, "获取成功",res)
	}
}

/**
删除数据
*/
func (c *MenuController) Delete() {
    if c.Ctx.Input.IsDelete() {
		model := models.NewMenu()
		data := c.Ctx.Input.RequestBody
		//json数据封装到user对象中
		
		err := json.Unmarshal(data, model)
		
		if err != nil {
			c.JsonResult(e.ERROR, err.Error())
		}
		
		post, err := models.NewMenu().FindById(model.Id)
		if err != nil||php2go.Empty(post) {
			c.JsonResult(e.ERROR, "没找到数据")
		}
		
		if err := model.Delete(); err != nil {
			c.JsonResult(e.ERROR, "删除失败")
		}
		c.JsonResult(e.SUCCESS, "删除成功")
	}
}

func (c *MenuController) BatchDelete() {
	model := models.NewMenu()

	var ids []int
	if err := c.Ctx.Input.Bind(&ids, "ids"); err != nil {
		c.JsonResult(e.ERROR, "赋值失败")
	}

	if err := model.DelBatch(ids); err != nil {
		c.JsonResult(e.ERROR, "删除失败")
	}
	c.JsonResult(e.SUCCESS, "删除成功")
}

func (c *MenuController) Menus()  {
	model := models.NewMenu()
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	
	err := json.Unmarshal(data, model)
	
	if err != nil {
		c.JsonResult(e.ERROR, err.Error())
	}

	if php2go.Empty(model.ParentId){
		menuData,_ := model.FindAll()
		c.JsonResult(e.SUCCESS, "获取成功",constructMenuTrees(menuData,0))
	}else {
		menuData,_ := model.FindAllByParentId(model.ParentId)
		c.JsonResult(e.SUCCESS, "获取成功",constructMenuTrees(menuData,0))
	}
}

func (c *MenuController) FindMenus()  {
	model := models.NewMenu()

	data := c.Ctx.Input.RequestBody

	err := json.Unmarshal(data, model)

	if err != nil {
		c.JsonResult(e.ERROR, err.Error())
	}

	var menus []*vo.TreeList
	if php2go.Empty(model.ParentId) {
		menus = model.FindTopMenu()
	}else {
		menus = model.FindMenus(model.ParentId)
	}

	if err != nil{
		c.JsonResult(e.ERROR, err.Error())
	}
	c.JsonResult(e.SUCCESS, "获取成功",menus)
}

func constructMenuTrees(menus []models.Menu, parentId int) []vo.MenuItem {

	branch := make([]vo.MenuItem, 0)
	
	for  _,menu := range menus {
		if menu.ParentId == parentId{
			childList := constructMenuTrees(menus, menu.Id)

			child := vo.MenuItem{
				Id:menu.Id,
				MenuName:menu.MenuName,
				OrderNum:menu.OrderNum,
				MenuType:menu.MenuType,
				Visible:menu.Visible,
				CreateBy:menu.CreateBy,
				CreatedAt:menu.CreatedAt,
				UpdateBy:menu.UpdateBy,
				Icon:menu.Icon,
				UpdatedAt:menu.UpdatedAt,
				Perms:menu.Perms,
				Remark:menu.Remark,
				Url:menu.Url,
				ParentId:menu.ParentId,
				ChildrenList: childList,
			}
			branch = append(branch, child)
		}
	}
	
	return branch
}