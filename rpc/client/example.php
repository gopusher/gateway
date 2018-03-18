<?php
class JsonRPC
{
    private $conn;

    function __construct($host, $port)
    {
        $this->conn = fsockopen($host, $port, $errno, $errstr, 3);
        if (!$this->conn) {
            throw new Exception('建立连接失败');
        }
    }

    public function sendToConnections($msg, $to)
    {
        if (!$this->conn ) {
            throw new Exception('连接不存在');
        }
        $err = fwrite($this->conn, json_encode(array(
                'method' => "Server.SendToConnections",
                'params' => [json_encode([
                    'to'    => $to,
                    'msg'   => $msg
                ])],

                'id'     => 0,
            )) . "\n");

        if ($err === false) {
            throw new Exception('调用方法失败' . __METHOD__);
        }

        stream_set_timeout($this->conn, 0, 3000);
        $line = fgets($this->conn);
        if ($line === false) {
            return NULL;
        }
        return json_decode($line,true);
    }
}

if (count($argv) < 3) {
    throw new Exception('消息发送参数格式: php ./example.php msg ...to');
}

//todo 这里的 ip 要注意
$client = new JsonRPC("127.0.0.1", 8901);

$r = $client->sendToConnections($argv[1], array_slice($argv, 2));
var_export($r);
