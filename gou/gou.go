package gou

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterTeam
	router *router
	teams  []*RouterTeam //存储所有路由分组
}

type RouterTeam struct {
	prefix      string
	middlewares []HandlerFunc //中间件支持
	parent      *RouterTeam   //支持继承
	engine      *Engine       //Engine接口支持
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterTeam = &RouterTeam{engine: engine}
	engine.teams = []*RouterTeam{engine.RouterTeam}
	return engine
}

func (team *RouterTeam) Team(prefix string) *RouterTeam {
	engine := team.engine
	newTeam := &RouterTeam{
		prefix: team.prefix + prefix,
		parent: team,
		engine: engine,
	}
	engine.teams = append(engine.teams, newTeam)
	return newTeam
}

func (team *RouterTeam) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := team.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	team.engine.router.addRoute(method, pattern, handler)
}

func (team *RouterTeam) GET(pattern string, handler HandlerFunc) {
	team.addRoute("GET", pattern, handler)
}

func (team *RouterTeam) POST(pattern string, handler HandlerFunc) {
	team.addRoute("POST", pattern, handler)
}

func (team *RouterTeam) Use(middlewares ...HandlerFunc) {
	team.middlewares = append(team.middlewares, middlewares...)
}

func (team *RouterTeam) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(team.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//判断文件存在和权限是否准许
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.ReW, c.Req)
	}
}

func (team *RouterTeam) Static(relativePath string, root string) {
	handler := team.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	team.GET(urlPattern, handler)
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, team := range engine.teams {
		if strings.HasPrefix(r.URL.Path, team.prefix) {
			middlewares = append(middlewares, team.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	engine.router.handle(c)
}
