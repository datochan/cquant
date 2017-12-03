# cquant

这是一个个人学习娱乐用的小程序，旨在自己学习和研究股票、基金、债券相关的投资方法。之所以放到网上来也是希望能侥幸碰到一个
和自己一样有此兴趣爱好的人能一起沟通进步。

本程序也非常简单，从网上获取到尽可能准确的原始数据，然后根据 @持有封基 的操作思路进行单因子或复合因子测试，各种策略进行回测，能及时找到
一种或者几种相对比较可靠的投资方式方便自己日后的投资理财。

由于我了解到的基金和债券相关的数据较少，不如股票数据丰富，所以本程序先已股票为主。

## 编译与执行

编译脚本的用法: `BASH ./build.sh -h`

>
> Usage: ./build.sh -o [linux|windows|darwin|freebsd] -a [386|amd64] <br />
>        ./build.sh -h for help.
>

示例:
```
BASH ./build.sh -a amd64 -o windows
```

## 数据目录格式

数据目录结构及对应文件的含义说明:

```
/data
├── base
│   ├── basics.csv.zip     # 通联数据中股票的基础信息(股票代码，名称，上市日期，状态等)
│   └── bonus.csv.zip      # 通达信中的高送转数据
│   └── calendar.csv.zip   # 通联数据的交易日历信息
│   └── st.csv.zip         # 通联数据中st股票信息
│   └── stocks.csv.zip     # 通达信中股票、基金、债券的基础信息
├── daily          # 按日期分类的股票数据 (方便筛选和查询)
│   ├── days               # 通达信中的日线数据
│   ├── fixed              # 自己计算的后复权价 (计算收益及每日涨跌幅等)
│   └── ...                # 根据情况随时扩展
├── history        # 按股票分类的股票数据
│   ├── days               # 每只股票的日线数据
│   ├── fixed              # 自己计算的股票后复权价 (计算收益用)
│   └── mins               # 通达信中每只股票的5分钟线数据
├── report         # 通达信的财报数据(为了方便查询，没有进行格式转换)
└── trade          # 个人的交易记录信息(以便日后计算收益等)
```

大体思路是， 通过现有操盘软件(通达信) 和 通联数据(http://www.datayes.com) 提供的一些免费接口为基础，封装自己的数据源。
* `base`文件夹中存放股票基础数据，如: 股票列表、交易日历、st股、高送转 等信息。
* `history` 中存放从通达信或者通联数据更新下来的原始数据, 及后续计算的复权价及各种指标数据。
* `daily` 中存放的其实就是history文件夹中存放的数据，只是不过是以日期归类，方便筛选。
* `trade` 中存放的是策略模拟的交易记录，以便日后统计收益等等。

程序比较简单，也不想将其搞复杂了。
