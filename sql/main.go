package main

import (
	"database/sql"
	"fmt"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:xc456789110@tcp(106.54.119.58:3306)/sql_test?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

//查询
//为了方便查询，我们事先定义好一个结构体来存储user表的数据
type user struct {
	id   int
	age  int
	name string
}

//单行查询
//单行查询 db.QueryRow() 执行依次查询，并期望返回最多一行结构（即Row）. QueryRow总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误。
//例如,未找到结果
//查询单条数据示例

func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	//非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放

	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Println("scan failed,err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)

}

//多行查询
//多行查询 db.Query执行一次查询，返回多行结果（即Rows）,一般用于执行select命令。参数args表示query中的占位参数
//具体示例代码，查询多条数据示例

func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id>? "
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed,err:%v\n", err)
		return
	}
	//非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	//循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed,err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

//插入数据
//插入、更新和删除操作都使用 Exec 方法

func insertRowDemo() {
	sqlStr := "insert into user(name,age) values (?,?)"
	ret,err := db.Exec(sqlStr,"王五",38)
	if err != nil {
		fmt.Printf("insert failed,err:%v\n",err)
		return
	}
	theID,err := ret.LastInsertId()    //新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed,err:%v\n",err)
		return
	}

	fmt.Printf("insert success,the id id %d.\n",theID)
}


//更新数据
func updateRowDemo()  {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr,39,3)
	if err != nil {
		fmt.Printf("update failed,err:%v\n",err)
		return
	}
	n,err := ret.RowsAffected()   //操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed,err:%v\n",err)
		return
	}
	fmt.Printf("update success,affected rows:%d\n",n)

}

//删除数据

func deleteRowDemo()  {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr,2)
	if err != nil {
		fmt.Printf("delete failed,err:%v\n",err)
		return
	}
	n,err := ret.RowsAffected()   //操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffect failed,err:%v\n",err)
		return
	}
	fmt.Printf("delete success,affected rows:%d\n",n)

}


//mysql预处理
/*
什么是预处理
普通sql语句执行过程:
1、客户端对sql语句进行占位符替换得到完整的sql语句
2、客户端发送完整sql语句到mysql服务端
3、mysql服务端执行完整的sql语句并将结果返回给客户端。
预处理执行过程：
1、把sql语句分成两部分，命令部分与数据部分
2、先把命令部分发送给mysql服务端，mysql服务端进行sql预处理
3、然后把数据部分发送给mysqk服务端，mysql服务端对sql语句进行占位符替换
4、mysql服务端执行完整的sql语句并将结果返回给客户端。
为什么需要预处理：
1、优化mysql服务器重复执行sql的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
2、避免sql注入的问题。

*/

func prepareQueryDemo()  {
	sqlStr := "select id,name,age from user where id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("papare failed,err:%v\n",err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed,err:%v\n",err)
		return
	}
	defer rows.Close()
	//循环读物结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id,&u.name,&u.age)
		if err != nil {
			fmt.Printf("scan failed,err:%v\n",err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n",u.id,u.name,u.age)
	}

}

//插入、更新和删除操作的预处理十分类似，这里以插入操作的预处理为例：
func prepareInsertDemo()  {
	sqlStr := "insert into user(name,age) values (?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed,,err:%v",err)
		return
	}
	defer stmt.Close()
	_,err = stmt.Exec("小王子",18)
	if err != nil {
		fmt.Printf("insert failed,err:%v\n",err)
		return
	}
	_,err = stmt.Exec("沙河娜扎", 19)
	if err != nil {
		fmt.Printf("insert failed,err:%v\n",err)
		return
	}
	fmt.Printf("insert success.")

}

//sql注入问题
//在任何时候都不应该自己拼接sql语句
//这里我们演示一个自行拼接sql语句的示例，编写一个根据name字段查询user表的函数如下：
//sql注入示例

func sqlInjectDemo(name string)  {
	sqlStr := fmt.Sprintf("select id,name,age from user where name='%s'",name)
	fmt.Printf("SQL:%s\n",sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id,&u.name,&u.age)
	if err != nil {
		fmt.Printf("exec failed.err:%v\n",err)
		return
	}
	fmt.Printf("user:%#v\n",u)
}

//go实现mysql事务
/*
事务：一个最小的不可再分的工作单元；通常一个事务对应一个完整的业务，同时这个完整的业务需要执行多次的DML(insert\update\delete)语句共同联合完成。A转账给B，
这里面就需要执行两次update操作。

在mysql中只有使用了Innodb数据库引擎的数据库或表才支持事务。事务处理可以用来维护数据库的完整性，保证成批的sql语句要么全部执行，要么全部不执行。

事务的ACID
通常事务必须满足4个条件(ACID):原子性、一致性、隔离性、持久性

*/

//下面的代码演示了一个简单的事务操作，该事务操作能够确保两次更新操作要么同时成功要么同时失败，不会存在中间状态
//事务操作示例
func transactionDemo()  {
	tx, err := db.Begin()    //开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback()   //回滚
		}
		fmt.Printf("begin trans failed, err:%v\n",err)
		return
	}

	sqlStr1 := "update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1,4)
	if err != nil {
		tx.Rollback()   //回滚
		fmt.Printf("exec sql1 failed, err:%v\n",err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback()   //回滚
		fmt.Printf("exec ret1.RowAffected() failed,err:%v\n",err)
		return
	}

	sqlStr2 := "update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2,5)
	if err != nil {
		tx.Rollback()   //回滚
		fmt.Printf("exec sql2 failed, err:%v\n",err)
		return
	}

	affRow2, err :=ret2.RowsAffected()
	if err != nil {
		tx.Rollback()  //回滚
		fmt.Printf("exec ret1.RowsAffect(),err:%v\n",err)
		return
	}

	fmt.Println(affRow1,affRow2)
	if affRow1 == 1 && affRow2 ==1 {
		fmt.Println("事务提交啦...")
		tx.Commit()   //提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}
	fmt.Println("exec trans success!")

}


func main() {
	err := initDB() //调用输出化数据库的函数
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}

	//queryRowDemo()
	//queryMultiRowDemo()
	//insertRowDemo()
	//updateRowDemo()
	deleteRowDemo()

	//prepareQueryDemo()
	//prepareInsertDemo()
	//sqlInjectDemo("xxx' and (select count(*) from user) <10 #")

	transactionDemo()
}
