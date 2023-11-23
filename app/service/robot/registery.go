package robot

import (
	"fmt"
	"math/rand"
	"strings"
)

type Registry struct {
	Robots map[string]Robot
}

func NewRegistry() *Registry {
	return &Registry{
		Robots: make(map[string]Robot),
	}
}

func (r *Registry) Register(robot Robot) error {
	if robot.GetRobotType() == "" {
		return fmt.Errorf("robot type is empty")
	}
	r.Robots[robot.GetRobotType()] = robot
	return nil
}

func (r *Registry) Unregister(robotType string) error {
	if _, ok := r.Robots[robotType]; !ok {
		return fmt.Errorf("robot type %s not found", robotType)
	}
	delete(r.Robots, robotType)
	return nil
}

func (r *Registry) GetRobot(robotName string) (robot Robot, ok bool) {
	var robotType string = robotName
	if strings.Contains(robotName, "(com)") {
		robotType = strings.Split(robotName, "(com)")[0]
	}
	robot, ok = r.Robots[robotType]
	return robot, ok
}

func (r *Registry) GetRobotTypes() []string {
	robotTypes := make([]string, 0, len(r.Robots))
	for robotType := range r.Robots {
		robotTypes = append(robotTypes, robotType)
	}
	return robotTypes
}

func (r *Registry) GetRobotCount() int {
	return len(r.Robots)
}

func (r *Registry) IsRobotTypeExist(robotType string) bool {
	_, ok := r.Robots[robotType]
	return ok
}

func (r *Registry) GetRandomRobotType() string {
	robotTypes := r.GetRobotTypes()
	return robotTypes[rand.Intn(len(robotTypes))]
}

func (r *Registry) GetRandomRobot() Robot {
	robot, _ := r.GetRobot(r.GetRandomRobotType())
	return robot
}
