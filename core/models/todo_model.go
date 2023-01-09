package models

type Todo struct {
	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title      string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Task       string `protobuf:"bytes,3,opt,name=task,proto3" json:"task,omitempty"`
	Signature  string `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	IsVerified string `protobuf:"bytes,3,opt,name=is_verified,proto3" json:"is_verified,omitempty"`
}

type Fetch struct {
	Query string `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
}

type File struct {
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Link string `protobuf:"bytes,2,opt,name=link,proto3" json:"link,omitempty"`
}
