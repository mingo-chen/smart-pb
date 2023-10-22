package smartpb

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
)

func Test_arrayPath(t *testing.T) {
	Convey("arrayPath_ok\n", t, func() {
		var path = "users[1]"
		ret := arrayPath.FindStringSubmatch(path)
		So(len(ret), ShouldEqual, 3)
		So(ret[0], ShouldEqual, "users[1]")
		So(ret[1], ShouldEqual, "users")
		So(ret[2], ShouldEqual, "1")
	})

	Convey("arrayPath_ko\n", t, func() {
		var path = "users[a1]"
		ret := arrayPath.FindStringSubmatch(path)
		So(len(ret), ShouldEqual, 0)
	})
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

	Convey("Unmarshal_ok\n", t, func() {
		msg := Marshal(api)
		So(msg, ShouldNotBeNil)

		buf := Unmarshal(msg)
		newPj := Sink(buf)
		So(newPj, ShouldNotBeNil)

		strV, err := newPj.GetString("info.email")
		So(err, ShouldBeNil)
		So(strV, ShouldEqual, "abc@ef.com")
	})
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

	msg := Marshal(api)
	msg.Pretty()

	t.Logf("pj:%+v", msg)
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

	Convey("array_path_ok\n", t, func() {
		msg := Marshal(api)
		id, err := msg.GetInt64("v3[0].v2.id")
		So(err, ShouldBeNil)
		So(id, ShouldEqual, int64(12345678))

		email, err := msg.GetString("v4[1].users[2].info.email")
		So(err, ShouldBeNil)
		So(email, ShouldEqual, "zzz@t.com")
	})
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

	msg := Marshal(api)

	Convey("GetInt64\n", t, func() {
		id, err := msg.GetInt64("users[1].id")
		So(err, ShouldBeNil)
		So(id, ShouldEqual, 101)
	})

	Convey("GetInt64_nested\n", t, func() {
		adminId, err := msg.GetInt64("admin.v2.id")
		So(err, ShouldBeNil)
		So(adminId, ShouldEqual, 12345)
	})

	Convey("GetString_array\n", t, func() {
		email, err := msg.GetString("users[2].info.email")
		So(err, ShouldBeNil)
		So(email, ShouldEqual, "cc@t.com")
	})

	Convey("Get\n", t, func() {
		var adminName string
		err := msg.Get("admin.v2.info.name", &adminName)
		So(err, ShouldBeNil)
		So(adminName, ShouldEqual, "tx")
	})

	Convey("Get_pb_array\n", t, func() {
		var users []*Api2
		err := msg.Get("users", &users)
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 3)
		t.Logf("users:%+v", users)
	})

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
	msg := Sink(buf)

	Convey("Get_enum\n", t, func() {
		var lv Level
		err := msg.Get("info.lv", &lv)
		So(err, ShouldBeNil)
		So(lv, ShouldEqual, Level_Hard)
	})

	Convey("Get_enums\n", t, func() {
		var lvs []Level
		err := msg.Get("lvs", &lvs)
		So(err, ShouldBeNil)
		So(lvs, ShouldHaveLength, 2)
		So(lvs[0], ShouldEqual, Level_Middle)
		So(lvs[1], ShouldEqual, Level_Hard)
	})

	Convey("Get_pb\n", t, func() {
		var info = new(Api1)
		err := msg.Get("info", info)
		So(err, ShouldBeNil)
		So(info.Name, ShouldEqual, "mingo")
		So(info.Email, ShouldEqual, "abc@ef.com")
		t.Logf("get val:%+v", info)
	})
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

	Convey("Get_pb\n", t, func() {
		msg := Marshal(api)
		So(msg, ShouldNotBeNil)

		buf, err := msg.Msg.Marshal()
		So(err, ShouldBeNil)

		var h = new(Api2)
		err = proto.Unmarshal(buf, h)
		So(err, ShouldBeNil)
		So(h.Id, ShouldEqual, 10086)

		t.Logf("val:%+v", h)
	})

}
