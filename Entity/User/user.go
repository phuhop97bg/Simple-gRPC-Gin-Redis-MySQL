package UserModel

type User struct {
	ID    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" db:"ID"`
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" db:"Name"`
	Age   int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty" db:"Age"`
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty" db:"Email"`
}
