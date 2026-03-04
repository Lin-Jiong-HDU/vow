# 总设计文档

`vow`是一个在终端中背单词的工具。

## 功能

- 每天第一次打开终端的时候，`vow`自动从某个单词表中抽取没背过的单词，自动提示功能暂时不实现，用户需要自己输入`vow`命令来查看当天的单词，也可以写到`zsh`的配置文件中，类似`fastfetch`
- 每天背诵的单词数量可以设置
- 使用`vow done`表示当天的单词背诵完成了
- 单词栈：用户导入单词表后，第一次启动`vow`时会自动生成单词栈，在`~/.vow/tasks/{YYYY/MM/DD}-tasks.json`，每天背诵完成后会把当天的单词从单词栈中移除，并记录到`~/.vow/done/{YYYY/MM/DD}-done.json`，同时更新总的已经背诵的单词列表`~/.vow/done.json`

## 相关文件结构设计

config：~/.vow/config.json

```json
{
  "dailyWordCount": 5
}
```

单词表：`~/.vow/vow.json`

```json
{
  "name": "words",
  "words": [
    {
      "word": "abandon",
      "meaning": "v. 放弃；抛弃；遗弃；放纵；沉溺于",
      "example": "He abandoned his family and went to live in the mountains."
    },
    {
      "word": "ability",
      "meaning": "n. 能力；才能；技能；本领",
      "example": "She has the ability to solve complex problems."
    }
  ]
}
```

记录已背诵的单词：`~/.vow/done/{YYYY/MM/DD}-done.json`

```json
{
  "date": "2024-06-01",
  "words": ["abandon", "ability"]
}
```

总的已经背诵的单词：~/.vow/done.json

```json
{
  "updateTime": "2024-06-01T12:00:00Z",
  "words": ["abandon", "ability"]
}
```

单词栈：`~/.vow/tasks/{YYYY/MM/DD}-tasks.json`

```json
{
  "date": "2024-06-01",
  "words": ["abandon", "ability", "abandon", "ability", "abandon"]
}
```

## 数据流

单词表 -> 单词栈 -> 已背诵单词列表 -> 总的已背诵单词列表

## 相关文件结构设计
