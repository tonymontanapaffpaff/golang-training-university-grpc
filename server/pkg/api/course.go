package api

import (
	"context"

	pb "github.com/tonymontanapaffpaff/golang-training-university-grpc/proto/go_proto"
	"github.com/tonymontanapaffpaff/golang-training-university-grpc/server/pkg/data"

	log "github.com/sirupsen/logrus"
)

type CourseServer struct {
	data *data.CourseData
}

func NewCourseServer(c data.CourseData) *CourseServer {
	return &CourseServer{data: &c}
}

func (c CourseServer) CreateCourse(ctx context.Context, request *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	course := data.Course{
		Code:           request.Course.GetCode(),
		Title:          request.Course.GetTitle(),
		DepartmentCode: request.Course.GetDepartmentCode(),
		Description:    request.Course.GetDescription(),
	}
	code, err := c.data.Add(course)
	if err != nil {
		log.WithFields(log.Fields{
			"course": course,
		}).Errorf("got an error when trying to create course, err: %s", err)
		return nil, err
	}
	log.WithFields(log.Fields{
		"course": course,
	}).Info("course successfully created")
	return &pb.CreateCourseResponse{Code: code}, nil
}

func (c CourseServer) GetCourse(ctx context.Context, request *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	code := request.GetCode()
	course, err := c.data.Read(code)
	if err != nil {
		log.WithFields(log.Fields{
			"course": course,
		}).Errorf("got an error when trying to get course, err: %s", err)
		return nil, err
	}
	log.WithFields(log.Fields{
		"course": course,
	}).Info("got course")
	return &pb.GetCourseResponse{Course: &pb.Course{
		Code:           course.Code,
		Title:          course.Title,
		DepartmentCode: course.DepartmentCode,
		Description:    course.Description,
	}}, err
}

func (c CourseServer) GetAllCourses(ctx context.Context, request *pb.GetAllCoursesRequest) (*pb.GetAllCoursesResponse, error) {
	courses, err := c.data.ReadAll()
	if err != nil {
		log.Errorf("got an error when trying to get course list, err: %s", err)
		return nil, err
	}
	log.WithFields(log.Fields{
		"courses": courses,
	}).Info("got courses list")

	var responseCourses []*pb.Course
	for _, value := range courses {
		course := &pb.Course{
			Code:           value.Code,
			Title:          value.Title,
			DepartmentCode: value.DepartmentCode,
			Description:    value.Description,
		}
		responseCourses = append(responseCourses, course)
	}
	return &pb.GetAllCoursesResponse{Courses: responseCourses}, nil
}

func (c CourseServer) UpdateCourseDescription(ctx context.Context, request *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	requestCode := request.GetCode()
	requestDescription := request.GetDescription()
	code, err := c.data.ChangeDescription(requestCode, requestDescription)
	if err != nil {
		log.WithFields(log.Fields{
			"code":        requestCode,
			"description": requestDescription,
		}).Errorf("got an error when trying to update course description, err: %s", err)
		return nil, err
	}
	log.WithFields(log.Fields{
		"code":        code,
		"description": requestDescription,
	}).Infof("course description with code %d successfully updated", code)
	return &pb.UpdateCourseResponse{Code: code}, nil
}

func (c CourseServer) DeleteCourse(ctx context.Context, request *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	code := request.GetCode()
	err := c.data.Delete(code)
	if err != nil {
		log.WithFields(log.Fields{
			"code": code,
		}).Errorf("got an error when trying to delete course, err: %s", err)
		return nil, err
	}
	log.Infof("course with code %d successfully deleted", code)
	return &pb.DeleteCourseResponse{}, nil
}
