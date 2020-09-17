package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"database/sql"

	//pb "proto"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

const (
	port = ":50051"
)

type MysqlConfig struct {
	Adress   string `yaml:"address"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ReadYamlConfig(path string) (*MysqlConfig, error) {
	conf := &MysqlConfig{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	return conf, nil
}

func echo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	//fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	output := ""
	for _, v := range r.Form {
		//	fmt.Println("key:", k)
		for _, word := range v {
			output = output + word
		}
	}
	if output == "" {
		return
	}
	fmt.Fprintf(w, output) //这个写入到w的是输出到客户端的
	insertData(output)
}

var DB *sql.DB

func insertData(chat string) {
	statement := fmt.Sprintf("INSERT INTO %s (talk) VALUES ('%s')", "chat_room", chat)
	log.Printf("statement is:%s \n", statement)

	ret, err := DB.Exec(statement)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		fmt.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:", RowsAffected)
	}
}

func sayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	//读取mysql的配置
	conf, err := ReadYamlConfig("config/mysql.ymal")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("conf of msql:%+v", conf)

	db_line := conf.Username + ":" + conf.Password + "@/test?charset=utf8/(" + conf.Adress + ":" + conf.Port + ")"
	db_line = "monty:some_pass@tcp(127.0.0.1:3306)/test?charset=utf8"
	db, err := sql.Open("mysql", db_line)
	if err != nil {
		fmt.Printf("failed to connect mysql:%s\n", err)
		return
	}
	DB = db

	defer DB.Close()

	//设置http服务器
	http.HandleFunc("/", echo)              //设置访问的路由
	err = http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//开启一个grpc
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &pb.GreeterService{SayHello: sayHello})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
