package api

import (
	"context"

	pb "github.com/tonymontanapaffpaff/golang-training-university-grpc/proto/go_proto"
	"github.com/tonymontanapaffpaff/golang-training-university-grpc/server/pkg/data"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourseServer struct {
	data *data.CourseData
}

func NewCourseServer(c data.CourseData) *CourseServer {
	return &CourseServer{data: &c}
}

func (c CourseServer) CreateCourse(ctx context.Context, request *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	if err := checkCreateRequest(request); err != nil {
		log.WithFields(log.Fields{
			"course": request.GetCourse(),
		}).Errorf("got an error when trying to create course, err: %s", err)
		return nil, err
	}
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

		s := status.Newf(codes.Internal, "can't create a new course %v", course)
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}

		return nil, errWithDetails.Err()
	}
	log.WithFields(log.Fields{
		"course": course,
	}).Info("course successfully created")
	return &pb.CreateCourseResponse{Code: code}, nil
}

func (c CourseServer) GetCourse(ctx context.Context, request *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	code := request.GetCode()
	if s := checkCode(code); s != nil {
		log.WithFields(log.Fields{
			"code": code,
		}).Error("got an error when trying to get course, err: %s", s.Err())
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}
		return nil, errWithDetails.Err()
	}
	course, err := c.data.Read(code)
	if err != nil {
		log.WithFields(log.Fields{
			"course": course,
		}).Errorf("got an error when trying to get course, err: %s", err)

		s := status.Newf(codes.Internal, "can't get course with following code %v", code)
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}

		return nil, errWithDetails.Err()
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

		s := status.New(codes.Internal, "can't get courses list")
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}

		return nil, errWithDetails.Err()
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
	if err := checkUpdateRequest(request); err != nil {
		log.WithFields(log.Fields{
			"code":        requestCode,
			"description": requestDescription,
		}).Errorf("got an error when trying to create course, err: %s", err)
		return nil, err
	}
	code, err := c.data.ChangeDescription(requestCode, requestDescription)
	if err != nil {
		log.WithFields(log.Fields{
			"code":        requestCode,
			"description": requestDescription,
		}).Errorf("got an error when trying to update course description, err: %s", err)

		s := status.New(codes.Internal, "can't update course description")
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}

		return nil, errWithDetails.Err()
	}
	log.WithFields(log.Fields{
		"code":        code,
		"description": requestDescription,
	}).Infof("course description with code %d successfully updated", code)
	return &pb.UpdateCourseResponse{Code: code}, nil
}

func (c CourseServer) DeleteCourse(ctx context.Context, request *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	code := request.GetCode()
	if s := checkCode(code); s != nil {
		log.WithFields(log.Fields{
			"code": code,
		}).Error("got an error when trying to delete course, err: %s", s.Err())
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}
		return nil, errWithDetails.Err()
	}
	err := c.data.Delete(code)
	if err != nil {
		log.WithFields(log.Fields{
			"code": code,
		}).Errorf("got an error when trying to delete course, err: %s", err)

		s := status.New(codes.Internal, "can't delete course")
		errWithDetails, err := s.WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "can't convert status to status with details %v", s)
		}

		return nil, errWithDetails.Err()
	}
	log.Infof("course with code %d successfully deleted", code)
	return &pb.DeleteCourseResponse{}, nil
}

func checkCode(code int32) *status.Status {
	if code <= 0 {
		s := status.Newf(codes.InvalidArgument, "not a positive number %d", code)
		return s
	}
	return nil
}

func checkCreateRequest(r *pb.CreateCourseRequest) error {
	if s := checkCode(r.GetCourse().GetCode()); s != nil {
		errWithDetails, err := s.WithDetails(r)
		if err != nil {
			return status.Errorf(codes.Unknown, "can't convert s to s with details %v", s)
		}
		return errWithDetails.Err()
	}
	if description := r.GetCourse().GetDescription(); description == "" {
		s := status.Newf(codes.InvalidArgument, "empty field %d", description)
		errWithDetails, err := s.WithDetails(r)
		if err != nil {
			return status.Errorf(codes.Unknown, "can't convert s to s with details %v", s)
		}
		return errWithDetails.Err()
	}
	if departmentCode := r.GetCourse().GetDepartmentCode(); departmentCode <= 0 {
		s := status.Newf(codes.InvalidArgument, "not a positive number %d", departmentCode)
		errWithDetails, err := s.WithDetails(r)
		if err != nil {
			return status.Errorf(codes.Unknown, "can't convert s to s with details %v", s)
		}
		return errWithDetails.Err()
	}
	if title := r.GetCourse().GetTitle(); title == "" {
		s := status.Newf(codes.InvalidArgument, "empty field %d", title)
		errWithDetails, err := s.WithDetails(r)
		if err != nil {
			return status.Errorf(codes.Unknown, "can't convert s to s with details %v", s)
		}
		return errWithDetails.Err()
	}
	return nil
}

func checkUpdateRequest(r *pb.UpdateCourseRequest) error {
	if s := checkCode(r.GetCode()); s != nil {
		errWithDetails, err := s.WithDetails(r)
		if err != nil {
			return status.Errorf(codes.Unknown, "can't convert s to s with details %v", s)
		}
		return errWithDetails.Err()
	}
	if description := r.GetDescription(); description == "" {
		s := status.Newf(codes.InvalidArgument, "empty field %d", description)
		errWithDetails, err := s.WithDetails(r)
		if err != nil {
			return status.Errorf(codes.Unknown, "can't convert s to s with details %v", s)
		}
		return errWithDetails.Err()
	}
	return nil
}
