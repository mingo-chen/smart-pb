package smartpb

import (
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
)

func Test_arrayPath(t *testing.T) {
	var path = "users[1]"
	ret := arrayPath.FindStringSubmatch(path)
	t.Logf("ret:%+v", ret)
}

func TestUnmarshal(t *testing.T) {
	var api = &Api2{
		Id: 10086,
		Info: &Api1{
			Name:  "mingo",
			Email: "abc@ef.com",
			Lv:    Level_Hard,
		},
	}

	pj := Marshal(api)
	buf := Unmarshal(pj)
	newPj := Sink(buf)
	if strV, err := newPj.GetString("info.email"); err != nil {
		t.Fatalf("get val err:%+v", err)
	} else {
		t.Logf("get val:%+v", strV)
	}

}

func TestSmart_Pretty(t *testing.T) {
	var api = &Api2{
		Id: 10086,
		Info: &Api1{
			Name:  "mingo",
			Email: "abc@ef.com",
			Lv:    Level_Hard,
		},
	}

	pj := Marshal(api)
	pj.Pretty()

	t.Logf("pj:%+v", pj)
}

func TestPJ_path_array2(t *testing.T) {
	var v31 = &Api3{
		Uuid:  "xxx_yyy-zz-01",
		Likes: []string{"a", "b", "c"},
		Times: 8,
		V2: &Api2{
			Id: 12345678,
			Info: &Api1{
				Name:  "mingochen",
				Email: "a@bb.ccc",
			},
		},
	}

	var v41 = &Api4{
		Users: []*Api2{
			{Id: 100, Info: &Api1{Name: "aa", Email: "aa@t.com"}},
			{Id: 101, Info: &Api1{Name: "bb", Email: "bb@t.com"}},
			{Id: 102, Info: &Api1{Name: "cc", Email: "cc@t.com"}},
		},
	}

	var v42 = &Api4{
		Users: []*Api2{
			{Id: 200, Info: &Api1{Name: "xxx", Email: "xxx@t.com"}},
			{Id: 201, Info: &Api1{Name: "yyy", Email: "yyy@t.com"}},
			{Id: 202, Info: &Api1{Name: "zzz", Email: "zzz@t.com"}},
		},
	}

	var api = &Api5{
		V3: []*Api3{v31},
		V4: []*Api4{v41, v42},
	}

	pj := Marshal(api)
	id, err := pj.GetInt64("v3[0].v2.id")
	t.Logf("id: %+v, err:%+v", id, err) // 12345678

	email, err := pj.GetString("v4[1].users[2].info.email")
	t.Logf("email: %+v, err:%+v", email, err) // zzz@t.com
}

func TestPJ_path_array(t *testing.T) {
	var api = &Api4{
		Users: []*Api2{
			{Id: 100, Info: &Api1{Name: "aa", Email: "aa@t.com"}},
			{Id: 101, Info: &Api1{Name: "bb", Email: "bb@t.com"}},
			{Id: 102, Info: &Api1{Name: "cc", Email: "cc@t.com"}},
		},
		Admin: &Api4Api6{
			V2: &Api2{
				Id: 12345,
				Info: &Api1{
					Name:  "tx",
					Email: "k@te.com",
				},
			},
		},
	}

	pj := Marshal(api)
	id, err := pj.GetInt64("users[1].id")
	t.Logf("id: %+v, err:%+v", id, err) // 101

	email, err := pj.GetString("users[2].info.email")
	t.Logf("email: %+v, err:%+v", email, err) // cc@t.com

	adminId, err := pj.GetInt64("admin.v2.id")
	t.Logf("id: %+v, err:%+v", adminId, err) // 12345

	var adminName string
	err = pj.Get("admin.v2.info.name", &adminName)
	t.Logf("name: %+v, err:%+v", adminName, err) // tx

	var users []*Api2
	if err := pj.Get("users", &users); err != nil {
		t.Fatalf("get users err:%+v", err)
	} else {
		t.Logf("get users:%+v", users)
	}
}

func Test_lv(t *testing.T) {
	var holder Level

	convert(&holder)
	t.Logf("holder:%+v", holder)
}

func convert(holder interface{}) {
	var v = int32(2)
	fmt.Printf("before:%+v\n", holder)
	th := reflect.ValueOf(holder)
	rv := reflect.ValueOf(v)
	fmt.Printf("holder type:%+v, holder kind:%+v, val type:%+v, val kind:%+v\n",
		th.Elem().Type(), th.Elem().Kind(), rv.Type(), rv.Kind())
	nv := rv.Convert(th.Elem().Type())
	th.Elem().Set(nv)

	var api2 = new(Api2)
	fmt.Printf("ProtoReflect:%+v \n", api2.ProtoReflect()) // &{p:{p:<nil>} mi:0xc000120948}
}

func TestSmart_Get(t *testing.T) {
	var api = &Api2{
		Id: 10086,
		Info: &Api1{
			Name:  "mingo",
			Email: "abc@ef.com",
			Lv:    Level_Hard,
		},
		Lvs: []Level{Level_Middle, Level_Hard},
	}

	pj := Marshal(api)
	buf := Unmarshal(pj)
	newPj := Sink(buf)

	// var lv Level
	// if err := newPj.Get("info.lv", &lv); err != nil {
	// 	t.Fatalf("get val err:%+v", err)
	// } else {
	// 	t.Logf("get val:%+v", lv)
	// }

	// var lvs []Level
	// if err := newPj.Get("lvs", &lvs); err != nil {
	// 	t.Fatalf("get val err:%+v", err)
	// } else {
	// 	t.Logf("get val:%+v", lvs)
	// }

	var info = new(Api1)
	if err := newPj.Get("info", info); err != nil {
		t.Fatalf("get val err:%+v", err)
	} else {
		t.Logf("get val:%+v", info)
	}
}

func TestSmart_dynamic_message(t *testing.T) {
	var api = &Api2{
		Id: 10086,
		Info: &Api1{
			Name:  "mingo",
			Email: "abc@ef.com",
			Lv:    Level_Hard,
		},
		Lvs: []Level{Level_Middle, Level_Hard},
	}

	pj := Marshal(api)
	if buf, err := pj.Msg.Marshal(); err != nil {
		t.Fatalf("marshal err:%+v", err)
	} else {
		var h = new(Api2)
		if err := proto.Unmarshal(buf, h); err != nil {
			t.Fatalf("unmarshal err:%+v", err)
		} else {
			t.Logf("val:%+v", h)
		}
	}
}
