package smartpb

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	protov2 "github.com/golang/protobuf/proto"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// 都是对象数组，没有二维及以上数组
var arrayPath = regexp.MustCompile(`(.*?)\[([0-9]+)\]`)

// Message ...
type Message struct {
	Data     []byte                  `json:"data"`
	Metadata []byte                  `json:"metadata"`
	DMD      *desc.MessageDescriptor `json:"dmd"`
	Msg      *dynamic.Message        `json:"msg"`
}

// Unmarshal PBMD -> Message
func Unmarshal(m *Message) []byte {
	payload := &Payload{
		Pmd:  m.Metadata,
		Data: m.Data,
	}

	buf, _ := proto.Marshal(payload)
	return buf
}

// Sink 直接[]byte -> PBMD
func Sink(buf []byte) *Message {
	var payload = new(Payload)
	if err := proto.Unmarshal(buf, payload); err != nil {
		return nil
	}

	var m = &Message{
		Data:     payload.Data,
		Metadata: payload.Pmd,
	}

	// 继续解析
	buildMetadata(m)
	return m
}

// Marshal Message -> PBMD
func Marshal(msg proto.Message) *Message {
	data, _ := proto.Marshal(msg)
	md, err := getMetaInfo(msg)
	if err != nil {
		return nil
	}

	var m = &Message{
		Data:     data,
		Metadata: md,
	}

	// 继续解析
	buildMetadata(m)
	return m
}

func getMetaInfo(msg proto.Message) ([]byte, error) {
	rv := reflect.ValueOf(msg)
	mtd := rv.MethodByName("ProtoReflect")
	ret := mtd.Call(nil)
	pfm := ret[0].Interface().(protoreflect.Message)

	prd := pfm.Descriptor()
	var has = make(map[string]bool)
	var dps, edps = analyzeDesc(prd, has)
	fdp := &descriptorpb.FileDescriptorProto{
		Syntax:      proto.String(prd.ParentFile().Syntax().String()),
		Name:        proto.String(string(prd.ParentFile().Name())),
		Package:     proto.String(string(prd.ParentFile().Package())),
		MessageType: dps,
		EnumType:    edps,
	}

	buf, _ := proto.Marshal(fdp)
	return buf, nil
}

// analyzeDesc 要递归生成，并且要防止循环
func analyzeDesc(root protoreflect.MessageDescriptor, hasRead map[string]bool) ([]*descriptorpb.DescriptorProto, []*descriptorpb.EnumDescriptorProto) {
	var dps []*descriptorpb.DescriptorProto
	var edps []*descriptorpb.EnumDescriptorProto
	dps = append(dps, protodesc.ToDescriptorProto(root))
	hasRead[string(root.FullName())] = true // 防止循环遍历

	fields := root.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if m := field.Message(); m != nil && !hasRead[string(m.FullName())] {
			tmpDps, tmpEdps := analyzeDesc(field.Message(), hasRead)
			dps = append(dps, tmpDps...)
			edps = append(edps, tmpEdps...)
			hasRead[string(m.FullName())] = true
		}

		if e := field.Enum(); e != nil && !hasRead[string(e.FullName())] {
			edps = append(edps, protodesc.ToEnumDescriptorProto(e))
			hasRead[string(e.FullName())] = true
		}
	}

	return dps, edps
}

func buildMetadata(m *Message) {
	var dfdp = new(descriptorpb.FileDescriptorProto)
	proto.Unmarshal(m.Metadata, dfdp)
	fd, err := protodesc.NewFile(dfdp, nil)
	if err != nil {
		panic(err)
	}
	md := fd.Messages().Get(0)
	//for i := 0; i < md.Fields().Len(); i++ {
	//	field := md.Fields().Get(i)
	//md.Fields().Get(i).(*filedesc.Field).L1.HasEnforceUTF8 = true
	//md.Fields().Get(i).(*filedesc.Field).L1.EnforceUTF8 = false
	//}

	dmd, err := desc.WrapMessage(md)
	if err != nil {
		panic(err)
	} else {
		m.DMD = dmd
	}

	msg := dynamic.NewMessage(dmd)
	if err := msg.Unmarshal(m.Data); err != nil {
		panic(err)
	} else {
		m.Msg = msg
	}
}

