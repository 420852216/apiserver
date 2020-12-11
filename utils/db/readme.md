
###关联查询步骤
1.查询出根数据

2.判断是否有有数据，如果没有直接返回

3.选取第一条数据，遍历它的每个字段

4.判断是否是匿名字段，如果是直接跳过

5.判断是否有tag: connection, 没有直接跳过

6.判断tag:connection的值是否是横杠分隔两个:string-string，如果不是直接返回错误

7.得到关联关系，connection的值[当前表的字段,关联表的字段]

8.遍历根数据，使用`关联表的字段`获取根数据里面的值

9.把值放进一个数组，当作sql的参数

10.判断字段类型：数组，结构体，指针

11.创建一个用来存放查询结果的切片：reflect.New(reflect.SliceOf(foreignModel.Type()))

12.准备好sql和前面的值数组，使用in，查询关联表

13.遍历查到的关联数据，生成一个map，key是根据connection的`当前表的字段`查询的数据，value是遍历的结果

14.遍历根数据，根据connection的`关联表的字段`获取的数据(关联数据的ID)

15.根据这个ID去13的map里面查找，如果有:设置字段的值为map里的value，没有则不做任何操作。
		