{{$rdoc := rdocdir "doc/types"}}
{{$rdoc.GetDescription (reflect.CallMethod Types .type) .type 1}}