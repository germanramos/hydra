// Automatically generated by MockGen. DO NOT EDIT!
// Source: etcd_base_repository.go

package mock_repository

import (
	entity "github.com/innotech/hydra/model/entity"
	gomock "github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
)

// Mock of EtcdAccessLayer interface
type MockEtcdAccessLayer struct {
	ctrl     *gomock.Controller
	recorder *_MockEtcdAccessLayerRecorder
}

// Recorder for MockEtcdAccessLayer (not exported)
type _MockEtcdAccessLayerRecorder struct {
	mock *MockEtcdAccessLayer
}

func NewMockEtcdAccessLayer(ctrl *gomock.Controller) *MockEtcdAccessLayer {
	mock := &MockEtcdAccessLayer{ctrl: ctrl}
	mock.recorder = &_MockEtcdAccessLayerRecorder{mock}
	return mock
}

func (_m *MockEtcdAccessLayer) EXPECT() *_MockEtcdAccessLayerRecorder {
	return _m.recorder
}

func (_m *MockEtcdAccessLayer) Delete(key string) error {
	ret := _m.ctrl.Call(_m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEtcdAccessLayerRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockEtcdAccessLayer) Get(key string) (*entity.EtcdBaseModel, error) {
	ret := _m.ctrl.Call(_m, "Get", key)
	ret0, _ := ret[0].(*entity.EtcdBaseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEtcdAccessLayerRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockEtcdAccessLayer) GetAll() (*entity.EtcdBaseModels, error) {
	ret := _m.ctrl.Call(_m, "GetAll")
	ret0, _ := ret[0].(*entity.EtcdBaseModels)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEtcdAccessLayerRecorder) GetAll() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAll")
}

func (_m *MockEtcdAccessLayer) GetCollection() string {
	ret := _m.ctrl.Call(_m, "GetCollection")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockEtcdAccessLayerRecorder) GetCollection() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCollection")
}

func (_m *MockEtcdAccessLayer) Set(entity *entity.EtcdBaseModel) error {
	ret := _m.ctrl.Call(_m, "Set", entity)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEtcdAccessLayerRecorder) Set(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Set", arg0)
}

func (_m *MockEtcdAccessLayer) SetCollection(collection string) {
	_m.ctrl.Call(_m, "SetCollection", collection)
}

func (_mr *_MockEtcdAccessLayerRecorder) SetCollection(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetCollection", arg0)
}
