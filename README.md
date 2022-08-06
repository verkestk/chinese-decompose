# chinese-decompose
clusters chinese characters by similarity

## How to run

```
go run main.go lesson1-1.csv database.tsv > clusters.md
```

This command line tool takes two arguments:
- a path to a CSV file with vocabularly containing characters the user wishes to cluster
- a decomposition database, like `database.tsv` in this repo

The vocabulary CSV file should be in the following format:

```
vocabulary, pinyin, part of speech, translation
```

It's ok to leave the part of speech or translation blank. The tool will still work. However, it's advised to include at least the translation.

For example, vocabulary from Level 1 Part 1 Lesson 1 from Integrated Chinese third edition:
```
你,nǐ,pr,you
好,hǎo,adj,fine; good; nice; O.K.; it's settled
请,qǐng,v,please (polite form of request); to treat or to invite (somebody)
问,wèn,v,to ask (a question)
贵,guì,adj,honorable; expensive
姓,xìng,v/n,(one's) surname is...; to be surnamed; surname
我,wǒ,pr,I; me
呢,ne,qp,(question particle)
小姐,xiǎojiě,n,Miss; young lady
叫,jiào,v,to be called; to call
什么,shénme,qpr,what
名字,míngzi,n,name
先生,xiānsheng,n,Mr.; husband; teacher
李友,Lǐ Yǒu,pn,(a personal name)
李,lǐ,pn,(a surname); plum
王朋,Wáng Péng,pn,(a personal name)
王,wáng,pn,(a surname); king
是,shì,v,to be
老师,lǎoshī,n,teacher
吗,ma,qp,(question particle)
不,bù,adv,not; no
学生,xuésheng,n,student
也,yě,adv,too; also
人,rén,n,people; person
中国,Zhōngguó,pn,China
北京,Běijīng,pn,Beijing
美国,Měiguó,pn,America
纽约,Niǔyuē,pn,New York
```

## Decomposition Database

The `database.tsv` decomposition database included in this repository was downloaded from https://commons.wikimedia.org/wiki/Commons:Chinese_characters_decomposition on 5 August 2022.

## Example output

# Cluster for 子 in position bottom
## 字李学
### 字 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 名字 | míngzi | n | name |
### 李 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 李 | lǐ | pn | (a surname); plum |
| 李友 | Lǐ Yǒu | pn | (a personal name) |
### 学 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 学生 | xuésheng | n | student |
# Cluster for 女 in position left
## 好姓姐
### 好 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 好 | hǎo | adj | fine; good; nice; O.K.; it's settled |
### 姓 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 姓 | xìng | v/n | (one's) surname is...; to be surnamed; surname |
### 姐 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 小姐 | xiǎojiě | n | Miss; young lady |
# Cluster for 口 in position left
## 叫吗呢
### 叫 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 叫 | jiào | v | to be called; to call |
### 吗 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 吗 | ma | qp | (question particle) |
### 呢 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 呢 | ne | qp | (question particle) |
# Cluster for 丨 in position secondary
## 也中
### 也 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 也 | yě | adv | too; also |
### 中 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 中国 | Zhōngguó | pn | China |
# Cluster for 纟 in position left
## 约纽
### 约 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 纽约 | Niǔyuē | pn | New York |
### 纽 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 纽约 | Niǔyuē | pn | New York |
# Cluster for 牛 in position top
## 生先
### 生 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 先生 | xiānsheng | n | Mr.; husband; teacher |
| 学生 | xuésheng | n | student |
### 先 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 先生 | xiānsheng | n | Mr.; husband; teacher |
# Cluster for 亻 in position left
## 什你
### 什 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 什么 | shénme | qpr | what |
### 你 vocabulary
| Term | Pinyin | PoS | Translation |
| --- | --- | --- | --- |
| 你 | nǐ | pr | you |
