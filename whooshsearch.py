#!/usr/bin/env python
# -*- coding: utf-8 -*-
# vim: ts=4 sw=4 tw=99 et:

"""
@Date   : July 21, 2019
@Author : corvo

"""


import os
import json
import logging
from collections import namedtuple

import click
import tornado
import tornado.web
import pymysql

from whoosh import index
from whoosh.qparser import QueryParser, MultifieldParser, FuzzyTermPlugin, MultifieldPlugin
from whoosh.index import create_in
from whoosh.analysis import RegexAnalyzer
from whoosh.fields import *

from jieba.analyse import ChineseAnalyzer

Record = namedtuple('Record', 'id title content')
debug = print
idx_dir = './index'
with open('db.json', 'r') as f:
    db_cfg = json.loads(f.read())
    db = pymysql.connect(**db_cfg)
chinaAnalyzer = ChineseAnalyzer()

# 定义索引schema,确定索引字段
schema = Schema(id=NUMERIC(stored=True),
                title=TEXT(stored=True, analyzer=chinaAnalyzer),
                content=TEXT(analyzer=chinaAnalyzer),
                )


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
                    cursor.execute(
                        "select * from topic where id = %s and active=1", (rlt['id']))
                    topic = cursor.fetchone()
                if topic:
                    highlight_str = rlt.highlights('content', topic['content'])

            q_array.append(
                {
                    'id': rlt['id'],
                    'title': rlt['title'],
                    'content': highlight_str,
                }
            )

    return q_array


class MainHandler(tornado.web.RequestHandler):
    def get(self):
        q_str = self.get_argument('query', '')
        rlt = {'items': []}
        if query:
            rlt['items'] = query_docs(q_str)

        self.write(json.dumps(rlt, ensure_ascii=False))

    def post(self):
        start = self.get_argument('start', 1)
        end = self.get_argument('end', 2)
        try:
            int(start)
            int(end)
        except ValueError as e:
            self.write("error")
            return
        incremental_index(start, end)
        self.write("make index from {} to {} success".format(start, end))


def make_app():
    return tornado.web.Application([
        (r"/", MainHandler),
    ])


@click.group()
def idx(): pass


@click.group()
def searchd(): pass


@idx.command()
@click.option('--start', type=int, default=1, help='topic start id')
@click.option('--end', type=int, default=2, help='topic end id')
def make_index(start, end):
    incremental_index(start, end)


@searchd.command()
@click.option('--query', type=str, help='The string for search')
def query(query):
    if query:
        debug(query_docs(query))


@searchd.command()
@click.option('--port', type=int, default=8888, help='The tornado http server')
def server(port):
    tornado.log.enable_pretty_logging()
    app = make_app()
    app.listen(port)
    logging.info("This time listen to %s", port)
    tornado.ioloop.IOLoop.current().start()


def main():
    # incremental_index(1, 2)
    debug(query_docs('细节'))


if __name__ == "__main__":
    cli = click.CommandCollection(sources=[idx, searchd])
    cli()
