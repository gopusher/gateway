package dingtalk

type Configs map[string]*Config

type Connections struct {
	configs     Configs
	connections map[string]*Robot
}

func InitConnections(configs Configs) *Connections {
	connections := make(map[string]*Robot)
	for name, config := range configs {
		robot := NewRobot(config)

		connections[name] = robot
	}

	return &Connections{
		configs:     configs,
		connections: connections,
	}
}

func (conns *Connections) Connection(conn string) *Robot {
	if conn == "" {
		conn = "default"
	}

	return conns.connections[conn]
}
