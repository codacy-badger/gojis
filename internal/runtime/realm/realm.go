package realm

import (
	"sync/atomic"

	"github.com/gojisvm/gojis/internal/runtime/binding"
	"github.com/gojisvm/gojis/internal/runtime/errors"
	"github.com/gojisvm/gojis/internal/runtime/lang"
)

const (
	IntrinsicNameObjectPrototype   = "ObjectPrototype"
	IntrinsicNameFunctionPrototype = "FunctionPrototype"
	IntrinsicNameThrowTypeError    = "ThrowTypeError"
)

var (
	currentRealm atomic.Value // holds *realm.Realm
)

// CurrentRealm returns the current realm as used in the specification.
func CurrentRealm() *Realm { return currentRealm.Load().(*Realm) }

type Realm struct {
	Intrinsics  *lang.Record
	GlobalObj   lang.Value                   // Object or Undefined
	GlobalEnv   lang.Value                   // Object or Undefined
	TemplateMap map[interface{}]*lang.Object // Parse Node -> Object
	HostDefined lang.Value
}

// Type returns lang.TypeInternal.
func (*Realm) Type() lang.Type { return lang.TypeInternal }

// Value returns the Realm itself.
func (r *Realm) Value() interface{} { return r }

func CreateRealm() *Realm {
	r := new(Realm)
	CreateIntrinsics(r)
	r.GlobalObj = lang.Undefined
	r.GlobalEnv = lang.Undefined
	r.TemplateMap = make(map[interface{}]*lang.Object)
	return r
}

func CreateIntrinsics(r *Realm) {
	r.Intrinsics = lang.NewRecord()
	objProto := lang.ObjectCreate(lang.Null)
	r.Intrinsics.SetField(IntrinsicNameObjectPrototype, objProto)
	// FIXME: %ThrowTypeError% as in 8.2.2 and 9.2.9.1

	panic("TODO: 8.2.2")
}

func (r *Realm) GetIntrinsicObject(n string) lang.Value {
	val, ok := r.Intrinsics.GetField(n)
	if !ok {
		return lang.Undefined
	}
	return val.(lang.Value)
}

func (r *Realm) SetRealmGlobalObject(globalObj, thisValue lang.Value) *Realm {
	if globalObj == lang.Undefined {
		panic("TODO: 8.2.3")
	}

	if thisValue == lang.Undefined {
		thisValue = globalObj
	}

	r.GlobalObj = globalObj

	GlobalEnv := binding.NewGlobalEnvironment(globalObj.(*lang.Object), thisValue.(*lang.Object))
	r.GlobalEnv = GlobalEnv

	return r
}

func (r *Realm) SetDefaultGlobalBindings() lang.Value {
	global := r.GlobalObj.(*lang.Object)
	panic("TODO: for every property\n" + `2. For each property of the Global Object specified in clause 18, do
	a. Let name be the String value of the property name.
	b. Let desc be the fully populated data property descriptor for the property containing the specified
	attributes for the property. For properties listed in 18.2, 18.3, or 18.4 the value of the [[Value]] attribute is
	the corresponding intrinsic object from realmRec.
	c. Perform ? DefinePropertyOrThrow(global, name, desc).`)

	return global
}

func GetFunctionRealm(obj *lang.Object) *Realm {
	panic("TODO")
}

func OrdinaryCreateFromConstructor(constructor *lang.Object, intrinsicDefaultProto lang.String, internalSlotsList ...lang.StringOrSymbol) (*lang.Object, errors.Error) {
	proto, err := GetPrototypeFromConstructor(constructor, intrinsicDefaultProto)
	if err != nil {
		return nil, err
	}

	return lang.ObjectCreate(proto, internalSlotsList...), nil
}

func GetPrototypeFromConstructor(constructor *lang.Object, intrinsicDefaultProto lang.String) (*lang.Object, errors.Error) {
	proto, err := lang.Get(constructor, lang.NewStringOrSymbol(lang.NewString("prototype")))
	if err != nil {
		return nil, err
	}

	if proto.Type() != lang.TypeObject {
		realm := GetFunctionRealm(constructor)
		proto = realm.GetIntrinsicObject(intrinsicDefaultProto.Value().(string))
	}
	return proto.(*lang.Object), nil
}