// String 序列化为json格式文本
func (m Message) String() string {
	return m.Msg.String()
}

// Pretty 直接格式化打印，用于调试
func (m Message) Pretty() {
	buf, err := m.Msg.MarshalTextIndent()
	if err != nil {
		fmt.Println("<invalid msg>")
	} else {
		fmt.Println(string(buf))
	}
}

// GetBool 从Message中按jsonpath风格取值
func (m Message) GetBool(path string) (bool, error) {
	return getVal[bool](m, path, false)
}

// GetInt8 从Message中按jsonpath风格取值
func (m Message) GetInt8(path string) (int8, error) {
	return getVal[int8](m, path, 0)
}

// GetUint8 从Message中按jsonpath风格取值
func (m Message) GetUint8(path string) (uint8, error) {
	return getVal[uint8](m, path, 0)
}

// GetInt16 从Message中按jsonpath风格取值
func (m Message) GetInt16(path string) (int16, error) {
	return getVal[int16](m, path, 0)
}

// GetUint16 从Message中按jsonpath风格取值
func (m Message) GetUint16(path string) (uint16, error) {
	return getVal[uint16](m, path, 0)
}

// GetInt32 从Message中按jsonpath风格取值
func (m Message) GetInt32(path string) (int32, error) {
	return getVal[int32](m, path, 0)
}

// GetUint32 从Message中按jsonpath风格取值
func (m Message) GetUint32(path string) (uint32, error) {
	return getVal[uint32](m, path, 0)
}

// GetInt64 从Message中按jsonpath风格取值
func (m Message) GetInt64(path string) (int64, error) {
	return getVal[int64](m, path, 0)
}

// GetUint64 从Message中按jsonpath风格取值
func (m Message) GetUint64(path string) (uint64, error) {
	return getVal[uint64](m, path, 0)
}

// GetFloat32 从Message中按jsonpath风格取值
func (m Message) GetFloat32(path string) (float32, error) {
	return getVal[float32](m, path, 0)
}

// GetFloat64 从Message中按jsonpath风格取值
func (m Message) GetFloat64(path string) (float64, error) {
	return getVal[float64](m, path, 0)
}

// GetBytes 从Message中按jsonpath风格取值
func (s Message) GetBytes(path string) ([]byte, error) {
	return getVal[[]byte](s, path, []byte{})
}

// GetString 从Message中按jsonpath风格取值
func (s Message) GetString(path string) (string, error) {
	return getVal[string](s, path, "")
}

