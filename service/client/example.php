<?php

class JsonRPC
{
    /**
     * redis实例
     */
    private static $obj = [];

    private function getConnect($host, $port)
    {
        $key = $host . ':' . $port;
        if (!isset(self::$obj[$key])) {
            $conn = fsockopen($host, $port, $errno, $errstr, 3);
            if (! $conn) {
                throw new \Exception('建立连接失败');
            }
            self::$obj[$key] = $conn;
        }
        return self::$obj[$key];
    }

    public function execute($host, $port, $data) {
        $conn = $this->getConnect($host, $port);

        if (! $conn) {
            throw new \Exception('连接不存在');
        }
        $err = fwrite($conn, json_encode(array_merge($data, ['id' => 0])) . "\n");

        if ($err === false) {
            throw new \Exception('调用方法失败' . __METHOD__);
        }

        stream_set_timeout($conn, 0, 3000);
        $line = fgets($conn);
        if ($line === false) {
            return NULL;
        }
        return json_decode($line,true);
    }

    public function SendToConnections($host, $port, $token, array $connections, $msg)
    {
        $data = array(
            'method' => "Server.SendToConnections",
            'params' => [[
                'connections'   => array_values(array_unique($connections)),
                'msg'           => $msg,
                'token'         => $token
            ]],
        );

        return $this->execute($host, $port, $data);
    }

    public function Broadcast($host, $port, $token, $msg)
    {
        $data = array(
            'method' => "Server.Broadcast",
            'params' => [[
                'msg'           => $msg,
                'token'         => $token
            ]],
        );

        return $this->execute($host, $port, $data);
    }

    public function KickConnections($host, $port, $token, array $connections)
    {
        $data = array(
            'method' => "Server.KickConnections",
            'params' => [[
                'connections'   => array_values(array_unique($connections)),
                'token'         => $token
            ]],
        );

        return $this->execute($host, $port, $data);
    }

    public function __construct()
    {
        foreach (self::$obj as $conn) {
            fclose($conn);
        }
    }
}

/** kick conn **/
// if (count($argv) < 2) {
//     throw new Exception('消息参数格式: php ./example.php ...connections');
// }
//
// $client = new JsonRPC();
// $r = $client->KickConnections("192.168.3.165", 8901, 'token', array_slice($argv, 1));
// var_export($r);
// exit;
/** kick conn end **/


/** 发消息 **/
// if (count($argv) < 3) {
//     throw new Exception('消息发送参数格式: php ./example.php msg ...to');
// }
// $client = new JsonRPC();
// $r = $client->SendToConnections("message.demo.com", 8901, 'token', array_slice($argv, 2), $argv[1]);
// var_export($r);
// exit;
/** 发消息 end **/

/** 广播 **/
if (count($argv) < 2) {
    throw new Exception('消息发送参数格式: php ./example.php msg ...to');
}
$client = new JsonRPC();
$r = $client->Broadcast("message.demo.com", 8901, 'token', $argv[1]);
var_export($r);
exit;
/** 广播 end **/
