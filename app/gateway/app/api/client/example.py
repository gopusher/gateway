# -*- coding: UTF-8 -*-
import json
import socket


class Rpc(object):
    host: str
    port: int
    token: str
    conn: socket.socket

    def __init__(self, host, port, token):
        self.host = host
        self.port = port
        self.token = token
        # self.conn = self.get_connect()

    def connect(self):
        if self.conn is None:
            self.conn = socket.create_connection((self.host, self.port), timeout=1.5)
        return self.conn

    def __getattr__(self, method):
        if method in ["conn"]:
            return None

        def func(**kwargs):
            # args = kwargs if len(kwargs) else args
            params = {
                'token': self.token,
            }
            if kwargs is not None:
                params.update(kwargs)

            data = {
                'method': "Server." + method,
                'params': [params],
            }

            return self.execute(data)

        return func

    def execute(self, data):
        self.connect()

        msg_id = 0
        data['id'] = msg_id
        msg = json.dumps(data)
        self.conn.sendall(msg.encode())

        resp = self.read_line()
        if not resp:
            self.close()
            raise Exception("rpc 获取数据失败, Not resp, server gone")

        resp = json.loads(resp)

        if resp.get('id') != msg_id:
            raise Exception("expected id=%s, received id=%s: %s"
                            % (msg_id, resp.get('id'), resp.get('error')))

        if resp.get('error') is not None:
            raise Exception(resp.get('error'))

        data = json.loads(resp.get('result'))
        if data['code'] != '0':
            raise Exception("rpc 获取数据失败: %s" % data['error'])

        return data

    def read_line(self):
        # return self.conn.makefile().readline()

        ret = b''
        while True:
            c = self.conn.recv(1)
            if c == b'\n' or c == b'':
                break
            else:
                ret += c

        return ret.decode("utf-8")

    def close(self):
        self.conn.close()
        self.conn = None


if __name__ == '__main__':
    rpc = Rpc('127.0.0.1', 8901, 'token')
    print(rpc.SendToConnections(connections=['a', 'b']))