// Get 从Message中按jsonpath风格取值
func (m Message) Get(path string, holder interface{}) error {
	node, err := m.getNode(path)
	if err != nil {
		return err
	}

	rth := reflect.TypeOf(holder)
	rtn := reflect.TypeOf(node)
	if rth.Kind() != reflect.Pointer {
		return errors.New("holder must be a pointer") // 必须是指针才能进行修改，设值
	} else {
		rth = rth.Elem()
	}

	fmt.Printf("holder type:%+v, target type:%+v \n", rth.Kind(), rtn.Kind())
	if rtn.Kind() == reflect.Pointer {
		if rtn.Elem().Kind() != rth.Kind() {
			return fmt.Errorf("value of type %v is not assignable to type %v",
				rtn.Kind(), rth.Kind())
		}
	} else if rth.Kind() != rtn.Kind() { // 类型必须相同才能进行Set
		// value of type int64 is not assignable to type int32
		return fmt.Errorf("value of type %v is not assignable to type %v",
			rtn.Kind(), rth.Kind())
	}

	rvh := reflect.ValueOf(holder)
	rvn := reflect.ValueOf(node)
	if rtn.Kind() == reflect.Array || rtn.Kind() == reflect.Slice { // array, slice要单独处理
		fmt.Printf("array len:%+v, htype:%+v, ntype:%+v\n",
			rvn.Len(), rth, rvn.Index(0).Type())
		array := reflect.MakeSlice(rth, rvn.Len(), rvn.Len())
		for i := 0; i < rvn.Len(); i++ {
			dst := array.Index(i)
			fmt.Printf("dst:%+v, type:%+v, zero:%+v\n", dst, dst.Type().Elem(), dst.IsZero())
			var tmp reflect.Value
			if dst.Kind() == reflect.Pointer {
				tmp = reflect.New(dst.Type().Elem())
			} else {
				tmp = reflect.New(dst.Type())
			}

			fmt.Printf("temp dst:%+v, type:%+v, zero:%+v\n", tmp, tmp.Type(), tmp.IsZero())
			rvni := rvn.Index(i).Interface() // interface{} -> 具体类型
			src := reflect.ValueOf(rvni)
			if err := convertTo(src, tmp); err != nil {
				return err
			} else {
				dst.Set(tmp)
			}
		}
		rvh.Elem().Set(array)

	} else {
		return convertTo(rvn, rvh)
	}

	return nil
}

func convertTo(src, dst reflect.Value) error {
	if dst.Type() == src.Type() {
		dst.Elem().Set(src)
		return nil
	} else if !dst.IsZero() && src.CanConvert(dst.Elem().Type()) { // 类型重定义
		dst.Elem().Set(src.Convert(dst.Elem().Type()))
		return nil
	}

	dm, ok := src.Interface().(*dynamic.Message) // Message
	if !ok {
		return fmt.Errorf("can't convert, [%v=>%v]", src.Type(), dst.Type())
	}

	pm, ok := dst.Interface().(protov2.Message)
	if !ok {
		return fmt.Errorf("dst[%T] is not proto.Message", dst.Interface())
	}

	fmt.Printf("interface%+v, Message:%+v \n", dst.Interface(), pm)
	return dm.ConvertTo(pm)
}

func getVal[T any](s Message, path string, defV T) (T, error) {
	if val, err := s.getNode(path); err != nil {
		return defV, err
	} else {
		if v, ok := val.(T); ok {
			return v, nil
		} else {
			return defV, fmt.Errorf("type not match, %T=>%T", val, defV)
		}
	}
}

// getNode 把path按.进行分隔，逐层取出Message中的Node对象
func (m Message) getNode(path string) (interface{}, error) {
	tokens := strings.Split(path, ".")
	var tmpMsg = m.Msg
	for i, token := range tokens {
		ret := arrayPath.FindStringSubmatch(token)
		prefix := strings.Join(tokens[:i+1], ".")
		var value interface{}
		if len(ret) == 3 { // 不会有多维数组，只会是一维对象数组
			token = ret[1]
			idx, _ := strconv.Atoi(ret[2])
			val, err := tmpMsg.TryGetFieldByName(token)
			if err != nil {
				// unknown field name[admin.v2.infos]
				return nil, fmt.Errorf("%+v[%+v]", err, prefix)
			}

			vals := val.([]interface{})
			if len(vals) <= idx { // 增加断言，防止越界
				return nil, fmt.Errorf("array index over, %d < %d", len(vals), idx)
			}
			value = vals[idx]

		} else {
			val, err := tmpMsg.TryGetFieldByName(token)
			if err != nil {
				// unknown field name[admin.v2.infos]
				return nil, fmt.Errorf("%+v[%+v]", err, prefix)
			}

			value = val
		}

		if i < len(tokens)-1 {
			if m, ok := value.(*dynamic.Message); ok {
				tmpMsg = m
			} else {
				return "", fmt.Errorf("invalid path:%+v", path)
			}
		} else {
			return value, nil
		}
	}

	return nil, nil
}
