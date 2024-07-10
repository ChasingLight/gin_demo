package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

type LoginReq struct {
	Username string `form:"username" json:"user" binding:"required"`
	Password string `form:"password" json:"pwd" binging:"required"`
	Code     string `form:"code" json:"code" binging:"required"`
}

func main() {
	// gin 默认路由引擎
	r := gin.Default()
	//---------------------------加载静态文件--------------------------------
	r.Static("/static", "./statics")
	//---------------------------声明自定义模版函数--------------------------------
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	// 指定用户GET请求访问/hello时，执行sayHello这个函数
	r.GET("/hello", sayHello)

	//---------------------------gin 响应 json--------------------------------
	// 方式1：使用gin.H{}
	r.GET("/person", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name": "jaden",
			"age":  18,
			"sex":  "男"})
	})
	// 方式2：使用 结构体 + tag定制化 + gin序列化
	type Person struct {
		Name string `json:"name"` //使用tag定制化key返回
		Age  int    `json:"age"`
		Sex  string `json:"sex"`
	}
	person := Person{Name: "ribbon", Age: 17, Sex: "女"}
	r.GET("/another_person", func(c *gin.Context) {
		c.JSON(http.StatusOK, person) //json序列化
	})

	//---------------------------RESTful 风格--------------------------------
	r.GET("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GET Book",
		})
	})

	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POST Book",
		})
	})

	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PUT Book",
		})
	})

	r.DELETE("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DELETE Book",
		})
	})

	//---------------------------Gin多个模版渲染--------------------------------
	// 模版定义---已定义在项目根目录下的 templates 目录中
	// 模版解析
	r.LoadHTMLFiles("./templates/posts/index.html",
		"./templates/users/index.html",
		"./templates/index.tmpl",
		"./templates/index.html")
	//r.LoadHTMLGlob("./templates/**/*")
	// 模版渲染
	// 1> 多个同名模版渲染
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(200, "posts/index.html", gin.H{
			"title": "JadenOliver1!",
		})
	})
	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(200, "users/index.html", gin.H{
			"title": "JadenOliver2!",
		})
	})
	// 2> 自定义模版函数---safe
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", "<a href='https://liwenzhou.com'>李文周的博客</a>")
	})
	// 3> 站长之家-前端模版---快速构建响应式网站首页
	r.GET("/home", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	//---------------------------获取querystring参数--------------------------------
	r.GET("/getQueryString", func(c *gin.Context) {
		name := c.Query("name")
		age := c.DefaultQuery("age", "17")
		sex, isExist := c.GetQuery("sex")
		if !isExist {
			sex = "未知"
		}
		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"age":  age,
			"sex":  sex,
		})
	})

	//---------------------------获取form参数--------------------------------
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		code := c.DefaultPostForm("code", "6688")
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"password": password,
			"code":     code,
		})
	})

	//---------------------------获取path参数--------------------------------
	r.GET("/user/search/:username/:address", func(c *gin.Context) {
		username := c.Param("username")
		address := c.Param("address")
		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"address":  address,
		})
	})

	//---------------------------参数绑定：gin.Context.shouldBind()------------------------------
	r.POST("/loginForm", func(c *gin.Context) {
		var loginReq LoginReq
		err := c.ShouldBind(&loginReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, loginReq)
		}
	})
	r.POST("/loginJson", func(c *gin.Context) {
		var loginReq LoginReq
		err := c.ShouldBind(&loginReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, loginReq)
		}
	})
	// 仍以 from tag 作为基准
	// http://127.0.0.1:8080/loginQueryStr?username=小明&password=123&code=111
	r.GET("/loginQueryStr", func(c *gin.Context) {
		var loginReq LoginReq
		err := c.ShouldBind(&loginReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, loginReq)
		}
	})

	//---------------------------文件上传------------------------------
	//r.MaxMultipartMemory = 8
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			msg := fmt.Sprintf("从请求获取文件信息报错：%s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": msg,
			})
		} else {
			// 保存文件到指定目录
			dst := path.Join("./", file.Filename)
			err := c.SaveUploadedFile(file, dst)
			if err != nil {
				msg := fmt.Sprintf("保存上传文件到本地报错：%s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": msg,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "文件上传与保存成功",
				})
			}
		}
	})
	//---------------------------重定向 与 请求转发------------------------------
	r.GET("/redirect", func(c *gin.Context) {
		//fmt.Println("进入 /redirect 方法内部，开始执行重定向")
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com/")
	})

	// 请求转发
	r.GET("/a", func(c *gin.Context) {
		fmt.Println("进入 /a 方法内部，开始执行请求转发")
		c.Request.URL.Path = "/b"
		r.HandleContext(c)
	})
	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello world",
		})
	})
	//---------------------------路由组 与 404 统一格式返回------------------------------
	shopGroup := r.Group("/shop")
	// 路由组 注册中间件
	shopGroup.Use(m1)
	{
		shopGroup.GET("/index", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "/shop/index"})
		})
		shopGroup.GET("/cart", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "/shop/cart"})
		})
		// 嵌套路由组---xx 和 /xx 均行
		xx := shopGroup.Group("xx")
		{
			xx.GET("/oo", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"msg": "/shop/xx/oo"})
			})
		}
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "统一404返回格式",
		})
	})

	//---------------------------中间件------------------------------
	// 单个路由注册中间件
	r.GET("/middleware", m1, m2, func(c *gin.Context) {
		var name string
		if value, exists := c.Get("name"); exists {
			name = value.(string)
		} else {
			name = "匿名用户"
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":  "hello world",
			"name": name,
		})
	})

	//---------------------------GORM连接mysql------------------------------
	gormMysqlTest()
	// 启动路由引擎
	r.Run(":8080")
} //end main

func sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello JadenOliver!",
	})
}

func m1(c *gin.Context) {
	fmt.Println("m1 in...")
	start := time.Now()
	c.Next() //调用后续处理器函数
	cost := time.Since(start)
	fmt.Printf("------m1 cost %s\n", cost)
	fmt.Println("m1 out...")
	fmt.Println()
}

func m2(c *gin.Context) {
	fmt.Println("m2 in...")
	// 可通过c.Set在请求上下文中设置值，后续处理函数可以获取到该值
	c.Set("name", "JadenOliver2")
	c.Next()
	//c.Abort()	//阻止调用后续的处理函数
	fmt.Println("m2 out...")
}

// ------------------------GORM连接mysql相关代码---begin--------------------------
type UserInfo struct {
	ID     uint
	Name   string
	Gender string
	Hobby  string
}

// TableName : 将 结构体 UserInfo 表名设置为 t_user_info
func (u UserInfo) TableName() string {
	return "t_user_info"
}

func gormMysqlTest() {
	// 1.连接本地mysql数据库
	dsn := "root:root@tcp(127.0.0.1:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	log.Println("Database connection established")
	// 新增
	tx := db.Create(&UserInfo{Name: "云中君", Gender: "男", Hobby: "行云布雨"})
	fmt.Printf("新增记录条数：%d\n", tx.RowsAffected)

	// 查询
	var dbUser UserInfo
	//db.First(&dbUser, 5)                 //根据整型主键查询
	db.First(&dbUser, "name = ?", "云中君") //根据name字段进行查找
	fmt.Printf("User info: %+v\n", dbUser)

	//更新
	db.Model(&dbUser).Update("hobby", "行云布雨666") //单个字段更新
	// 多个字段更新
	db.Model(&dbUser).Updates(UserInfo{Name: "云中君666", Gender: "男666"})
	db.Model(&dbUser).Updates(map[string]interface{}{"Hobby": "行云布雨666888"})
	fmt.Printf("User info: %+v\n", dbUser)
	//删除
	//db.Delete(&dbUser)
}

// ------------------------GORM连接mysql相关代码---end--------------------------
