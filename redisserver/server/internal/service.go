package internal

import "redisserver/pb"

func (dbm *RedisDB) Stream(stream pb.RedisDBService_StreamServer) error {
	agent := NewDBAgent(stream, dbm.Processor, dbm.router)
	dbm.agents.Store(agent.sessionID, agent)
	agent.Run()
	// 收到退出指令后完成关闭操作
	agent.OnClose()
	dbm.agents.Delete(agent.sessionID)
	return nil
}
