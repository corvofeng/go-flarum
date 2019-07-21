#!/usr/bin/env python
# -*- coding: utf-8 -*-
# vim: ts=4 sw=4 tw=99 et:

"""
@Date   : July 21, 2019
@Author : corvo

"""


import os
from collections import namedtuple
import pymysql

from whoosh import index
from whoosh.qparser import QueryParser, MultifieldParser, FuzzyTermPlugin, MultifieldPlugin
from whoosh.index import create_in
from whoosh.analysis import RegexAnalyzer
from whoosh.fields import *

from jieba.analyse import ChineseAnalyzer

# analyzer = RegexAnalyzer(r"([\u4e00-\u9fa5])|(\w+(\.?\w +)*)")
Record = namedtuple('Record', 'id title content')

chinaAnalyzer = ChineseAnalyzer()

# 定义索引schema,确定索引字段
schema = Schema(id=NUMERIC(stored=True),
                title=TEXT(stored=True, analyzer=chinaAnalyzer),
                content=TEXT(analyzer=chinaAnalyzer),
                )

db = pymysql.connect("localhost", "root", "fengyuhao", "collipa")
debug = print


idx_dir = './index'


def get_index():

    if not os.path.exists(idx_dir):
        os.mkdir(idx_dir)
        # 创建索引对象
        debug("Create index dir")
        ix = create_in(idx_dir, schema)
    else:
        ix = index.open_dir(idx_dir)
    return ix


def incremental_index(start, end):
    ix = get_index()
    with ix.searcher() as searcher:
        writer = ix.writer()
        writer.add_document(
        )
        results = []
        with db.cursor(cursor=pymysql.cursors.DictCursor) as cursor:
            cursor.execute(
                "SELECT * FROM topic WHERE id >= %s and id < %s", (start, end))
            results = cursor.fetchall()
        for row in results:
            writer.delete_by_term('id', row['id'])
            writer.add_document(
                id=row['id'],
                title=row['title'],
                content=row['content'],
            )
        writer.commit()


def query_docs(q_str):
    ix = get_index()

    q_array = []

    # q = QueryParser(["content", ix.schema).parse(q_str)
    # parser = MultifieldParser(["content", "title"], ix.schema)
    p = QueryParser(None, ix.schema)
    p.add_plugin(MultifieldPlugin(['content', 'title']))
    p.add_plugin(FuzzyTermPlugin())
    q = p.parse(q_str)
    with ix.searcher() as searcher:
        results = searcher.search(q)
        for rlt in results:
            highlight_str = ''
            # 首先搜索标题, 而后搜索内容
            if not highlight_str:
                try:
                    highlight_str = rlt.highlights('title')
                except KeyError as e:
                    pass

            if not highlight_str:
                topic = None
                with db.cursor(cursor=pymysql.cursors.DictCursor) as cursor:
                    cursor.execute("select * from topic where id = %s", (rlt['id']))
                    topic = cursor.fetchone()
                if topic:
                    highlight_str = rlt.highlights('content', topic['content'])

            q_array.append(
                Record(
                    id=rlt['id'],
                    title=rlt['title'],
                    content=highlight_str,
                )
            )

    return q_array


def main():
    incremental_index(1, 2)
    debug(query_docs('细节'))


if __name__ == "__main__":
    main()
