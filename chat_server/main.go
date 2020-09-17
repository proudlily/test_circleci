package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"database/sql"

	pb "github.com/proudlily/test_circleci/proto"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
)

const (
	port = ":50051"
)

type Yaml struct {
	Mysql MysqlConfig `yaml:"mysql"`
}
type MysqlConfig struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ReadYamlConfig(fileName string) (*Yaml, error) {
	conf := &Yaml{}

	if f, err := os.Open(fileName); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(conf)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil

}

func echo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	output := ""
	for _, v := range r.Form {
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

type sayHello struct {
}

func (this *sayHello) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Message != "" {
		insertData(in.Message)
	}
	return &pb.HelloReply{Message: in.Message}, nil
}

func httpServer() {
	//设置http服务器
	http.HandleFunc("/", echo)               //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func tcpServer() {
	//开启一个grpc
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("----start tcp")

	s := grpc.NewServer()

	a := &sayHello{}
	pb.RegisterGreeterServer(s, a)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	//读取mysql的配置
	conf, err := ReadYamlConfig("../config/mysql.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}

	//db_line = "monty:some_pass@tcp(127.0.0.1:3306)/test?charset=utf8"
	db_line := conf.Mysql.Username + ":" + conf.Mysql.Password + "@tcp(" + conf.Mysql.Address + ":" + conf.Mysql.Port + ")/test?charset=utf8"
	db, err := sql.Open("mysql", db_line)
	if err != nil {
		fmt.Printf("failed to connect mysql:%s\n", err)
		return
	}
	DB = db

	defer DB.Close()

	go httpServer()
	go tcpServer()

	select {}
}
